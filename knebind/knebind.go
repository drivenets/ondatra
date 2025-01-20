// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package knebind provides an Ondatra binding for KNE devices.
package knebind

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/net/context"

	log "github.com/golang/glog"
	"github.com/open-traffic-generator/snappi/gosnappi"
	"github.com/openconfig/gnoigo"
	closer "github.com/openconfig/gocloser"
	"github.com/openconfig/kne/topo"
	"github.com/openconfig/kne/topo/node"
	"github.com/openconfig/ondatra/binding"
	"github.com/openconfig/ondatra/binding/introspect"
	"github.com/openconfig/ondatra/knebind/creds"
	"github.com/openconfig/ondatra/knebind/solver"
	"golang.org/x/crypto/ssh"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	attestzpb "github.com/openconfig/attestz/proto/tpm_attestz"
	enrollzpb "github.com/openconfig/attestz/proto/tpm_enrollz"
	gpb "github.com/openconfig/gnmi/proto/gnmi"
	acctzpb "github.com/openconfig/gnsi/acctz"
	authzpb "github.com/openconfig/gnsi/authz"
	certzpb "github.com/openconfig/gnsi/certz"
	credzpb "github.com/openconfig/gnsi/credentialz"
	pathzpb "github.com/openconfig/gnsi/pathz"
	grpb "github.com/openconfig/gribi/v1/proto/service"
	cpb "github.com/openconfig/kne/proto/controller"
	tpb "github.com/openconfig/kne/proto/topo"
	opb "github.com/openconfig/ondatra/proto"
	p4pb "github.com/p4lang/p4runtime/go/p4/v1"
)

const grpcDialTimeout = 30 * time.Second

// To be stubbed out by unit tests
var (
	userCurrFn = user.Current
)

// New returns a new KNE bind instance.
func New(cfg *Config) (*Bind, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}
	if cfg.Kubeconfig == "" {
		user, err := userCurrFn()
		if err != nil {
			return nil, err
		}
		cfg.Kubeconfig = filepath.Join(user.HomeDir, ".kube/config")
	}
	top, err := topo.Load(cfg.Topology)
	if err != nil {
		return nil, fmt.Errorf("error loading topology: %w", err)
	}
	tm, err := topo.New(top, topo.WithKubecfg(cfg.Kubeconfig))
	if err != nil {
		return nil, fmt.Errorf("error creating topology manager: %w", err)
	}
	return &Bind{cfg: cfg, tm: tm}, nil
}

// Bind implements the ondatra Binding interface for KNE
type Bind struct {
	binding.Binding
	cfg *Config
	tm  topoManager
}

// Reserve implements the binding Reserve method by finding nodes and links in
// the topology specified in the config file that match the requested testbed.
func (b *Bind) Reserve(ctx context.Context, tb *opb.Testbed, runTime time.Duration, waitTime time.Duration, partial map[string]string) (*binding.Reservation, error) {
	resp, err := b.tm.Show(ctx)
	if err != nil {
		return nil, err
	}
	res, err := solver.Solve(ctx, tb, resp.GetTopology(), partial)
	if err != nil {
		return nil, err
	}
	for i, dut := range res.DUTs {
		if dut.Vendor() == opb.Device_DRIVENETS {
			dnDut := &dnDUT{
				ServiceDUT: dut.(*solver.ServiceDUT),
				bind:       b,
			}
			res.DUTs[i] = dnDut
		} else {
			kdut := &kneDUT{
				ServiceDUT: dut.(*solver.ServiceDUT),
				bind:       b,
			}
			res.DUTs[i] = kdut
			if b.cfg.SkipReset {
				continue
			}
			if err := kdut.resetConfig(ctx); err != nil {
				return nil, err
			}
		}
	}
	for i, ate := range res.ATEs {
		res.ATEs[i] = &kneATE{
			ServiceATE: ate.(*solver.ServiceATE),
			bind:       b,
		}
	}
	return res, nil
}

// Release is a no-op because there's no need to reserve local VMs.
func (b *Bind) Release(context.Context) error {
	return nil
}

