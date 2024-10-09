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

package dnbind_test

import (
	"context"
	"path/filepath"
	"regexp"
	"testing"
	"time"

	"github.com/openconfig/ondatra/cli"
	"github.com/openconfig/ondatra/config"
	dninit "github.com/openconfig/ondatra/dnbind/init"
	"github.com/openconfig/ondatra/internal/flags"
)

func TestDrivenetsBinding(t *testing.T) {
	_, err := flags.Parse()
	if err != nil {
		t.Fatalf("failed to parse flags: %s", err.Error())
	}

	bind, err := dninit.Init()
	if err != nil {
		t.Fatalf("failed to create binding: %s", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := bind.Reserve(ctx, nil, 0, 0, nil)
	if err != nil {
		t.Fatalf("failed to reserve binding: %s", err.Error())
	}

	for _, dut := range res.DUTs {
		cli, err := dut.DialCLI(ctx)
		if err != nil {
			t.Fatal(err.Error())
		}

		res, err := cli.RunCommand(ctx, "show system")
		if err != nil {
			t.Fatal(err.Error())
		}
		t.Log(res.Output())

		config := `interfaces
					 ge100-0/0/0.123
					   admin-state enabled
					   vlan-id 123
					 !
				   !`
		dut.PushConfig(ctx, config, false)

		res, err = cli.RunCommand(ctx, "show interfaces ge100-0/0/0.123")
		if err != nil {
			t.Fatal(err.Error())
		}
		t.Log(res.Output())

		config = `no interfaces ge100-0/0/0.123`
		dut.PushConfig(ctx, config, false)

		res, err = cli.RunCommand(ctx, "show interfaces ge100-0/0/0.123")
		if err != nil {
			t.Fatal(err.Error())
		}
		t.Log(res.Output())
	}
}

// TODO: check push config outputs?
func TestDrivenetsVendorConfig(t *testing.T) {
	_, err := flags.Parse()
	if err != nil {
		t.Fatalf("failed to parse flags: %s", err.Error())
	}

	bind, err := dninit.Init()
	if err != nil {
		t.Fatalf("failed to create binding: %s", err.Error())
	}

	ctx := context.Background()

	res, err := bind.Reserve(ctx, nil, 0, 0, nil)
	if err != nil {
		t.Fatalf("failed to reserve binding: %s", err.Error())
	}

	for _, dut := range res.DUTs {
		// updates DUT config with static config below
		config.NewVendorConfig(dut).
			WithDrivenetsText(
				`interfaces
				   ge100-0/0/0.321
				     admin-state enabled
					 vlan-id 321
				   !
				 !`).
			Append(t)

		// updates DUT config with replace config below
		config.NewVendorConfig(dut).
			WithDrivenetsText(
				`interfaces
				   ge100-0/0/0.{{ var "vlan" }}
				     admin-state {{ var "state" }}
					 vlan-id {{ var "vlan" }}
				   !
				 !`).
			WithVarMap(map[string]string{
				"vlan":  "888",
				"state": "disabled",
			}).
			Append(t)

		// update DUT with multi-vendor config
		config.NewVendorConfig(dut).
			WithCienaText(`should skip this`).
			WithCiscoText(`should also skip this`).
			WithAristaText(`should also skip this`).
			WithJuniperText(`should also skip this`).
			WithDrivenetsText(
				`interfaces
				   ge100-0/0/0.333
				     admin-state enabled
					 vlan-id 333
				   !
				 !`).
			Append(t)

		// replace DUT config with static config below
		config.NewVendorConfig(dut).
			WithDrivenetsText(
				`interfaces
				   ge100-0/0/0
				     admin-state enabled
				   !
				!`).
			Push(t)

		// replace DUT config with static config below
		config.NewVendorConfig(dut).
			WithDrivenetsFile(filepath.Join("testdata", "example_config_1.txt")).
			Push(t)
	}
}

func TestDrivenetsCLI(t *testing.T) {
	_, err := flags.Parse()
	if err != nil {
		t.Fatalf("failed to parse flags: %s", err.Error())
	}

	bind, err := dninit.Init()
	if err != nil {
		t.Fatalf("failed to create binding: %s", err.Error())
	}

	ctx := context.Background()

	res, err := bind.Reserve(ctx, nil, 0, 0, nil)
	if err != nil {
		t.Fatalf("failed to reserve binding: %s", err.Error())
	}

	sysname, _ := regexp.Compile("System Name: w")
	for _, dut := range res.DUTs {
		cli := cli.New(dut)

		result := cli.RunResult(t, "show system name")
		if !sysname.MatchString(result.Output()) {
			t.Fatalf("unexpected command output: %s", result)
		}
	}
}
