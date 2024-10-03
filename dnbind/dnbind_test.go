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
	"fmt"
	"testing"

	"github.com/openconfig/ondatra/binding"
	dninit "github.com/openconfig/ondatra/dnbind/init"
	"github.com/openconfig/ondatra/internal/flags"
)

func TestMain(m *testing.M) {
	err := RunDnBindTests(m, dninit.Init)
	if err != nil {
		fmt.Println(err)
	}
}

func RunDnBindTests(m *testing.M, newBindFn func() (binding.Binding, error)) error {
	_, err := flags.Parse()
	if err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	fmt.Println("=== RUN   DnBindTests")
	bind, err := newBindFn()
	if err != nil {
		return fmt.Errorf("failed to create binding: %w", err)
	}

	ctx := context.Background()

	res, err := bind.Reserve(ctx, nil, 0, 0, nil)
	if err != nil {
		return fmt.Errorf("failed to reserve binding: %w", err)
	}

	for _, dut := range res.DUTs {
		cli, err := dut.DialCLI(ctx)
		if err != nil {
			return err
		}

		res, err := cli.RunCommand(ctx, "show system")
		if err != nil {
			return err
		}
		fmt.Println(res.Output())

		config := `interfaces
					 ge100-0/0/0.123
					   admin-state enabled
					   vlan-id 123
					 !
				   !`
		dut.PushConfig(ctx, config, false)

		res, err = cli.RunCommand(ctx, "show interfaces ge100-0/0/0.123")
		if err != nil {
			return err
		}
		fmt.Println(res.Output())

		config = `no interfaces ge100-0/0/0.123`
		dut.PushConfig(ctx, config, false)

		res, err = cli.RunCommand(ctx, "show interfaces ge100-0/0/0.123")
		if err != nil {
			return err
		}
		fmt.Println(res.Output())
	}

	return nil
}