type dnDUT struct {
	*solver.ServiceDUT
	bind *Bind
	cli  *dnCLI
}

var _ introspect.Introspector = (*dnDUT)(nil)

func (d *dnDUT) Dialer(svc introspect.Service) (*introspect.Dialer, error) {
	svcName, ok := solver.ServiceName(svc)
	if !ok {
		svcName = string(svc)
		log.Warningf("Service %q has no known KNE name, trying %q", svc, svcName)
	}
	svcPB, ok := d.Services[svcName]
	if !ok {
		return nil, fmt.Errorf("service %q not found on DUT %q", svcName, d.Name())
	}
	opts := []grpc.DialOption{grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true}))} // NOLINT
	if creds := d.newRPCCredentials(); creds != nil {
		opts = append(opts, grpc.WithPerRPCCredentials(creds))
	}
	return makeDialer(svcPB, opts...), nil
}

// newRPCCredentials determines the correct credentials used to access a node via rpc
// from a given node name, node vendor, and knebind config. The precedence order for determining
// the credentials is as follows:
//
// 1. credential provided for a specific node by name from the knebind config
// 2. credential provided for a vendor of the node from the knebind config
// 3. credential from the default username and password fields from the knebind config
// 4. no credentials
func (d *dnDUT) newRPCCredentials() *rpcCredentials {
	cfg := d.bind.cfg
	if userPass := cfg.Credentials.Lookup(d.Name(), d.NodeVendor); userPass != nil {
		return &rpcCredentials{userPass}
	}
	// TODO(team): Deprecate username and password fields.
	if cfg.Username != "" {
		return &rpcCredentials{&creds.UserPass{Username: cfg.Username, Password: cfg.Password}}
	}
	return nil
}

func (d *dnDUT) dialGRPC(ctx context.Context, svc introspect.Service, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	dialer, err := d.Dialer(svc)
	if err != nil {
		return nil, err
	}
	log.Infof("Dialing service %q on DUT %s with options %v", svc, d.Name(), opts)
	return dialer.Dial(ctx, opts...)
}

func (dut *dnDUT) DialCLI(ctx context.Context) (binding.CLIClient, error) {
	if dut.cli != nil {
		return dut.cli, nil
	}

	var err error
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, time.Minute)
	defer cancel()

	var timeout time.Duration = 0
	deadline, ok := ctx.Deadline()
	if ok {
		timeout = time.Until(deadline)
	}

	s, err := dut.Service("ssh")
	if err != nil {
		return nil, err
	}
	addr := serviceAddr(s)

	userPass := dut.newRPCCredentials()
	if userPass == nil {
		return nil, errors.New("RunCommand requires node credentials be provided")
	}

	c := &ssh.ClientConfig{
		User:            userPass.Username,
		Auth:            []ssh.AuthMethod{ssh.Password(userPass.Password)},
		Timeout:         timeout,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	dialCli := func() (err error) {
		dut.cli = &dnCLI{}

		if dut.cli.client, err = ssh.Dial("tcp", addr, c); err != nil {
			return err
		}

		if dut.cli.session, err = dut.cli.client.NewSession(); err != nil {
			return err
		}

		if err = dut.cli.session.RequestPty("xterm", 80, 40, modes); err != nil {
			return err
		}

		if dut.cli.stdin, err = dut.cli.session.StdinPipe(); err != nil {
			return err
		}

		if dut.cli.stdout, err = dut.cli.session.StdoutPipe(); err != nil {
			return err
		}

		if err = dut.cli.session.Shell(); err != nil {
			return err
		}

		return nil
	}

	for {
		select {
		case <-ctx.Done():
			return nil, err
		default:
			if err = dialCli(); err != nil {
				dut.cli.Close()
				dut.cli = nil
				log.Errorf("%s: %s; retrying for %s\n", dut.Name(),
					err, time.Until(deadline).Round(time.Millisecond))
				time.Sleep(time.Second * 10)
				break
			}

			// wait for prompt
			if _, err = dut.cli.CommandResult(ctx); err != nil {
				return nil, err
			}

			return dut.cli, nil
		}
	}
}

