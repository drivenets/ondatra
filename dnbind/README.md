# Drivenets Binding

Drivenets Binding is an implementation of the Ondatra binding interface with support
limited to [Config](https://pkg.go.dev/github.com/openconfig/ondatra/config) interface.


## Flags

Drivenets Binding integration tests requires a `--config` flag be passed that specifies device information.  
Running Drivenets Binding tests requires passing device credential by using `--node_cres` flag.


### Device Credentials

An example of credentials flags:

```
--node_creds=hostname/user/pass
```

An example of yaml configuration file where `id` needs to match testbed `id`:

```
nodes:
  - id: testbed_id
    hostname:  foo
    credentials:
      username: name
      password: pass
```


## Running the Integration Test

To execute the test, you must update config.yaml with your Drivenets device details
and pass both the testbed and config files as flags to the test:

```
go test github.com/openconfig/ondatra/dnbind/integration --testbed=testbed.textproto --config=config.yaml
```

This repo includes an
[example integration test](integration/integration_test.go) that uses the Drivenets
binding, a [testbed file](integration/testbed.textproto) for that test, and a
[mock configuration file](integration/config.yaml) that is matched by the
testbed.