func (d *dnDUT) DialGNMI(ctx context.Context, opts ...grpc.DialOption) (gpb.GNMIClient, error) {
	conn, err := d.dialGRPC(ctx, introspect.GNMI, opts)
	if err != nil {
		return nil, err
	}
	return gpb.NewGNMIClient(conn), nil
}

func (d *dnDUT) DialGNOI(ctx context.Context, opts ...grpc.DialOption) (gnoigo.Clients, error) {
	return nil, errors.New("GNOI is not supported on Drivenets DUT")
}

func (d *dnDUT) DialGNSI(ctx context.Context, opts ...grpc.DialOption) (binding.GNSIClients, error) {
	return nil, errors.New("GNSI is not supported on Drivenets DUT")
}

func (d *dnDUT) DialGRIBI(ctx context.Context, opts ...grpc.DialOption) (grpb.GRIBIClient, error) {
	return nil, errors.New("GRIBI is not supported on Drivenets DUT")
}

func (d *dnDUT) DialP4RT(ctx context.Context, opts ...grpc.DialOption) (p4pb.P4RuntimeClient, error) {
	return nil, errors.New("P4RT is not supported on Drivenets DUT")
}

func (dut *dnDUT) PushConfig(ctx context.Context, config string, reset bool) error {
	if _, err := dut.DialCLI(ctx); err != nil {
		return err
	}

	log.Infof("Pushing config:\n\"%s\"", config)

	check := func(res binding.CommandResult, err error) error {
		if err == nil && len(res.Error()) == 0 {
			return nil
		}
		// revert current changes
		dut.cli.RunCommand(ctx, "rollback 0")
		// exit configure menu
		dut.cli.RunCommand(ctx, "end")
		// propagate error or stderr
		if err != nil {
			return err
		}
		return errors.New(res.Error())
	}

	commands := []string{"configure"}
	if reset {
		commands = append(commands, "load override factory-default")
	}
	commands = append(commands, strings.Split(config, "\n")...)
	commands = append(commands, []string{"commit check", "commit", "end"}...)

	for _, command := range commands {
		if err := check(dut.cli.RunCommand(ctx, command)); err != nil {
			return err
		}
	}

	return nil
}

type dnCLI struct {
	*binding.AbstractCLIClient
	client  *ssh.Client
	session *ssh.Session
	stdin   io.WriteCloser
	stdout  io.Reader
}

func (c *dnCLI) Close() error {
	if c.stdin != nil {
		c.stdin.Close()
	}
	if c.session != nil {
		c.session.Close()
	}
	if c.client != nil {
		c.client.Close()
	}
	return nil
}

func (c *dnCLI) RunCommand(ctx context.Context, cmd string) (binding.CommandResult, error) {
	cmd = strings.Replace(cmd, "\t", " ", -1) + "\n"
	if _, err := c.stdin.Write([]byte(cmd)); err != nil {
		return nil, err
	}

	return c.CommandResult(ctx)
}

func (c *dnCLI) CommandResult(ctx context.Context) (r *cmdResult, err error) {
	buffer := make([]byte, 10000000)
	r = &cmdResult{}

	for {
		byteCount, err := c.stdout.Read(buffer)
		if err != nil {
			return r, err
		}
		log.Infof(string(buffer[:byteCount]))
		lines := strings.Split(string(buffer[:byteCount]), "\n")

		for _, line := range lines {
			line = strings.TrimRight(line, " \r")
			if len(line) > 0 {
				if line[len(line)-1:] == "#" {
					return r, nil
				}
			}
			r.output += line + "\n"
			if strings.HasPrefix(line, "ERROR:") {
				r.error = r.output
				r.output = ""
			}
		}
	}
}

type kneDUT struct {
	*solver.ServiceDUT
	bind *Bind
}

var _ introspect.Introspector = (*kneDUT)(nil)

func (d *kneDUT) Dialer(svc introspect.Service) (*introspect.Dialer, error) {
	svcName, ok := solver.ServiceName(svc)
	if !ok {
		svcName = string(svc)
		log.Warningf("Service %q has no known KNE name, trying %q", svc, svcName)
	}
	svcPB, ok := d.Services[svcName]
	if !ok {
		return nil, fmt.Errorf("service %q not found on DUT %q", svcName, d.Name())
	}
	opts := []grpc.DialOption{grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true}))} // NOLINT
	if creds := d.newRPCCredentials(); creds != nil {
		opts = append(opts, grpc.WithPerRPCCredentials(creds))
	}
	return makeDialer(svcPB, opts...), nil
}

func (d *kneDUT) resetConfig(ctx context.Context) error {
	// TODO(team): Reduce duplication between this and CLI reset implementation.
	name := d.Name()
	switch err := d.bind.tm.ResetCfg(ctx, name); {
	case status.Code(err) == codes.Unimplemented:
		log.Warningf("Node %q does not support config reset, skipping", name)
	case err != nil:
		return err
	}
	bp, err := fileRelative(d.bind.cfg.Topology)
	if err != nil {
		return fmt.Errorf("failed to find relative path for topology: %v", err)
	}
	node := d.bind.tm.Nodes()[name]
	cd := node.GetProto().GetConfig().GetConfigData()
	if cd == nil {
		log.Infof("Skipping node %q no initial config found", name)
		return nil
	}
	var b []byte
	switch v := cd.(type) {
	case *tpb.Config_Data:
		b = v.Data
	case *tpb.Config_File:
		cPath := v.File
		if !filepath.IsAbs(cPath) {
			cPath = filepath.Join(bp, cPath)
		}
		log.Infof("Pushing configuration %q to %q", cPath, name)
		var err error
		b, err = os.ReadFile(cPath)
		if err != nil {
			return err
		}
	}
	return d.pushConfig(ctx, b)
}

func (d *kneDUT) pushConfig(ctx context.Context, config []byte) error {
	switch err := d.bind.tm.ConfigPush(ctx, d.Name(), bytes.NewBuffer(config)); {
	case status.Code(err) == codes.Unimplemented:
		log.Warningf("Node %q does not support pushing config, skipping", d.Name())
	case err != nil:
		return err
	}
	return nil
}

func (d *kneDUT) DialGNMI(ctx context.Context, opts ...grpc.DialOption) (gpb.GNMIClient, error) {
	conn, err := d.dialGRPC(ctx, introspect.GNMI, opts)
	if err != nil {
		return nil, err
	}
	return gpb.NewGNMIClient(conn), nil
}

func (d *kneDUT) DialGNOI(ctx context.Context, opts ...grpc.DialOption) (gnoigo.Clients, error) {
	conn, err := d.dialGRPC(ctx, introspect.GNOI, opts)
	if err != nil {
		return nil, err
	}
	return gnoigo.NewClients(conn), nil
}

func (d *kneDUT) DialGRIBI(ctx context.Context, opts ...grpc.DialOption) (grpb.GRIBIClient, error) {
	conn, err := d.dialGRPC(ctx, introspect.GRIBI, opts)
	if err != nil {
		return nil, err
	}
	return grpb.NewGRIBIClient(conn), nil
}

func (d *kneDUT) DialP4RT(ctx context.Context, opts ...grpc.DialOption) (p4pb.P4RuntimeClient, error) {
	conn, err := d.dialGRPC(ctx, introspect.P4RT, opts)
	if err != nil {
		return nil, err
	}
	return p4pb.NewP4RuntimeClient(conn), nil
}

// gnsiConn implements the stub builder needed by the Ondatra
// binding.Binding interface.
type gnsiConn struct {
	*binding.AbstractGNSIClients
	conn *grpc.ClientConn
}

func (c *gnsiConn) Authz() authzpb.AuthzClient {
	return authzpb.NewAuthzClient(c.conn)
}

func (c *gnsiConn) Pathz() pathzpb.PathzClient {
	return pathzpb.NewPathzClient(c.conn)
}

func (c *gnsiConn) Certz() certzpb.CertzClient {
	return certzpb.NewCertzClient(c.conn)
}

func (c *gnsiConn) Credentialz() credzpb.CredentialzClient {
	return credzpb.NewCredentialzClient(c.conn)
}

func (c *gnsiConn) Acctz() acctzpb.AcctzClient {
	return acctzpb.NewAcctzClient(c.conn)
}

func (c *gnsiConn) AcctzStream() acctzpb.AcctzStreamClient {
	return acctzpb.NewAcctzStreamClient(c.conn)
}

func (c *gnsiConn) Attestz() attestzpb.TpmAttestzServiceClient {
	return attestzpb.NewTpmAttestzServiceClient(c.conn)
}

func (c *gnsiConn) Enrollz() enrollzpb.TpmEnrollzServiceClient {
	return enrollzpb.NewTpmEnrollzServiceClient(c.conn)
}

func (d *kneDUT) DialGNSI(ctx context.Context, opts ...grpc.DialOption) (binding.GNSIClients, error) {
	conn, err := d.dialGRPC(ctx, introspect.GNSI, opts)
	if err != nil {
		return nil, err
	}
	return &gnsiConn{conn: conn}, nil
}

// DialGRPC dials the service with the specified name.
// It is exported so that test may use a type assertion to dial custom services.
func (d *kneDUT) DialGRPC(ctx context.Context, serviceName string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return d.dialGRPC(ctx, introspect.Service(serviceName), opts)
}

func (d *kneDUT) dialGRPC(ctx context.Context, svc introspect.Service, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	dialer, err := d.Dialer(svc)
	if err != nil {
		return nil, err
	}
	log.Infof("Dialing service %q on DUT %s with options %v", svc, d.Name(), opts)
	return dialer.Dial(ctx, opts...)
}

// serviceAddr returns the external IP address of a KNE service.
func serviceAddr(s *tpb.Service) string {
	return fmt.Sprintf("%s:%d", s.GetOutsideIp(), s.GetOutside())
}

type rpcCredentials struct {
	*creds.UserPass
}

func (r *rpcCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"username": r.Username,
		"password": r.Password,
	}, nil
}

func (r *rpcCredentials) RequireTransportSecurity() bool {
	return true
}

// newRPCCredentials determines the correct credentials used to access a node via rpc
// from a given node name, node vendor, and knebind config. The precedence order for determining
// the credentials is as follows:
//
// 1. credential provided for a specific node by name from the knebind config
// 2. credential provided for a vendor of the node from the knebind config
// 3. credential from the default username and password fields from the knebind config
// 4. no credentials
func (d *kneDUT) newRPCCredentials() *rpcCredentials {
	cfg := d.bind.cfg
	if userPass := cfg.Credentials.Lookup(d.Name(), d.NodeVendor); userPass != nil {
		return &rpcCredentials{userPass}
	}
	// TODO(team): Deprecate username and password fields.
	if cfg.Username != "" {
		return &rpcCredentials{&creds.UserPass{Username: cfg.Username, Password: cfg.Password}}
	}
	return nil
}

func (d *kneDUT) PushConfig(ctx context.Context, config string, reset bool) error {
	if reset {
		if err := d.resetConfig(ctx); err != nil {
			return err
		}
	}
	return d.pushConfig(ctx, []byte(config))
}

func (d *kneDUT) DialCLI(context.Context) (binding.CLIClient, error) {
	return &kneCLI{dut: d}, nil
}

type kneCLI struct {
	*binding.AbstractCLIClient
	dut *kneDUT
}

func (c *kneCLI) RunCommand(_ context.Context, cmd string) (_ binding.CommandResult, rerr error) {
	s, err := c.dut.Service("ssh")
	if err != nil {
		return nil, err
	}
	addr := serviceAddr(s)
	userPass := c.dut.newRPCCredentials()
	if userPass == nil {
		return nil, errors.New("RunCommand requires node credentials be provided")
	}
	log.Infof("Using credentials %v to SSH", userPass)
	cfg := &ssh.ClientConfig{
		User: userPass.Username,
		Auth: []ssh.AuthMethod{ssh.KeyboardInteractive(func(user, instruction string, questions []string, echos []bool) ([]string, error) {
			if len(questions) > 0 {
				return []string{userPass.Password}, nil
			}
			return nil, nil
		})},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", addr, cfg)
	if err != nil {
		return nil, fmt.Errorf("could not dial SSH server %s: %w", addr, err)
	}
	defer closer.Close(&rerr, client.Close, "error closing SSH client")
	session, err := client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("could not create ssh session: %w", err)
	}
	defer closer.Close(&rerr, func() error {
		if err := session.Close(); err != io.EOF {
			return err
		}
		return nil
	}, "error closing SSH session")

	var buf bytes.Buffer
	session.Stdout = &buf
	switch err := session.Run(cmd).(type) {
	case nil:
		return &cmdResult{output: buf.String()}, nil
	case *ssh.ExitError, *ssh.ExitMissingError:
		return &cmdResult{output: buf.String(), error: err.Error()}, nil
	default:
		return nil, err
	}
}

type cmdResult struct {
	*binding.AbstractCommandResult
	output, error string
}

func (r *cmdResult) Output() string {
	return r.output
}

func (r *cmdResult) Error() string {
	return r.error
}

type kneATE struct {
	*solver.ServiceATE
	bind *Bind
}

var _ introspect.Introspector = (*kneATE)(nil)

func (a *kneATE) Dialer(svc introspect.Service) (*introspect.Dialer, error) {
	svcName, ok := solver.ServiceName(svc)
	if !ok {
		return nil, fmt.Errorf("service %q has no known KNE name", svc)
	}
	svcPB, ok := a.Services[svcName]
	if !ok {
		return nil, fmt.Errorf("service %q not found on ATE %q", svcName, a.Name())
	}
	var creds credentials.TransportCredentials
	switch svc {
	case introspect.OTG:
		creds = insecure.NewCredentials()
	case introspect.GNMI:
		creds = credentials.NewTLS(&tls.Config{InsecureSkipVerify: true}) // NOLINT
	default:
		return nil, fmt.Errorf("credential for service %q on ATE %q not found", svc, a.Name())
	}
	return makeDialer(svcPB, grpc.WithTransportCredentials(creds)), nil
}

func (a *kneATE) DialOTG(ctx context.Context, opts ...grpc.DialOption) (gosnappi.Api, error) {
	conn, err := a.dialGRPC(ctx, introspect.OTG, opts)
	if err != nil {
		return nil, err
	}
	api := gosnappi.NewApi()
	api.NewGrpcTransport().
		SetClientConnection(conn).
		SetRequestTimeout(30 * time.Second)
	return api, nil
}

func (a *kneATE) DialGNMI(ctx context.Context, opts ...grpc.DialOption) (gpb.GNMIClient, error) {
	conn, err := a.dialGRPC(ctx, introspect.GNMI, opts)
	if err != nil {
		return nil, err
	}
	return gpb.NewGNMIClient(conn), nil
}

func (a *kneATE) dialGRPC(ctx context.Context, svc introspect.Service, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	dialer, err := a.Dialer(svc)
	if err != nil {
		return nil, err
	}
	log.Infof("Dialing service %q on ATE %s with options %v", svc, a.Name(), opts)
	return dialer.Dial(ctx, opts...)
}

func makeDialer(svcPb *tpb.Service, opts ...grpc.DialOption) *introspect.Dialer {
	return &introspect.Dialer{
		DialFunc: func(ctx context.Context, addr string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
			ctx, cancel := context.WithTimeout(ctx, grpcDialTimeout)
			defer cancel()
			return grpc.DialContext(ctx, addr, opts...)
		},
		DialTarget: serviceAddr(svcPb),
		DialOpts:   opts,
		DevicePort: int(svcPb.GetInside()),
	}
}

func fileRelative(p string) (string, error) {
	bp, err := filepath.Abs(p)
	if err != nil {
		return "", err
	}
	return filepath.Dir(bp), nil
}

type topoManager interface {
	Nodes() map[string]node.Node
	Show(context.Context) (*cpb.ShowTopologyResponse, error)
	ConfigPush(context.Context, string, io.Reader) error
	ResetCfg(context.Context, string) error
}
