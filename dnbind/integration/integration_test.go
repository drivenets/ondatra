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

package integration_test

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/openconfig/ondatra"
	dninit "github.com/openconfig/ondatra/dnbind/init"
)

func TestMain(m *testing.M) {
	ondatra.RunTests(m, dninit.Init)
}

func TestDrivenetsConfig(t *testing.T) {
	dut := ondatra.DUT(t, "dut")

	dut.Config().New().
		WithDrivenetsText(
			`interfaces
               ge100-0/0/0.321
                 admin-state enabled
                 vlan-id 321
               !
             !`).
		Append(t)

	// check duplicate commit
	dut.Config().New().
		WithDrivenetsText(
			`interfaces
               ge100-0/0/0.321
                 admin-state enabled
                 vlan-id 321
               !
             !`).
		Append(t)

	// updates DUT config with replace config below
	dut.Config().New().
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
	dut.Config().New().
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
	dut.Config().New().
		WithDrivenetsText(
			`interfaces
               ge100-0/0/0
                 admin-state enabled
                 description {{ var "desc" }}
               !
             !`).
		WithVarValue("desc", "ondatra").
		Push(t)

	// replace DUT config with static config below
	dut.Config().New().
		WithDrivenetsFile(filepath.Join("../testdata", "example_config_1.txt")).
		Push(t)
}

func TestDrivenetsCLI(t *testing.T) {
	cli := ondatra.DUT(t, "dut").CLI()

	result := cli.RunResult(t, "show system name")
	if !strings.Contains(result.Output(), "System Name: ") {
		t.Fatalf("unexpected command output: stdout: %s stderr: %s",
			result.Output(), result.Error())
	}

	result = cli.RunResult(t, "show interfaces ge100-0/0/a")
	if len(result.Error()) == 0 {
		t.Fatalf("unexpected command output: stdout: %s stderr: %s",
			result.Output(), result.Error())
	}

	stdout := cli.Run(t, "show system name")
	if !strings.Contains(stdout, "System Name: ") {
		t.Fatalf("unexpected command output: %s", stdout)
	}
}

func Test_TC_1_1(t *testing.T) {
	// 1.1
	command := ""
	cli := ondatra.DUT(t, "dal00").CLI()
	command = "show system ntp detail; show system hardware | no-more"
	fmt.Printf("Executing command %s\n", command)
	result := cli.RunResult(t, command)
	fmt.Println(result.Output())
	fmt.Println("++++++++++++++++++++++++++++++++++")
	fmt.Println(result.Error())
	//
}

func Test_TC_1_4(t *testing.T) {
	command := ""
	cli := ondatra.DUT(t, "dal00").CLI()
	command = "show interface bundle-60011 | no-more"
	fmt.Printf("Executing command %s\n", command)
	result := cli.RunResult(t, command)
	fmt.Println(result.Output())
	fmt.Printf("entering loop")
	for i := 1; i <= 100; i++ {
		command = "show interface bundle-60011 | no-more"
		fmt.Printf("Executing command %s\n", command)
		result = cli.RunResult(t, command)
		fmt.Print(i)
		fmt.Println(result.Output())
		fmt.Printf("waiting 5 seconds")
		time.Sleep(5000 * time.Millisecond)
	}
}

func Test_TC_1_5(t *testing.T) {
	// configuration error cause exit from CLI
	dut := ondatra.DUT(t, "dal00")
	dut.Config().New().
		WithDrivenetsText(
			`interfaces lo9876
				!
				!
				!`).
		Append(t)
	// the test is not expected to reach this state
	cli := ondatra.DUT(t, "dal00").CLI()
	command := "set cli-timestamp"
	fmt.Printf("Executing command %s\n", command)
	result := cli.RunResult(t, command)
	fmt.Println(result.Output())
	fmt.Println("++++++++++++++++++++++++++++++++++")
	fmt.Println(result.Error())
	//
}

func Test_TC_1_6(t *testing.T) {
	// 1.6
	command := ""
	cli := ondatra.DUT(t, "dal00").CLI()
	command = ""
	fmt.Printf("Executing command %s\n", command)
	result := cli.RunResult(t, command)
	fmt.Println(result.Output())
	fmt.Println("++++++++++++++++++++++++++++++++++")
	fmt.Println(result.Error())
	//
}

func Test_TC_2_1(t *testing.T) {
	// set system name to 1 character
	cli := ondatra.DUT(t, "dal00").CLI()
	command := "show system name"
	fmt.Printf("Executing command %s\n", command)
	result := cli.RunResult(t, command)
	fmt.Println(result.Output())
	fmt.Println("++++++++++++++++++++++++++++++++++")
	fmt.Println(result.Error())
	//
	dut := ondatra.DUT(t, "dal00")
	// updates DUT config with replace config below
	dut.Config().New().
		WithDrivenetsText(
			`system name A
				 `).
		Append(t)
	//
	command = "show system name"
	fmt.Printf("Executing command %s\n", command)
	result = cli.RunResult(t, command)
	fmt.Println(result.Output())
	fmt.Println("++++++++++++++++++++++++++++++++++")
	fmt.Println(result.Error())
	//
	command = "show system | no-more"
	fmt.Printf("Executing command %s\n", command)
	result = cli.RunResult(t, command)
	fmt.Println(result.Output())
	fmt.Println("++++++++++++++++++++++++++++++++++")
	fmt.Println(result.Error())
	//

}

func Test_TC_2_2(t *testing.T) {
	// set system name to 48 character
	cli := ondatra.DUT(t, "dal00").CLI()
	command := "show system name"
	fmt.Printf("Executing command %s\n", command)
	result := cli.RunResult(t, command)
	fmt.Println(result.Output())
	fmt.Println("++++++++++++++++++++++++++++++++++")
	fmt.Println(result.Error())
	//
	dut := ondatra.DUT(t, "dal00")
	// updates DUT config with replace config below
	dut.Config().New().
		WithDrivenetsText(
			`system name X9F7qqM6LQ8D3Z1_5V0C4B7feW2J6K5T3H9S1A0U3Y6L8R7
	
				 `).
		Append(t)
	//
	command = "show system name"
	fmt.Printf("Executing command %s\n", command)
	result = cli.RunResult(t, command)
	fmt.Println(result.Output())
	fmt.Println("++++++++++++++++++++++++++++++++++")
	fmt.Println(result.Error())
	//
	command = "show system | no-more"
	fmt.Printf("Executing command %s\n", command)
	result = cli.RunResult(t, command)
	fmt.Println(result.Output())
	fmt.Println("++++++++++++++++++++++++++++++++++")
	fmt.Println(result.Error())
	//

}

func Test_TC_2_3(t *testing.T) {
	// 2.3
	cli := ondatra.DUT(t, "dal00").CLI()
	command := "set cli-timestamp"
	fmt.Printf("Executing command %s\n", command)
	result := cli.RunResult(t, command)
	fmt.Println(result.Output())
	fmt.Println("++++++++++++++++++++++++++++++++++")
	fmt.Println(result.Error())
	//
	command = "show system name"
	fmt.Printf("Executing command %s\n", command)
	result = cli.RunResult(t, command)
	fmt.Println(result.Output())
	fmt.Println("++++++++++++++++++++++++++++++++++")
	fmt.Println(result.Error())
	//
	dut := ondatra.DUT(t, "dal00")
	// updates DUT config with replace config below
	dut.Config().New().
		WithDrivenetsText(
			`interfaces lo9876
				 !
				 !`).
		Append(t)

}

func Test_TC_2_4(t *testing.T) {
	// set login banner with #
	dut := ondatra.DUT(t, "dal00")
	dut.Config().New().
		WithDrivenetsText(
			`system login banner "######################\n#                    #\n#   TC 2.4           #\n######################"
				`).
		Append(t)
	cli := ondatra.DUT(t, "dal00").CLI()
	command := "show system name"
	fmt.Printf("Executing command %s\n", command)
	result := cli.RunResult(t, command)
	fmt.Println(result.Output())
	fmt.Println("++++++++++++++++++++++++++++++++++")
	fmt.Println(result.Error())
	//
	dut2 := ondatra.DUT(t, "dal00")
	dut2.Config().New().
		WithDrivenetsText(
			`system login banner "######################\n#                    #\n#   TC 2.4           #\nDAL(cfg)#"
			`).
		Append(t)
	cli2 := ondatra.DUT(t, "dal00").CLI()
	command = "show system name"
	fmt.Printf("Executing command %s\n", command)
	result = cli2.RunResult(t, command)
	fmt.Println(result.Output())
	fmt.Println("++++++++++++++++++++++++++++++++++")
	fmt.Println(result.Error())
	//
	dut3 := ondatra.DUT(t, "dal00")
	dut3.Config().New().
		WithDrivenetsText(
			`system login banner "######################\n#                    #\n#   TC 2.4           #\nERROR:#"
			`).
		Append(t)
	cli3 := ondatra.DUT(t, "dal00").CLI()
	command = "show system name"
	fmt.Printf("Executing command %s\n", command)
	result = cli3.RunResult(t, command)
	fmt.Println(result.Output())
	fmt.Println("++++++++++++++++++++++++++++++++++")
	fmt.Println(result.Error())
	//
}

func Test_TC_2_5(t *testing.T) {
	// change device name during execution
	cli := ondatra.DUT(t, "dal00").CLI()
	command := "show system name"
	fmt.Printf("Executing command %s\n", command)
	result := cli.RunResult(t, command)
	fmt.Println(result.Output())
	fmt.Println("++++++++++++++++++++++++++++++++++")
	fmt.Println(result.Error())
	//
	dut := ondatra.DUT(t, "dal00")
	// updates DUT config with replace config below
	dut.Config().New().
		WithDrivenetsText(
			`system name DALLAS00
				 `).
		Append(t)
	//
	command = "show system name"
	fmt.Printf("Executing command %s\n", command)
	result = cli.RunResult(t, command)
	fmt.Println(result.Output())
	fmt.Println("++++++++++++++++++++++++++++++++++")
	fmt.Println(result.Error())
	//
	command = "show system | no-more"
	fmt.Printf("Executing command %s\n", command)
	result = cli.RunResult(t, command)
	fmt.Println(result.Output())
	fmt.Println("++++++++++++++++++++++++++++++++++")
	fmt.Println(result.Error())
	//

}

func Test_TC_2_6(t *testing.T) {
	// set logging terminal and trigger command that generates print to terminal. dut must have established isis adjacencies
	cli := ondatra.DUT(t, "dal00").CLI()
	command := "set logging terminal"
	fmt.Printf("Executing command %s\n", command)
	result := cli.RunResult(t, command)
	fmt.Println(result.Output())
	fmt.Println("++++++++++++++++++++++++++++++++++")
	fmt.Println(result.Error())
	//
	command = "show isis neighbors"
	fmt.Printf("Executing command %s\n", command)
	result = cli.RunResult(t, command)
	fmt.Println(result.Output())
	fmt.Println("++++++++++++++++++++++++++++++++++")
	fmt.Println(result.Error())
	//
	command = "clear isis neighbor"
	fmt.Printf("Executing command %s\n", command)
	result = cli.RunResult(t, command)
	fmt.Println(result.Output())
	fmt.Println("++++++++++++++++++++++++++++++++++")
	fmt.Println(result.Error())
	//
	command = "show isis neighbors"
	fmt.Printf("Executing command %s\n", command)
	result = cli.RunResult(t, command)
	fmt.Println(result.Output())
	fmt.Println("++++++++++++++++++++++++++++++++++")
	fmt.Println(result.Error())
	//

}

func Test_TC_2_7(t *testing.T) {
	// change device name to "ERROR"
	cli := ondatra.DUT(t, "dal00").CLI()
	command := "show system name"
	fmt.Printf("Executing command %s\n", command)
	result := cli.RunResult(t, command)
	fmt.Println(result.Output())
	fmt.Println("++++++++++++++++++++++++++++++++++")
	fmt.Println(result.Error())
	//
	dut := ondatra.DUT(t, "dal00")
	// updates DUT config with replace config below
	dut.Config().New().
		WithDrivenetsText(
			`system name ERROR
				 `).
		Append(t)
	//
	command = "show system name"
	fmt.Printf("Executing command %s\n", command)
	result = cli.RunResult(t, command)
	fmt.Println(result.Output())
	fmt.Println("++++++++++++++++++++++++++++++++++")
	fmt.Println(result.Error())
	//
	command = "show system | no-more"
	fmt.Printf("Executing command %s\n", command)
	result = cli.RunResult(t, command)
	fmt.Println(result.Output())
	fmt.Println("++++++++++++++++++++++++++++++++++")
	fmt.Println(result.Error())
	//

}

func Test_TC_3_1(t *testing.T) {
	// execute 100 CLI commands
	cli := ondatra.DUT(t, "dal00").CLI()
	command := "show system hardware | no-more"
	for i := 4001; i <= 4101; i++ {
		fmt.Printf("Executing command %s\n", command)
		result := cli.RunResult(t, command)
		fmt.Println(result.Output())
		fmt.Println("++++++++++++++++++++++++++++++++++")
		fmt.Println(result.Error())
		//
	}
}

func Test_TC_3_3(t *testing.T) {
	S := ""
	dut := ondatra.DUT(t, "dal00")
	for i := 4001; i <= 4101; i++ {
		S = strconv.Itoa(i)
		// updates DUT config with replace config below
		dut.Config().New().
			WithDrivenetsText(
				`interfaces
					 lo{{var "lid"}}
					 description "loopback {{ var "lid" }}"
					 !`).
			WithVarMap(map[string]string{
				"lid":         S,
				"description": "fofo",
			}).
			Append(t)
		//
	}
	cli := ondatra.DUT(t, "dal00").CLI()
	command := "show system name"
	fmt.Printf("Executing command %s\n", command)
	result := cli.RunResult(t, command)
	fmt.Println(result.Output())
	fmt.Println("++++++++++++++++++++++++++++++++++")
	fmt.Println(result.Error())
	//
	command = "show interfaces | include lo4 | count"
	fmt.Printf("Executing command %s\n", command)
	result = cli.RunResult(t, command)
	fmt.Println(result.Output())
	fmt.Println("++++++++++++++++++++++++++++++++++")
	fmt.Println(result.Error())
	//

}

func Test_TC_4(t *testing.T) {
	//4.x
	command := ""
	cli := ondatra.DUT(t, "dal00").CLI()
	//
	command = "set cli-terminal-length 0"
	fmt.Printf("Executing command %s\n", command)
	result := cli.Run(t, command)
	fmt.Println(result)
	// T1
	command = "clear isis neighbor"
	fmt.Printf("Executing command %s\n", command)
	result = cli.Run(t, command)
	fmt.Println(result)
	//verification
	command = "show isis neighbors detail | i Last"
	fmt.Printf("Executing command %s\n", command)
	result = cli.Run(t, command)
	fmt.Println(result)
	// T2
	command = "show interfaces counters"
	fmt.Printf("Executing command %s\n", command)
	result = cli.Run(t, command)
	fmt.Println(result)

	command = "clear interfaces counters"
	fmt.Printf("Executing command %s\n", command)
	result = cli.Run(t, command)
	fmt.Println(result)
	//verification
	command = "show interfaces counters"
	fmt.Printf("Executing command %s\n", command)
	result = cli.Run(t, command)
	fmt.Println(result)
	// T3
	command = "request mpls traffic-engineering pcep revoke"
	fmt.Printf("Executing command %s\n", command)
	result = cli.Run(t, command)
	fmt.Println(result)
	// T4
	command = "request system tech-support T1 component routing after 2024-11-03T00:00:00 force"
	fmt.Printf("Executing command %s\n", command)
	result = cli.Run(t, command)
	fmt.Println(result)
	//verification
	command = "show system tech-support status"
	fmt.Printf("Executing command %s\n", command)
	result = cli.Run(t, command)
	fmt.Println(result)
	// T5
	command = "show interfaces bundle-60011"
	fmt.Printf("Executing command %s\n", command)
	result = cli.Run(t, command)
	fmt.Println(result)
	// set cli timestamp
	command = "set cli-timestamp"
	fmt.Printf("Executing command %s\n", command)
	result = cli.Run(t, command)
	fmt.Println(result)
	// T6
	command = "show system ntp detail"
	fmt.Printf("Executing command %s\n", command)
	result = cli.Run(t, command)
	fmt.Println(result)
	// T7
	command = "show dnos-internal oper-data /drivenets-top/interfaces/interface[name='bundle-60011.1']/oper-items/counters"
	fmt.Printf("Executing command %s\n", command)
	result = cli.Run(t, command)
	fmt.Println(result)
	// T8
	command = "run ping 10.0.2.0 count 5"
	fmt.Printf("Executing command %s\n", command)
	result = cli.Run(t, command)
	fmt.Println(result)
	// T9
	command = "show route protocol bgp"
	fmt.Printf("Executing command %s\n", command)
	result = cli.Run(t, command)
	fmt.Println(result)
	// T10
	command = "show config"
	fmt.Printf("Executing command %s\n", command)
	result = cli.Run(t, command)
	fmt.Println(result)
	// T11
	command = "show file log routing_engine/system-events.log"
	fmt.Printf("Executing command %s\n", command)
	result = cli.Run(t, command)
	fmt.Println(result)
	// T12
	command = "show interfaces bundle-60011 ; show interfaces bundle-60008"
	fmt.Printf("Executing command %s\n", command)
	result = cli.Run(t, command)
	fmt.Println(result)
	// T13
	command = "show system ntp detail ; show system hardware | no-more"
	fmt.Printf("Executing command %s\n", command)
	result = cli.Run(t, command)
	fmt.Println(result)
	// T14
	command = "sho rou protocol bg | n"
	fmt.Printf("Executing command %s\n", command)
	result = cli.Run(t, command)
	fmt.Println(result)
	// T15
	command = "sh int cou | n"
	fmt.Printf("Executing command %s\n", command)
	result = cli.Run(t, command)
	fmt.Println(result)
	// T16
	command = "cl is n "
	fmt.Printf("Executing command %s\n", command)
	result = cli.Run(t, command)
	fmt.Println(result)

}

func Test_TC_4_5_1(t *testing.T) {
	// execute CLI command with error
	cli := ondatra.DUT(t, "dal00").CLI()
	command := "show bgp nei 1.1.1"
	fmt.Printf("Executing command %s\n", command)
	result := cli.Run(t, command)
	fmt.Println(result)
}

func Test_TC_4_7_1(t *testing.T) {
	// execute CLI command exit
	cli := ondatra.DUT(t, "dal00").CLI()
	command := "exit"
	fmt.Printf("Executing command %s\n", command)
	result := cli.Run(t, command)
	fmt.Println(result)
}

func Test_TC_4_7_2(t *testing.T) {
	// execute CLI command quit
	cli := ondatra.DUT(t, "dal00").CLI()
	command := "quit"
	fmt.Printf("Executing command %s\n", command)
	result := cli.Run(t, command)
	fmt.Println(result)
}

func Test_TC_5(t *testing.T) {
	// 5X
	command := ""
	cli := ondatra.DUT(t, "dal00").CLI()
	//
	command = "set cli-terminal-length 0"
	fmt.Printf("Executing command %s\n", command)
	result := cli.RunResult(t, command)
	fmt.Println(result.Output())
	// T1
	command = "clear isis neighbor"
	fmt.Printf("Executing command %s\n", command)
	result = cli.RunResult(t, command)
	fmt.Println(result.Output())
	//verification
	command = "show isis neighbors detail | i Last"
	fmt.Printf("Executing command %s\n", command)
	result = cli.RunResult(t, command)
	fmt.Println(result.Output())
	// T2
	command = "show interfaces counters"
	fmt.Printf("Executing command %s\n", command)
	result = cli.RunResult(t, command)
	fmt.Println(result.Output())

	command = "clear interfaces counters"
	fmt.Printf("Executing command %s\n", command)
	result = cli.RunResult(t, command)
	fmt.Println(result.Output())
	//verification
	command = "show interfaces counters"
	fmt.Printf("Executing command %s\n", command)
	result = cli.RunResult(t, command)
	fmt.Println(result.Output())
	// T3
	command = "request mpls traffic-engineering pcep revoke"
	fmt.Printf("Executing command %s\n", command)
	result = cli.RunResult(t, command)
	fmt.Println(result.Output())
	// T4
	command = "request system tech-support TC5X component routing after 2024-11-03T00:00:00 force"
	fmt.Printf("Executing command %s\n", command)
	result = cli.RunResult(t, command)
	fmt.Println(result.Output())
	//verification
	command = "show system tech-support status"
	fmt.Printf("Executing command %s\n", command)
	result = cli.RunResult(t, command)
	fmt.Println(result.Output())
	// T5
	command = "show interfaces bundle-60011"
	fmt.Printf("Executing command %s\n", command)
	result = cli.RunResult(t, command)
	fmt.Println(result.Output())
	// set cli timestamp
	command = "set cli-timestamp"
	fmt.Printf("Executing command %s\n", command)
	result = cli.RunResult(t, command)
	fmt.Println(result.Output())
	// T6
	command = "show system ntp detail"
	fmt.Printf("Executing command %s\n", command)
	result = cli.RunResult(t, command)
	fmt.Println(result.Output())
	// T7
	command = "show dnos-internal oper-data /drivenets-top/interfaces/interface[name='bundle-60011.1']/oper-items/counters"
	fmt.Printf("Executing command %s\n", command)
	result = cli.RunResult(t, command)
	fmt.Println(result.Output())
	// T8
	command = "run ping 10.0.2.0 count 5"
	fmt.Printf("Executing command %s\n", command)
	result = cli.RunResult(t, command)
	fmt.Println(result.Output())
	// T9
	command = "show route protocol bgp"
	fmt.Printf("Executing command %s\n", command)
	result = cli.RunResult(t, command)
	fmt.Println(result.Output())
	// T10
	command = "show config"
	fmt.Printf("Executing command %s\n", command)
	result = cli.RunResult(t, command)
	fmt.Println(result.Output())
	// T11
	command = "show file log routing_engine/system-events.log"
	fmt.Printf("Executing command %s\n", command)
	result = cli.RunResult(t, command)
	fmt.Println(result.Output())
	// T12
	command = "show interfaces bundle-60011 ; show interfaces bundle-60008"
	fmt.Printf("Executing command %s\n", command)
	result = cli.RunResult(t, command)
	fmt.Println(result.Output())
	// T13
	command = "show system ntp detail ; show system hardware | no-more"
	fmt.Printf("Executing command %s\n", command)
	result = cli.RunResult(t, command)
	fmt.Println(result.Output())
	// T17 - failed CLI command
	command = "show bgp nei 1.1.1"
	fmt.Printf("Executing command %s\n", command)
	result = cli.RunResult(t, command)
	fmt.Println(result.Output())
	fmt.Println(result.Error())
	// T14
	command = "sho rou protocol bg | n"
	fmt.Printf("Executing command %s\n", command)
	result = cli.RunResult(t, command)
	fmt.Println(result.Output())
	// T15
	command = "sh int cou | n"
	fmt.Printf("Executing command %s\n", command)
	result = cli.RunResult(t, command)
	fmt.Println(result.Output())
	// T16
	command = "cl is n "
	fmt.Printf("Executing command %s\n", command)
	result = cli.RunResult(t, command)
	fmt.Println(result.Output())

}

func Test_TC_5_7_1(t *testing.T) {
	// execute CLI command exit
	cli := ondatra.DUT(t, "dal00").CLI()
	command := "exit"
	fmt.Printf("Executing command %s\n", command)
	result := cli.RunResult(t, command)
	fmt.Println(result.Output())
	fmt.Println("++++++++++++++++++++++++++++++")
	fmt.Println(result.Error())
}

func Test_TC_5_7_2(t *testing.T) {
	// execute CLI command quit
	cli := ondatra.DUT(t, "dal00").CLI()
	command := "quit"
	fmt.Printf("Executing command %s\n", command)
	result := cli.RunResult(t, command)
	fmt.Println(result.Output())
	fmt.Println("++++++++++++++++++++++++++++++")
	fmt.Println(result.Error())
}

func Test_TC_6_1(t *testing.T) {
	// set system name to 1 character
	cli := ondatra.DUT(t, "dal00").CLI()
	command := "show system name"
	fmt.Printf("Executing command %s\n", command)
	result := cli.RunResult(t, command)
	fmt.Println(result.Output())
	fmt.Println("++++++++++++++++++++++++++++++++++")
	fmt.Println(result.Error())
	//
	dut := ondatra.DUT(t, "dal00")
	// updates DUT config with replace config below
	dut.Config().New().
		WithDrivenetsText(
			`protocols bgp 1
				 `).
		Append(t)
	//
	command = "show system name"
	fmt.Printf("Executing command %s\n", command)
	result = cli.RunResult(t, command)
	fmt.Println(result.Output())
	fmt.Println("++++++++++++++++++++++++++++++++++")
	fmt.Println(result.Error())
	//
	command = "show system | no-more"
	fmt.Printf("Executing command %s\n", command)
	result = cli.RunResult(t, command)
	fmt.Println(result.Output())
	fmt.Println("++++++++++++++++++++++++++++++++++")
	fmt.Println(result.Error())
	//

}

func Test_TC_6_2(t *testing.T) {
	// set system name to 1 character
	cli := ondatra.DUT(t, "dal00").CLI()
	command := "show system name"
	fmt.Printf("Executing command %s\n", command)
	result := cli.RunResult(t, command)
	fmt.Println(result.Output())
	fmt.Println("++++++++++++++++++++++++++++++++++")
	fmt.Println(result.Error())
	//
	dut := ondatra.DUT(t, "dal00")
	// updates DUT config with replace config below
	dut.Config().New().
		WithDrivenetsText(
			`protocols bgp 65000
				 `).
		Append(t)
	//
	command = "show system name"
	fmt.Printf("Executing command %s\n", command)
	result = cli.RunResult(t, command)
	fmt.Println(result.Output())
	fmt.Println("++++++++++++++++++++++++++++++++++")
	fmt.Println(result.Error())
	//
	command = "show system | no-more"
	fmt.Printf("Executing command %s\n", command)
	result = cli.RunResult(t, command)
	fmt.Println(result.Output())
	fmt.Println("++++++++++++++++++++++++++++++++++")
	fmt.Println(result.Error())
	//

}

func Test_TC_7_1(t *testing.T) {
	// TC 7.1
	dut := ondatra.DUT(t, "dal00")
	// updates DUT config with single configuration
	dut.Config().New().
		WithDrivenetsText(
			`interfaces
	lo9001
	admin-state enabled
	description fofo
	!
	!
				`).
		Append(t)
}

func Test_TC_7_2(t *testing.T) {
	// TC 7.2
	dut := ondatra.DUT(t, "dal00")
	// updates DUT config with replace config below
	dut.Config().New().
		WithDrivenetsText(
			`system
	in-band-management source-interface lo0
	timezone UTC
	aaa-server
	admin-state enabled
	tacacs
	server priority 1 address 100.64.2.252
	accounting enabled
	authentication enabled
	authorization enabled
	password enc-gAAAAABlpNp2XeC6yfok8VmyeB0s4LfD1T98KX1nand2vcKzQdD8Clvmoh4rIVDrxq0VL_XAWGgaxJme-Kb3sPJesZhKVaCPgg==
	port 4949
	vrf mgmt0
	!
	timers
	hold-down 500
	!
	!
	!
	alarms
	admin-state enabled
	!
	cprl
	icmp
	burst 2
	rate 2
	!
	!
	dns
	domain-name dev.drivenets.com
	server priority 1 ip-address 100.64.0.1
	vrf mgmt0
	!
	server priority 2 ip-address 100.64.15.2
	vrf mgmt0
	!
	!
	ftp
	server
	admin-state disabled
	!
	!
	grpc
	max-sessions 6
	security
	tls server-certificate grpc_dnor
	!
	vrf default
	admin-state disabled
	!
	vrf mgmt0
	admin-state enabled
	!
	!
	high-availability
	disable-recovery-mode
	!
	info
	location WOPA_R410
	!
	logging
	syslog
	event-group all severity info
	event-group aaa severity info
	event-group bfd severity info
	event-group bgp severity info
	event-group clock severity info
	event-group diagnostics severity info
	event-group efm-oam severity info
	event-group fib-manager severity info
	event-group interfaces severity info
	event-group isis severity info
	event-group l2vpn severity info
	event-group lacp severity info
	event-group ldp severity info
	event-group lldp severity info
	event-group management severity info
	event-group monitoring severity info
	event-group mpls severity info
	event-group mpls-oam severity info
	event-group msdp severity info
	event-group multicast severity info
	event-group nat severity info
	event-group ospf severity info
	event-group ospfv3 severity info
	event-group pcep severity info
	event-group pim severity info
	event-group platform severity info
	event-group qos severity info
	event-group rib severity info
	event-group rsvp severity info
	event-group segment-routing severity info
	event-group services severity info
	event-group static-route severity info
	event-group system severity info
	event-group tcp severity info
	event-group vrrp severity info
	source-interface lo0
	suppress-event-list BGP_IPV4_NEIGHBOR_ADJACENCY_BACKWARD_CHANGE,BGP_IPV6_NEIGHBOR_ADJACENCY_BACKWARD_CHANGE
	server 10.10.75.158 vrf mgmt0
	admin-state enabled
	port 514
	protocol udp
	severity warning
	!
	server 10.144.28.196 vrf mgmt0
	!
	server 10.253.4.113 vrf mgmt0
	!
	server 10.253.16.232 vrf mgmt0
	!
	server 10.255.250.35 vrf mgmt0
	!
	server 100.64.2.252 vrf mgmt0
	admin-state enabled
	port 1234
	protocol tcp
	severity info
	!
	!
	!
	login
	ncm
	user dnroot
	password enc-gAAAAABgUIs_1U7_1nOEseJ75fkb7UVQSFOD5hi_wyLlHHb41toEsyUgw15QaCqpU2RYwe0CDuNkp_oLw27NNTC1oOk9jWRJXQ==
	role admin
	!
	!
	user dnroot
	password $6$m7b5lYW.PtBTZ0lr$ibN6YvwGYpPauDUIbNu7XMYRkKOnG3H.gyCQVweYv/lgsXMj65ClO8yeGTYX/Rmy80Bx23AeXOWM47nDXr9sN.
	!
	user dntechsupport
	password $6$m0NpQOgH8CusogMK$7JYzbJt55BxEOTkCG6z.9DNTPhDEbX8TN4Po1/mUJzlu/KaywxSoIRLvNYDdfE/WlVZVS2nfXd5J6S2Vin9D3/
	!
	ipmi
	user admin
	password enc-gAAAAABeANLvFw0fQC2rmQLMAF4UtTDyWIjWR1mDWrZmyNoLByecxKV7XxqJVBwoZUB0Y_3NmicHiFlT__YZX_48pt2ghenEnA==
	!
	user dnroot
	password enc-gAAAAABeANIYYvDD_EFSSRwAsxfNj9gHnUlqQnzhE4-GvvCjRuEb21-WA1LQk3V4YvQPabuWjwnrje9jwWNYsKj10czQ7G7-8A==
	!
	!
	!
	ncp 3
	model NCP-40C hw-model S9700-53DX
	admin-state enabled
	description dn-ncp-3
	hardware
	usb 0
	!
	!
	!
	netconf
	max-sessions 12
	port 830
	session-timeout 60
	vrf default
	admin-state enabled
	!
	vrf mgmt0
	admin-state enabled
	!
	!
	ntp
	server 100.64.0.2 vrf mgmt0
	admin-state enabled
	!
	!
	snmp
	community 32946t8pyrehfdlb vrf mgmt0
	!
	!
	ssh
	server
	vrf mgmt0
	client-list type deny
	client-list 100.64.2.252/32
	!
	!
	!
	telnet
	server
	admin-state disabled
	!
	!
	!
	
	services
	performance-monitoring
	profiles
	endpoint-delay STW_P1
	test-duration probes probe-count 10 probe-interval 1 repeat-interval 60
	!
	!
	simple-twamp session DAL00_DAL01
	admin-state enabled
	destination-address 20.0.2.0
	profile STW_P1
	source-address 20.0.1.0
	!
	!
	simple-twamp
	session-reflector
	admin-state enabled
	!
	session-sender
	admin-state enabled
	!
	!
	twamp
	admin-state enabled
	!
	!
	
	interfaces
	bundle-111
	admin-state enabled
	description bundle_with_member
	mtu 9000
	!
	bundle-666
	admin-state enabled
	!
	bundle-60002
	admin-state enabled
	description SA-5
	mtu 9222
	qos policy INGRESS direction in
	qos policy EGRESS direction out
	!
	bundle-60002.1
	admin-state enabled
	description "TOPO-1 MPLS PEERING"
	ipv4-address 10.1.5.1/30
	qos ip-marking trusted
	qos policy INGRESS direction in
	qos policy EGRESS direction out
	vlan-id 1
	!
	bundle-60002.2
	admin-state enabled
	description "TOPO-1 INTERNET PEERING"
	ipv4-address 10.1.5.5/30
	ipv6-address 10:1:5::5/126
	qos ip-marking trusted
	qos policy INGRESS direction in
	qos policy EGRESS direction out
	vlan-id 2
	!
	bundle-60007
	admin-state enabled
	description "IXIA PORT_3_2"
	mtu 9222
	qos policy INGRESS direction in
	qos policy EGRESS direction out
	!
	bundle-60007.3
	admin-state enabled
	ipv4-address 20.1.26.0/31
	ipv6-address 20:1:26::/127
	qos ip-marking trusted
	qos policy INGRESS direction in
	qos policy EGRESS direction out
	vlan-tags outer-tag 7 inner-tag 3
	!
	bundle-60007.10
	admin-state enabled
	ipv4-address 10.1.20.1/24
	qos ip-marking trusted
	qos policy INGRESS direction in
	qos policy EGRESS direction out
	vlan-tags outer-tag 7 inner-tag 20
	!
	bundle-60008
	admin-state enabled
	description ATL00
	mtu 9222
	qos policy INGRESS direction in
	qos policy EGRESS direction out
	!
	bundle-60008.1
	admin-state enabled
	ipv4-address 10.1.3.1/30
	mpls enabled
	qos ip-marking trusted
	qos policy INGRESS direction in
	qos policy EGRESS direction out
	vlan-id 1
	!
	bundle-60008.3
	admin-state enabled
	ipv4-address 20.1.3.0/31
	ipv6-address 20:1:3::/127
	mpls enabled
	qos ip-marking trusted
	qos policy INGRESS direction in
	qos policy EGRESS direction out
	vlan-id 3
	!
	bundle-60009
	admin-state enabled
	description ATL01
	mtu 9222
	qos policy INGRESS direction in
	qos policy EGRESS direction out
	!
	bundle-60009.1
	admin-state enabled
	ipv4-address 10.1.4.1/30
	mpls enabled
	qos ip-marking trusted
	qos policy INGRESS direction in
	qos policy EGRESS direction out
	vlan-id 1
	!
	bundle-60009.3
	admin-state enabled
	ipv4-address 20.1.4.0/31
	ipv6-address 20:1:4::/127
	mpls enabled
	qos ip-marking trusted
	qos policy INGRESS direction in
	qos policy EGRESS direction out
	vlan-id 3
	!
	bundle-60010
	admin-state enabled
	description "IXIA PORT_3_1"
	mtu 9222
	qos policy INGRESS direction in
	qos policy EGRESS direction out
	!
	bundle-60010.3
	admin-state enabled
	ipv4-address 20.1.20.0/31
	ipv6-address 20:1:20::/127
	mpls enabled
	qos ip-marking trusted
	qos policy INGRESS direction in
	qos policy EGRESS direction out
	vlan-tags outer-tag 3 inner-tag 3
	!
	bundle-60010.10
	admin-state enabled
	ipv4-address 10.1.10.1/24
	qos ip-marking trusted
	qos policy INGRESS direction in
	qos policy EGRESS direction out
	vlan-tags outer-tag 3 inner-tag 10
	!
	bundle-60010.100
	admin-state enabled
	ipv4-address 31.100.0.1/30
	vlan-tags outer-tag 3 inner-tag 100
	!
	bundle-60010.101
	admin-state enabled
	ipv4-address 31.101.0.1/30
	vlan-tags outer-tag 3 inner-tag 101
	!
	bundle-60010.102
	admin-state enabled
	ipv4-address 31.102.0.1/30
	vlan-tags outer-tag 3 inner-tag 102
	!
	bundle-60010.103
	admin-state enabled
	ipv4-address 31.103.0.1/30
	vlan-tags outer-tag 3 inner-tag 103
	!
	bundle-60010.104
	admin-state enabled
	ipv4-address 31.104.0.1/30
	vlan-tags outer-tag 3 inner-tag 104
	!
	bundle-60010.105
	admin-state enabled
	ipv4-address 31.105.0.1/30
	vlan-tags outer-tag 3 inner-tag 105
	!
	bundle-60010.106
	admin-state enabled
	ipv4-address 31.106.0.1/30
	vlan-tags outer-tag 3 inner-tag 106
	!
	bundle-60010.107
	admin-state enabled
	ipv4-address 31.107.0.1/30
	vlan-tags outer-tag 3 inner-tag 107
	!
	bundle-60010.108
	admin-state enabled
	ipv4-address 31.108.0.1/30
	vlan-tags outer-tag 3 inner-tag 108
	!
	bundle-60010.109
	admin-state enabled
	ipv4-address 31.109.0.1/30
	vlan-tags outer-tag 3 inner-tag 109
	!
	bundle-60010.110
	admin-state enabled
	ipv4-address 31.110.0.1/30
	vlan-tags outer-tag 3 inner-tag 110
	!
	bundle-60010.111
	admin-state enabled
	ipv4-address 31.111.0.1/30
	vlan-tags outer-tag 3 inner-tag 111
	!
	bundle-60010.112
	admin-state enabled
	ipv4-address 31.112.0.1/30
	vlan-tags outer-tag 3 inner-tag 112
	!
	bundle-60010.113
	admin-state enabled
	ipv4-address 31.113.0.1/30
	vlan-tags outer-tag 3 inner-tag 113
	!
	bundle-60010.114
	admin-state enabled
	ipv4-address 31.114.0.1/30
	vlan-tags outer-tag 3 inner-tag 114
	!
	bundle-60010.115
	admin-state enabled
	ipv4-address 31.115.0.1/30
	vlan-tags outer-tag 3 inner-tag 115
	!
	bundle-60010.116
	admin-state enabled
	ipv4-address 31.116.0.1/30
	vlan-tags outer-tag 3 inner-tag 116
	!
	bundle-60010.117
	admin-state enabled
	ipv4-address 31.117.0.1/30
	vlan-tags outer-tag 3 inner-tag 117
	!
	bundle-60010.118
	admin-state enabled
	ipv4-address 31.118.0.1/30
	vlan-tags outer-tag 3 inner-tag 118
	!
	bundle-60010.119
	admin-state enabled
	ipv4-address 31.119.0.1/30
	vlan-tags outer-tag 3 inner-tag 119
	!
	bundle-60010.120
	admin-state enabled
	ipv4-address 31.120.0.1/30
	vlan-tags outer-tag 3 inner-tag 120
	!
	bundle-60010.121
	admin-state enabled
	ipv4-address 31.121.0.1/30
	vlan-tags outer-tag 3 inner-tag 121
	!
	bundle-60010.122
	admin-state enabled
	ipv4-address 31.122.0.1/30
	vlan-tags outer-tag 3 inner-tag 122
	!
	bundle-60010.123
	admin-state enabled
	ipv4-address 31.123.0.1/30
	vlan-tags outer-tag 3 inner-tag 123
	!
	bundle-60010.124
	admin-state enabled
	ipv4-address 31.124.0.1/30
	vlan-tags outer-tag 3 inner-tag 124
	!
	bundle-60010.125
	admin-state enabled
	ipv4-address 31.125.0.1/30
	vlan-tags outer-tag 3 inner-tag 125
	!
	bundle-60010.126
	admin-state enabled
	ipv4-address 31.126.0.1/30
	vlan-tags outer-tag 3 inner-tag 126
	!
	bundle-60010.127
	admin-state enabled
	ipv4-address 31.127.0.1/30
	vlan-tags outer-tag 3 inner-tag 127
	!
	bundle-60010.128
	admin-state enabled
	ipv4-address 31.128.0.1/30
	vlan-tags outer-tag 3 inner-tag 128
	!
	bundle-60010.129
	admin-state enabled
	ipv4-address 31.129.0.1/30
	vlan-tags outer-tag 3 inner-tag 129
	!
	bundle-60010.130
	admin-state enabled
	ipv4-address 31.130.0.1/30
	vlan-tags outer-tag 3 inner-tag 130
	!
	bundle-60010.131
	admin-state enabled
	ipv4-address 31.131.0.1/30
	vlan-tags outer-tag 3 inner-tag 131
	!
	bundle-60010.132
	admin-state enabled
	ipv4-address 31.132.0.1/30
	vlan-tags outer-tag 3 inner-tag 132
	!
	bundle-60010.133
	admin-state enabled
	ipv4-address 31.133.0.1/30
	vlan-tags outer-tag 3 inner-tag 133
	!
	bundle-60010.134
	admin-state enabled
	ipv4-address 31.134.0.1/30
	vlan-tags outer-tag 3 inner-tag 134
	!
	bundle-60010.135
	admin-state enabled
	ipv4-address 31.135.0.1/30
	vlan-tags outer-tag 3 inner-tag 135
	!
	bundle-60010.136
	admin-state enabled
	ipv4-address 31.136.0.1/30
	vlan-tags outer-tag 3 inner-tag 136
	!
	bundle-60010.137
	admin-state enabled
	ipv4-address 31.137.0.1/30
	vlan-tags outer-tag 3 inner-tag 137
	!
	bundle-60010.138
	admin-state enabled
	ipv4-address 31.138.0.1/30
	vlan-tags outer-tag 3 inner-tag 138
	!
	bundle-60010.139
	admin-state enabled
	ipv4-address 31.139.0.1/30
	vlan-tags outer-tag 3 inner-tag 139
	!
	bundle-60010.140
	admin-state enabled
	ipv4-address 31.140.0.1/30
	vlan-tags outer-tag 3 inner-tag 140
	!
	bundle-60010.141
	admin-state enabled
	ipv4-address 31.141.0.1/30
	vlan-tags outer-tag 3 inner-tag 141
	!
	bundle-60010.142
	admin-state enabled
	ipv4-address 31.142.0.1/30
	vlan-tags outer-tag 3 inner-tag 142
	!
	bundle-60010.143
	admin-state enabled
	ipv4-address 31.143.0.1/30
	vlan-tags outer-tag 3 inner-tag 143
	!
	bundle-60010.144
	admin-state enabled
	ipv4-address 31.144.0.1/30
	vlan-tags outer-tag 3 inner-tag 144
	!
	bundle-60010.145
	admin-state enabled
	ipv4-address 31.145.0.1/30
	vlan-tags outer-tag 3 inner-tag 145
	!
	bundle-60010.146
	admin-state enabled
	ipv4-address 31.146.0.1/30
	vlan-tags outer-tag 3 inner-tag 146
	!
	bundle-60010.147
	admin-state enabled
	ipv4-address 31.147.0.1/30
	vlan-tags outer-tag 3 inner-tag 147
	!
	bundle-60010.148
	admin-state enabled
	ipv4-address 31.148.0.1/30
	vlan-tags outer-tag 3 inner-tag 148
	!
	bundle-60010.149
	admin-state enabled
	ipv4-address 31.149.0.1/30
	vlan-tags outer-tag 3 inner-tag 149
	!
	bundle-60010.150
	admin-state enabled
	ipv4-address 31.150.0.1/30
	vlan-tags outer-tag 3 inner-tag 150
	!
	bundle-60011
	admin-state enabled
	description DAL00
	ipv4-address 192.168.0.1/30
	mpls disabled
	mtu 9222
	qos policy INGRESS direction in
	qos policy EGRESS direction out
	!
	bundle-60011.1
	admin-state enabled
	ipv4-address 10.1.2.1/30
	mpls enabled
	qos ip-marking trusted
	qos policy INGRESS direction in
	qos policy EGRESS direction out
	vlan-id 1
	!
	bundle-60011.2
	admin-state enabled
	ipv4-address 10.1.2.5/30	qos ip-marking trusted
	qos policy INGRESS direction in
	qos policy EGRESS direction out
	vlan-id 2
	!
	bundle-60011.3
	admin-state enabled
	ipv4-address 20.1.2.0/31
	ipv6-address 20:1:2::/127
	mpls enabled
	qos ip-marking trusted
	qos policy INGRESS direction in
	qos policy EGRESS direction out
	vlan-id 3
	!
	bundle-60012
	admin-state enabled
	description SA-5
	mtu 9222
	qos policy INGRESS direction in
	qos policy EGRESS direction out
	!
	bundle-60012.1
	admin-state enabled
	ipv4-address 10.100.200.0/31
	ipv6-address 10:100:200::/127
	qos ip-marking trusted
	qos policy INGRESS direction in
	qos policy EGRESS direction out
	vlan-id 1
	!
	bundle-60012.2
	admin-state enabled
	ipv4-address 10.101.201.0/31
	ipv6-address 10:101:201::/127
	qos ip-marking trusted
	qos policy INGRESS direction in
	qos policy EGRESS direction out
	vlan-id 2
	!
	bundle-60012.3
	admin-state enabled
	ipv4-address 20.103.203.0/31
	ipv6-address 20:103:203::/127
	qos ip-marking trusted
	qos policy INGRESS direction in
	qos policy EGRESS direction out
	vlan-id 3
	!
	bundle-60038
	admin-state enabled
	description "TOPO-1 LOOP_BUNDLE_60038"
	ipv4-address 10.1.1.0/31
	mtu 9222
	qos policy INGRESS direction in
	qos policy EGRESS direction out
	!
	bundle-60038.100
	admin-state enabled
	ipv4-address 10.2.100.0/31
	ipv6-address 10:2:100::/127
	qos ip-marking trusted
	qos policy INGRESS direction in
	qos policy EGRESS direction out
	vlan-id 100
	!
	bundle-60038.101
	admin-state enabled
	ipv4-address 10.2.101.0/31
	ipv6-address 10:2:101::/127
	qos ip-marking trusted
	qos policy INGRESS direction in
	qos policy EGRESS direction out
	vlan-id 101
	!
	bundle-60038.999
	admin-state enabled
	ipv4-address 169.254.0.0/31
	qos ip-marking trusted
	qos policy INGRESS direction in
	qos policy EGRESS direction out
	vlan-id 999
	!
	bundle-60039
	admin-state enabled
	description "TOPO-1 LOOP_BUNDLE_60039"
	ipv4-address 10.1.1.1/31
	mtu 9222
	qos policy INGRESS direction in
	qos policy EGRESS direction out
	!
	bundle-60039.100
	admin-state enabled
	ipv4-address 10.2.100.1/31
	ipv6-address 10:2:100::1/127
	qos ip-marking trusted
	qos policy INGRESS direction in
	qos policy EGRESS direction out
	vlan-id 100
	!
	bundle-60039.101
	admin-state enabled
	ipv4-address 10.2.101.1/31
	ipv6-address 10:2:101::1/127
	qos ip-marking trusted
	qos policy INGRESS direction in
	qos policy EGRESS direction out
	vlan-id 101
	!
	bundle-60039.999
	admin-state enabled
	ipv4-address 169.254.0.1/31
	qos ip-marking trusted
	qos policy INGRESS direction in
	qos policy EGRESS direction out
	vlan-id 999
	!
	console-ncp-3/0
	admin-state enabled
	!
	ctrl-ncp-3/0
	admin-state enabled
	!
	ctrl-ncp-3/1
	admin-state enabled
	!
	fab-ncp400-3/0/0
	admin-state enabled
	!
	fab-ncp400-3/0/1
	admin-state enabled
	!
	fab-ncp400-3/0/2
	admin-state enabled
	!
	fab-ncp400-3/0/3
	admin-state enabled
	!
	fab-ncp400-3/0/4
	admin-state enabled
	!
	fab-ncp400-3/0/5
	admin-state enabled
	!
	fab-ncp400-3/0/6
	admin-state enabled
	!
	fab-ncp400-3/0/7
	admin-state enabled
	!
	fab-ncp400-3/0/8
	admin-state enabled
	!
	fab-ncp400-3/0/9
	admin-state enabled
	!
	fab-ncp400-3/0/10
	admin-state enabled
	!
	fab-ncp400-3/0/11
	admin-state enabled
	!
	fab-ncp400-3/0/12
	admin-state enabled
	!
	ge100-3/0/0
	admin-state enabled
	bundle-id 60010
	fec none
	!
	ge100-3/0/1
	admin-state enabled
	bundle-id 60007
	fec none
	!
	ge100-3/0/2
	admin-state enabled
	fec none
	!
	ge100-3/0/3
	admin-state enabled
	fec none
	!
	ge100-3/0/4
	admin-state enabled
	bundle-id 60008
	fec rs-fec-528-514
	!
	ge100-3/0/5
	admin-state enabled
	bundle-id 60009
	fec none
	!
	ge100-3/0/6
	admin-state enabled
	bundle-id 60002
	fec none
	!
	ge100-3/0/7
	admin-state enabled
	bundle-id 60012
	fec none
	!
	ge100-3/0/8
	admin-state enabled
	fec none
	!
	ge100-3/0/9
	admin-state enabled
	bundle-id 60011
	fec none
	!
	ge100-3/0/10
	admin-state enabled
	fec none
	!
	ge100-3/0/11
	admin-state enabled
	bundle-id 60011
	fec none
	!
	ge100-3/0/12
	admin-state enabled
	fec none
	!
	ge100-3/0/13
	admin-state enabled
	fec none
	!
	ge100-3/0/14
	admin-state enabled
	fec none
	!
	ge100-3/0/15
	admin-state enabled
	fec none
	!
	ge100-3/0/16
	admin-state enabled
	fec none
	!
	ge100-3/0/17
	admin-state enabled
	fec none
	!
	ge100-3/0/18
	admin-state enabled
	fec none
	!
	ge100-3/0/19
	admin-state enabled
	fec none
	!
	ge100-3/0/20
	admin-state enabled
	fec none
	!
	ge100-3/0/21
	admin-state enabled
	fec none
	!
	ge100-3/0/22
	admin-state enabled
	fec none
	!
	ge100-3/0/23
	admin-state enabled
	fec none
	!
	ge100-3/0/24
	admin-state enabled
	fec none
	!
	ge100-3/0/25
	admin-state enabled
	fec none
	!
	ge100-3/0/26
	admin-state enabled
	fec none
	!
	ge100-3/0/27
	admin-state enabled
	fec none
	!
	ge100-3/0/28
	admin-state disabled
	fec none
	util-rate-threshold 100
	!
	ge100-3/0/29
	admin-state enabled
	fec none
	!
	ge100-3/0/30
	admin-state enabled
	fec none
	!
	ge100-3/0/31
	admin-state enabled
	fec none
	!
	ge100-3/0/32
	admin-state enabled
	fec none
	!
	ge100-3/0/33
	admin-state enabled
	fec none
	!
	ge100-3/0/34
	admin-state enabled
	fec none
	!
	ge100-3/0/35
	admin-state enabled
	fec none
	!
	ge100-3/0/36
	admin-state enabled
	fec none
	!
	ge100-3/0/37
	admin-state enabled
	fec none
	!
	ge100-3/0/38
	admin-state enabled
	bundle-id 60038
	fec none
	dampening
	admin-state enabled
	!
	!
	ge100-3/0/39
	admin-state disabled
	bundle-id 60039
	fec none
	!
	ipmi-ncp-0/0
	admin-state enabled
	!
	ipmi-ncp-3/0
	admin-state enabled
	!
	lo0
	admin-state enabled
	description TOPO-1
	ipv4-address 10.0.1.0/32
	!
	lo2
	admin-state enabled
	description TOPO-2
	ipv4-address 20.0.1.0/32
	ipv6-address 20:0:1::/128
	!
	lo10
	admin-state enabled
	ipv4-address 100.0.0.1/32
	ipv6-address 100::1/128
	!
	lo900
	admin-state enabled
	description new_loopback
	!
	lo999
	admin-state enabled
	ipv4-address 9.9.9.9/32
	!
	lo1000
	admin-state enabled
	ipv4-address 1.1.1.1/32
	!
	lo12345
	admin-state enabled
	!
	mgmt0
	admin-state enabled
	!
	!
	
	routing-options
	maximum-routes 100000 threshold 20
	!
	
	protocols
	bgp nsr disabled
	bfd
	interface bundle-60011
	local-address 192.168.0.1
	neighbor 192.168.0.2
	!
	interface bundle-60038
	local-address 10.1.1.0
	neighbor 10.1.1.1
	!
	interface bundle-60039
	local-address 10.1.1.1
	neighbor 10.1.1.0
	!
	!
	bgp 65000
	network import-check disabled
	route-reflection policy-out attribute-change enabled
	router-id 10.0.1.0
	address-family ipv4-vpn
	fast-reroute enabled
	label-allocation per-nexthop
	!
	address-family ipv6-vpn
	fast-reroute enabled
	label-allocation per-nexthop
	!
	neighbor 20.1.20.1
	remote-as 65000
	update-source bundle-60010.3
	address-family ipv4-unicast
	labeled-unicast
	!
	!
	neighbor 31.100.0.2
	remote-as 65000
	admin-state disabled
	update-source bundle-60010.100
	bfd
	admin-state enabled
	min-rx 500
	min-tx 500
	multiplier 5
	!
	address-family ipv4-unicast
	!
	!
	neighbor 31.101.0.2
	remote-as 65000
	admin-state disabled
	update-source bundle-60010.101
	bfd
	admin-state enabled
	min-rx 500
	min-tx 500
	multiplier 5
	!
	address-family ipv4-unicast
	!
	!
	neighbor 31.102.0.2
	remote-as 65000
	admin-state disabled
	update-source bundle-60010.102
	bfd
	admin-state enabled
	min-rx 500
	min-tx 500
	multiplier 5
	!
	address-family ipv4-unicast
	!
	!
	neighbor 31.103.0.2
	remote-as 65000
	admin-state disabled
	update-source bundle-60010.103
	bfd
	admin-state enabled
	min-rx 500
	min-tx 500
	multiplier 5
	!
	address-family ipv4-unicast
	!
	!
	neighbor 31.104.0.2
	remote-as 65000
	admin-state disabled
	update-source bundle-60010.104
	bfd
	admin-state enabled
	min-rx 500
	min-tx 500
	multiplier 5
	!
	address-family ipv4-unicast
	!
	!
	neighbor 31.105.0.2
	remote-as 65000
	admin-state disabled
	update-source bundle-60010.105
	bfd
	admin-state enabled
	min-rx 500
	min-tx 500
	multiplier 5
	!
	address-family ipv4-unicast
	!
	!
	neighbor 31.106.0.2
	remote-as 65000
	admin-state disabled
	update-source bundle-60010.106
	bfd
	admin-state enabled
	min-rx 500
	min-tx 500
	multiplier 5
	!
	address-family ipv4-unicast
	!
	!
	neighbor 31.107.0.2
	remote-as 65000
	admin-state disabled
	update-source bundle-60010.107
	bfd
	admin-state enabled
	min-rx 500
	min-tx 500
	multiplier 5
	!
	address-family ipv4-unicast
	!
	!
	neighbor 31.108.0.2
	remote-as 65000
	admin-state disabled
	update-source bundle-60010.108
	bfd
	admin-state enabled
	min-rx 500
	min-tx 500
	multiplier 5
	!
	address-family ipv4-unicast
	!
	!
	neighbor 31.109.0.2
	remote-as 65000
	admin-state disabled
	update-source bundle-60010.109
	bfd
	admin-state enabled
	min-rx 500
	min-tx 500
	multiplier 5
	!
	address-family ipv4-unicast
	!
	!
	neighbor 31.110.0.2
	remote-as 65000
	admin-state disabled
	update-source bundle-60010.110
	bfd
	admin-state enabled
	min-rx 500
	min-tx 500
	multiplier 5
	!
	address-family ipv4-unicast
	!
	!
	neighbor 31.111.0.2
	remote-as 65000
	admin-state disabled
	update-source bundle-60010.111
	bfd
	admin-state enabled
	min-rx 500
	min-tx 500
	multiplier 5
	!
	address-family ipv4-unicast
	!
	!
	neighbor 31.112.0.2
	remote-as 65000
	admin-state disabled
	update-source bundle-60010.112
	bfd
	admin-state enabled
	min-rx 500
	min-tx 500
	multiplier 5
	!
	address-family ipv4-unicast
	!
	!
	neighbor 31.113.0.2
	remote-as 65000
	admin-state disabled
	update-source bundle-60010.113
	bfd
	admin-state enabled
	min-rx 500
	min-tx 500
	multiplier 5
	!
	address-family ipv4-unicast
	!
	!
	neighbor 31.114.0.2
	remote-as 65000
	admin-state disabled
	update-source bundle-60010.114
	bfd
	admin-state enabled
	min-rx 500
	min-tx 500
	multiplier 5
	!
	address-family ipv4-unicast
	!
	!
	neighbor 31.115.0.2
	remote-as 65000
	admin-state disabled
	update-source bundle-60010.115
	bfd
	admin-state enabled
	min-rx 500
	min-tx 500
	multiplier 5
	!
	address-family ipv4-unicast
	!
	!
	neighbor 31.116.0.2
	remote-as 65000
	admin-state disabled
	update-source bundle-60010.116
	bfd
	admin-state enabled
	min-rx 500
	min-tx 500
	multiplier 5
	!
	address-family ipv4-unicast
	!
	!
	neighbor 31.117.0.2
	remote-as 65000
	admin-state disabled
	update-source bundle-60010.117
	bfd
	admin-state enabled
	min-rx 500
	min-tx 500
	multiplier 5
	!
	address-family ipv4-unicast
	!
	!
	neighbor 31.118.0.2
	remote-as 65000
	admin-state disabled
	update-source bundle-60010.118
	bfd
	admin-state enabled
	min-rx 500
	min-tx 500
	multiplier 5
	!
	address-family ipv4-unicast
	!
	!
	neighbor 31.119.0.2
	remote-as 65000
	admin-state disabled
	update-source bundle-60010.119
	bfd
	admin-state enabled
	min-rx 500
	min-tx 500
	multiplier 5
	!
	address-family ipv4-unicast
	!
	!
	neighbor 31.120.0.2
	remote-as 65000
	admin-state disabled
	update-source bundle-60010.120
	bfd
	admin-state enabled
	min-rx 500
	min-tx 500
	multiplier 5
	!
	address-family ipv4-unicast
	!
	!
	neighbor 31.121.0.2
	remote-as 65000
	admin-state disabled
	update-source bundle-60010.121
	bfd
	admin-state enabled
	min-rx 500
	min-tx 500
	multiplier 5
	!
	address-family ipv4-unicast
	!
	!
	neighbor 31.122.0.2
	remote-as 65000
	admin-state disabled
	update-source bundle-60010.122
	bfd
	admin-state enabled
	min-rx 500
	min-tx 500
	multiplier 5
	!
	address-family ipv4-unicast
	!
	!
	neighbor 31.123.0.2
	remote-as 65000
	admin-state disabled
	update-source bundle-60010.123
	bfd
	admin-state enabled
	min-rx 500
	min-tx 500
	multiplier 5
	!
	address-family ipv4-unicast
	!
	!
	neighbor 31.124.0.2
	remote-as 65000
	admin-state disabled
	update-source bundle-60010.124
	bfd
	admin-state enabled
	min-rx 500
	min-tx 500
	multiplier 5
	!
	address-family ipv4-unicast
	!
	!
	neighbor 31.125.0.2
	remote-as 65000
	admin-state disabled
	update-source bundle-60010.125
	bfd
	admin-state enabled
	min-rx 500
	min-tx 500
	multiplier 5
	!
	address-family ipv4-unicast
	!
	!
	neighbor 31.126.0.2
	remote-as 65000
	admin-state disabled
	update-source bundle-60010.126
	bfd
	admin-state enabled
	min-rx 500
	min-tx 500
	multiplier 5
	!
	address-family ipv4-unicast
	!
	!
	neighbor 31.127.0.2
	remote-as 65000
	admin-state disabled
	update-source bundle-60010.127
	bfd
	admin-state enabled
	min-rx 500
	min-tx 500
	multiplier 5
	!
	address-family ipv4-unicast
	!
	!
	neighbor 31.128.0.2
	remote-as 65000
	admin-state disabled
	update-source bundle-60010.128
	bfd
	admin-state enabled
	min-rx 500
	min-tx 500
	multiplier 5
	!
	address-family ipv4-unicast
	!
	!
	neighbor 31.129.0.2
	remote-as 65000
	admin-state disabled
	update-source bundle-60010.129
	bfd
	admin-state enabled
	min-rx 500
	min-tx 500
	multiplier 5
	!
	address-family ipv4-unicast
	!
	!
	neighbor 31.130.0.2
	remote-as 65000
	admin-state disabled
	update-source bundle-60010.130
	bfd
	admin-state enabled
	min-rx 500
	min-tx 500
	multiplier 5
	!
	address-family ipv4-unicast
	!
	!
	neighbor 31.131.0.2
	remote-as 65000
	admin-state disabled
	update-source bundle-60010.131
	bfd
	admin-state enabled
	min-rx 500
	min-tx 500
	multiplier 5
	!
	address-family ipv4-unicast
	!
	!
	neighbor 31.132.0.2
	remote-as 65000
	admin-state disabled
	update-source bundle-60010.132
	bfd
	admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.133.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.133
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.134.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.134
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.135.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.135
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.136.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.136
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.137.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.137
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.138.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.138
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.139.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.139
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.140.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.140
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.141.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.141
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.142.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.142
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.143.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.143
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.144.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.144
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.145.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.145
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.146.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.146
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.147.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.147
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.148.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.148
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.149.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.149
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.150.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.150
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor-group NC_CLIENTS
 remote-as 65000
 update-source lo0
 address-family ipv4-rt-constrains
 route-reflector-client
 !
 address-family ipv4-vpn
 policy TOPO-1-OUT out
 route-reflector-client
 send-community
 send-large-community enabled
 !
 address-family ipv6-vpn
 policy TOPO-1-OUT out
 route-reflector-client
 send-community
 send-large-community enabled
 !
 neighbor 10.0.2.0
 description DAL01
 !
 neighbor 10.0.4.0
 description ATL01
 !
 !
 neighbor-group RR_FULL_MESH
 remote-as 65000
 update-source lo0
 address-family ipv4-rt-constrains
 !
 address-family ipv4-vpn
 policy TOPO-1-OUT out
 send-community
 send-large-community enabled
 !
 address-family ipv6-vpn
 policy TOPO-1-IN out
 send-community
 send-large-community enabled
 !
 neighbor 10.0.3.0
 description ATL00
 !
 !
 neighbor-group TOPO-2-EBGP
 remote-as 65500
 local-as 65499 type no-prepend-replace-as
 address-family ipv4-unicast
 dampening enabled
 !
 neighbor 20.1.2.1
 !
 neighbor 20.1.3.1
 !
 neighbor 20.1.4.1
 !
 !
 neighbor-group TOPO-2-FULL-MESH
 remote-as 65000
 update-source lo2
 address-family ipv4-vpn
 policy TOPO-2-OUT out
 send-community
 !
 address-family ipv6-vpn
 policy TOPO-2-OUT out
 send-community
 !
 neighbor 20.0.2.0
 !
 neighbor 20.0.3.0
 !
 neighbor 20.0.4.0
 !
 !
 !
 isis
 maximum-routes 300 threshold 50
 nsr enabled
 instance CORE
 iso-network 49.0000.0100.0000.1000.00
 level level-2
 lsp-mtu 1400
 max-metric 1000000
 address-family ipv4-unicast
 traffic-engineering enabled
 segment-routing
 admin-state enabled
 !
 !
 interface lo0
 address-family ipv4-unicast
 prefix-sid index 1 prefix-type node
 !
 !
 interface bundle-60007.10
 address-family ipv4-unicast
 !
 !
 interface bundle-60008.1
 network-type point-to-point
 address-family ipv4-unicast
 metric 5
 !
 bfd
 admin-state enabled
 min-rx 100
 min-tx 100
 multiplier 3
 !
 !
 interface bundle-60009.1
 network-type point-to-point
 address-family ipv4-unicast
 metric 100
 !
 bfd
 admin-state enabled
 min-rx 100
 min-tx 100
 multiplier 3
 !
 !
 interface bundle-60010.10
 address-family ipv4-unicast
 !
 !
 interface bundle-60010.100
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.101
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.102
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.103
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.104
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.105
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.106
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.107
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.108
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.109
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.110
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.111
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.112
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.113
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.114
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.115
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.116
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.117
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.118
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.119
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.120
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.121
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.122
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.123
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.124
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.125
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.126
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.127
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.128
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.129
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.130
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.131
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.132
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.133
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.134
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.135
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.136
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.137
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.138
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.139
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.140
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.141
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.142
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.143
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.144
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.145
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.146
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.147
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.148
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.149
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.150
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60011.1
 network-type point-to-point
 address-family ipv4-unicast
 metric 10
 !
 address-family ipv6-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 100
 min-tx 100
 multiplier 3
 !
 !
 interface bundle-60011.2
 network-type point-to-point
 address-family ipv4-unicast
 metric 800
 !
 !
 timers
 lsp-lifetime 65535
 lsp-refresh 65000
 !
 !
 instance TOPO-2
 administrative-distance 117
 iso-network 49.0000.0100.0000.1000.00
 level level-2
 lsp-mtu 1400
 address-family ipv4-unicast
 traffic-engineering enabled
 !
 interface lo2
 address-family ipv4-unicast
 !
 address-family ipv6-unicast
 !
 !
 interface bundle-60007.3
 address-family ipv4-unicast
 metric 10
 !
 !
 interface bundle-60008.3
 address-family ipv4-unicast
 metric 10
 !
 address-family ipv6-unicast
 metric 10
 !
 !
 interface bundle-60009.3
 address-family ipv4-unicast
 metric 100
 !
 address-family ipv6-unicast
 metric 100
 !
 bfd
 admin-state enabled
 !
 !
 interface bundle-60010.3
 address-family ipv4-unicast
 metric 10
 !
 address-family ipv6-unicast
 metric 10
 !
 !
 interface bundle-60011.3
 address-family ipv4-unicast
 metric 10
 !
 address-family ipv6-unicast
 metric 10
 !
 bfd
 admin-state enabled
 !
 !
 overload on-startup
 admin-state enabled
 advertisement-type max-metric
 interval 600
 wait-for-bgp bgp-delay 60
 !
 !
 !
 lacp
 interface bundle-60002
 mode active
 period short
 !
 interface bundle-60008
 mode active
 period short
 !
 interface bundle-60009
 mode active
 period short
 !
 interface bundle-60011
 mode active
 period short
 !
 interface bundle-60012
 mode active
 period short
 !
 !
 lldp
 interface ge100-3/0/2
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/3
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/4
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/5
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/6
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/7
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/8
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/9
 !
 interface ge100-3/0/10
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/11
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/12
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/13
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/14
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/15
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/16
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/17
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/18
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/19
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/20
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/21
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/22
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/23
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/24
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/25
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/26
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/27
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/29
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/32
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/33
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/34
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/35
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/36
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/37
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/38
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/39
 receive enabled
 transmit enabled
 !
 !
 mpls
 traffic-engineering
 router-id 20.0.1.0
 interface bundle-60008.3
 !
 interface bundle-60009.3
 !
 interface bundle-60010.3
 !
 interface bundle-60011.3
 !
 !
 !
 pim
 address-family ipv4
 interface bundle-60007.3
 admin-state enabled
 !
 interface bundle-60008.3
 admin-state enabled
 !
 interface bundle-60010.3
 admin-state enabled
 !
 interface bundle-60011.3
 admin-state enabled
 !
 !
 static-rp 20.0.22.0
 !
 !
 rsvp
 auto-mesh
 tunnel-template FULL_MESH_TOPO_2
 admin-state enabled
 destination-address TOPO-2-LO0
 igp-instance isis TOPO-2
 rib-unicast-install enabled
 source-address 20.0.1.0
 bfd
 admin-state enabled
 !
 primary
 protection node-protection
 path-options 1
 path dynamic
 !
 !
 !
 tunnel-template FULL_MESH_TOPO_2_IXIA
 admin-state enabled
 destination-address TOPO-2-IXIA
 igp-instance isis TOPO-2
 rib-unicast-install enabled
 source-address 20.0.1.0
 primary
 path-options 1
 path dynamic
 !
 !
 !
 !
 interface bundle-60008.3
 protection
 auto-bypass enabled
 !
 !
 interface bundle-60009.3
 protection
 auto-bypass enabled
 !
 !
 interface bundle-60010.3
 !
 interface bundle-60011.3
 protection
 auto-bypass enabled
 !
 !
 tunnel BBB bypass
 admin-state enabled
 destination-address 10.0.3.0
 primary
 path-options 1
 path dynamic
 !
 !
 !
 tunnel AAA
 admin-state enabled
 destination-address 20.0.2.0
 igp-instance isis TOPO-2
 primary
 path-options 1
 path dynamic
 !
 !
 !
 tunnel DAL00-ATL01
 admin-state disabled
 destination-address 20.0.4.0
 igp-instance isis TOPO-2
 source-address 20.0.1.0
 primary
 path-options 1
 path dynamic
 !
 !
 !
 !
 static
 maximum-routes 6 threshold 4
 !
 !
 
 multicast
 rpf-intact admin-state enabled
 !
 
 routing-policy
 community-list TOPO-1
 rule 1 allow value 65000:10
 !
 community-list TOPO-2
 rule 1 allow value 65000:20
 !
 policy TOPO-1-IN
 rule 1 allow
 set community 65000:10
 !
 !
 policy TOPO-1-OUT
 rule 10 allow
 match community TOPO-1
 !
 rule 100 allow
 !
 !
 policy TOPO-2-IN
 rule 1 allow
 set community 65000:20
 !
 !
 policy TOPO-2-OUT
 rule 10 allow
 match community TOPO-2
 !
 !
 prefix-list ipv4 TOPO-2-IXIA
 rule 1 deny 20.0.1.0/32,
 rule 2 deny 20.0.2.0/32,
 rule 3 deny 20.0.3.0/32,
 rule 4 deny 20.0.4.0/32,
 rule 5 allow 20.0.0.0/16 matching-len eq 32
 !
 prefix-list ipv4 TOPO-2-LO0
 rule 1 allow 20.0.2.0/32,
 rule 2 allow 20.0.4.0/32
 !
 !
 
 qos
 hw-mapping
 queue-size
 speed-ranges
 upto 50 mbps use 50 mbps
 upto 100 mbps use 100 mbps
 upto 250 mbps use 250 mbps
 upto 500 mbps use 500 mbps
 upto 750 mbps use 750 mbps
 upto 1 gbps use 1 gbps
 !
 !
 !
 wred-profile WRED-20K-DUAL
 curve green
 min 16 milliseconds max 20 milliseconds
 !
 curve yellow
 min 8 milliseconds max 12 milliseconds
 !
 !
 wred-profile WRED-25K-SINGLE
 curve green
 min 20 milliseconds max 25 milliseconds
 !
 !
 policy EGRESS
 rule 6
 match traffic-class CLASS-CONTROL
 action
 queue
 forwarding-class af
 bandwidth 4 percent
 size 5 milliseconds
 !
 !
 !
 !
 rule default
 action
 queue
 forwarding-class df
 bandwidth 30 percent
 size 25 milliseconds
 wred-profile WRED-25K-SINGLE
 !
 !
 !
 !
 !
 policy INGRESS
 rule 6
 match traffic-class CLASS6
 action
 set
 qos-tag 6
 !
 !
 !
 rule default
 !
 !
 traffic-class-map CLASS-CONTROL
 qos-tag 6
 !
 traffic-class-map CLASS6
 mpls-exp 6
 precedence 6
 !
 !
 
 network-services
 vrf
 instance DAL00_A
 interface bundle-60038.999
 protocols
 bgp 65000
 router-id 169.254.0.0
 neighbor 169.254.0.1
 remote-as 4290000002
 local-as 4290000001 type no-prepend-replace-as
 address-family ipv4-unicast
 !
 !
 !
 !
 !
 instance DAL00_B
 interface bundle-60039.999
 !
 instance INTERNET
 interface bundle-60002.2
 interface bundle-60038.100
 interface bundle-60038.101
 protocols
 bgp 65000
 route-distinguisher 10.0.1.0:2
 router-id 10.0.1.2
 address-family ipv4-unicast
 export-vpn route-target 65002:2
 fast-reroute enabled
 import-vpn route-target 65002:2
 !
 address-family ipv6-unicast
 export-vpn route-target 65002:2
 fast-reroute enabled
 import-vpn route-target 65002:2
 !
 neighbor 10.1.5.6
 remote-as 64512
 advertisement-interval 30
 address-family ipv4-unicast
 as-loop-check disabled
 policy TOPO-1-IN in
 !
 !
 neighbor 10.2.100.1
 remote-as 65100
 local-as 65002 type no-prepend-replace-as
 address-family ipv4-unicast
 policy TOPO-1-IN in
 !
 address-family ipv6-unicast
 policy TOPO-1-IN in
 !
 !
 neighbor 10.2.101.1
 remote-as 65101
 local-as 65002 type no-prepend-replace-as
 address-family ipv4-unicast
 policy TOPO-1-IN in
 !
 address-family ipv6-unicast
 policy TOPO-1-IN in
 !
 !
 neighbor 10:1:5::6
 remote-as 64512
 address-family ipv6-unicast
 policy TOPO-1-IN in
 !
 !
 !
 static
 address-family ipv4-unicast
 route 200.0.0.0/24
 next-hop 10.2.100.1
 !
 !
 !
 !
 !
 instance LEFT
 interface bundle-60038
 !
 instance RIGHT
 interface bundle-60039
 !
 instance VRF_100
 interface bundle-60012.1
 interface bundle-60039.100
 protocols
 bgp 65000
 route-distinguisher 10.0.1.0:100
 router-id 10.0.1.100
 address-family ipv4-unicast
 export-vpn route-target 65100:100
 fast-reroute enabled
 import-vpn route-target 65100:100
 !
 address-family ipv6-unicast
 export-vpn route-target 65100:100
 fast-reroute enabled
 import-vpn route-target 65100:100
 !
 neighbor 10.2.100.0
 remote-as 65002
 local-as 65100 type no-prepend-replace-as
 address-family ipv4-unicast
 policy TOPO-1-IN in
 !
 address-family ipv6-unicast
 policy TOPO-1-IN in
 !
 !
 neighbor 10.100.200.1
 remote-as 65200
 local-as 65100 type no-prepend-replace-as
 address-family ipv4-unicast
 policy TOPO-1-IN in
 !
 !
 neighbor 10:100:200::1
 remote-as 65200
 local-as 65100 type no-prepend-replace-as
 address-family ipv6-unicast
 policy TOPO-1-IN in
 !
 !
 !
 static
 maximum-routes 6 threshold 4
 address-family ipv4-unicast
 route 172.16.16.16/32
 null0
 !
 route 172.17.17.17/32
 null0
 !
 route 172.18.18.18/32
 null0
 !
 route 172.19.19.19/32
 null0
 !
 route 172.20.20.20/32
 null0
 !
 !
 !
 !
 !
 instance VRF_101
 interface bundle-60012.2
 interface bundle-60039.101
 protocols
 bgp 65000
 route-distinguisher 10.0.1.0:101
 router-id 10.0.1.101
 address-family ipv4-unicast
 export-vpn route-target 65101:101
 fast-reroute enabled
 import-vpn route-target 65101:101
 !
 address-family ipv6-unicast
 export-vpn route-target 65101:101
 fast-reroute enabled
 import-vpn route-target 65101:101
 !
 neighbor 10.2.101.0
 remote-as 65002
 local-as 65101 type no-prepend-replace-as
 address-family ipv4-unicast
 policy TOPO-1-IN in
 !
 address-family ipv6-unicast
 policy TOPO-1-IN in
 !
 !
 neighbor 10.101.201.1
 remote-as 65201
 local-as 65101 type no-prepend-replace-as
 address-family ipv4-unicast
 policy TOPO-1-IN in
 !
 !
 neighbor 10:101:201::1
 remote-as 65201
 local-as 65101 type no-prepend-replace-as
 address-family ipv6-unicast
 policy TOPO-1-IN in
 !
 !
 !
 !
 !
 instance VRF_103
 interface bundle-60012.3
 protocols
 bgp 65000
 route-distinguisher 20.0.1.0:103
 router-id 20.0.1.103
 address-family ipv4-unicast
 export-vpn route-target 65103:103
 fast-reroute enabled
 import-vpn route-target 65103:103
 !
 address-family ipv6-unicast
 export-vpn route-target 65103:103
 fast-reroute enabled
 import-vpn route-target 65103:103
 !
 neighbor 20.103.203.1
 remote-as 65203
 local-as 65103 type no-prepend-replace-as
 address-family ipv4-unicast
 policy TOPO-2-IN in
 !
 !
 neighbor 20:103:203::1
 remote-as 65203
 local-as 65103 type no-prepend-replace-as
 address-family ipv6-unicast
 policy TOPO-2-IN in
 !
 !
 !
 !
 !
 !
 !
 		`).
		Append(t)
}

func Test_TC_8_1(t *testing.T) {
	//8.1
	S1 := "./commands/command"
	S2 := ".txt"
	path := ""
	for i := 8001; i <= 8001; i++ {
		path = S1 + strconv.Itoa(i) + S2
		dut := ondatra.DUT(t, "dal00")
		// updates DUT config with replace config below
		dut.Config().New().
			WithDrivenetsFile(path).
			Append(t)
		//fmt.Println("test ended")
	}
}

func Test_TC_8_2(t *testing.T) {
	//8.1
	path := "./commands/last_good"
	dut := ondatra.DUT(t, "dal00")
	dut.Config().New().
		WithDrivenetsFile(path).
		Append(t)
}

func Test_TC_8_3(t *testing.T) {
	//8.3
	S1 := "./commands/command"
	S2 := ".txt"
	path := ""
	for i := 8001; i <= 8101; i++ {
		path = S1 + strconv.Itoa(i) + S2
		dut := ondatra.DUT(t, "dal00")
		// updates DUT config with replace config below
		dut.Config().New().
			WithDrivenetsFile(path).
			Append(t)
		//fmt.Println("test ended")
	}
}

func Test_TC_8_5(t *testing.T) {
	//8.1
	path := "./commands/wrong_filename"
	dut := ondatra.DUT(t, "dal00")
	dut.Config().New().
		WithDrivenetsFile(path).
		Append(t)
}

func Test_TC_8_6(t *testing.T) {
	//8.1
	path := "./commands/empty.txt"
	dut := ondatra.DUT(t, "dal00")
	dut.Config().New().
		WithDrivenetsFile(path).
		Append(t)
}

func Test_verify_config_fix(t *testing.T) {
	S := ""
	dut := ondatra.DUT(t, "dal00")
	for i := 4001; i <= 4015; i++ {
		S = strconv.Itoa(i)
		// updates DUT config with replace config below
		dut.Config().New().
			WithDrivenetsText(
				`interfaces
 					 lo{{var "lid"}}
 						 description "loopback {{ var "lid" }}"
 				 !`).
			WithVarMap(map[string]string{
				"lid":         S,
				"description": "fofo",
			}).
			Append(t)
		//
	}
	cli := ondatra.DUT(t, "dal00").CLI()
	command := "show system name"
	fmt.Printf("Executing command %s\n", command)
	result := cli.RunResult(t, command)
	fmt.Println(result.Output())
	fmt.Println("++++++++++++++++++++++++++++++++++")
	fmt.Println(result.Error())
	//
	command = "show interfaces | include lo4 | count"
	fmt.Printf("Executing command %s\n", command)
	result = cli.RunResult(t, command)
	fmt.Println(result.Output())
	fmt.Println("++++++++++++++++++++++++++++++++++")
	fmt.Println(result.Error())
	//

}

func Test_TC_9_1(t *testing.T) {
	// TC 7.1
	dut := ondatra.DUT(t, "dal00")
	// updates DUT config with replace config below
	dut.Config().New().
		WithDrivenetsText(
			`interfaces
 			 lo9001
 			 admin-state enabled
 			 description fofo
 			 !
 			 !
 			`).
		Push(t)
}

func Test_TC_9_2(t *testing.T) {
	// TC 7.2
	dut := ondatra.DUT(t, "dal00")
	// updates DUT config with replace config below
	dut.Config().New().
		WithDrivenetsText(
			`system
 in-band-management source-interface lo0
 timezone UTC
 aaa-server
 admin-state enabled
 tacacs
 server priority 1 address 100.64.2.252
 accounting enabled
 authentication enabled
 authorization enabled
 password enc-gAAAAABlpNp2XeC6yfok8VmyeB0s4LfD1T98KX1nand2vcKzQdD8Clvmoh4rIVDrxq0VL_XAWGgaxJme-Kb3sPJesZhKVaCPgg==
 port 4949
 vrf mgmt0
 !
 timers
 hold-down 500
 !
 !
 !
 alarms
 admin-state enabled
 !
 cprl
 icmp
 burst 2
 rate 2
 !
 !
 dns
 domain-name dev.drivenets.com
 server priority 1 ip-address 100.64.0.1
 vrf mgmt0
 !
 server priority 2 ip-address 100.64.15.2
 vrf mgmt0
 !
 !
 ftp
 server
 admin-state disabled
 !
 !
 grpc
 max-sessions 6
 security
 tls server-certificate grpc_dnor
 !
 vrf default
 admin-state disabled
 !
 vrf mgmt0
 admin-state enabled
 !
 !
 high-availability
 disable-recovery-mode
 !
 info
 location WOPA_R410
 !
 logging
 syslog
 event-group all severity info
 event-group aaa severity info
 event-group bfd severity info
 event-group bgp severity info
 event-group clock severity info
 event-group diagnostics severity info
 event-group efm-oam severity info
 event-group fib-manager severity info
 event-group interfaces severity info
 event-group isis severity info
 event-group l2vpn severity info
 event-group lacp severity info
 event-group ldp severity info
 event-group lldp severity info
 event-group management severity info
 event-group monitoring severity info
 event-group mpls severity info
 event-group mpls-oam severity info
 event-group msdp severity info
 event-group multicast severity info
 event-group nat severity info
 event-group ospf severity info
 event-group ospfv3 severity info
 event-group pcep severity info
 event-group pim severity info
 event-group platform severity info
 event-group qos severity info
 event-group rib severity info
 event-group rsvp severity info
 event-group segment-routing severity info
 event-group services severity info
 event-group static-route severity info
 event-group system severity info
 event-group tcp severity info
 event-group vrrp severity info
 source-interface lo0
 suppress-event-list BGP_IPV4_NEIGHBOR_ADJACENCY_BACKWARD_CHANGE,BGP_IPV6_NEIGHBOR_ADJACENCY_BACKWARD_CHANGE
 server 10.10.75.158 vrf mgmt0
 admin-state enabled
 port 514
 protocol udp
 severity warning
 !
 server 10.144.28.196 vrf mgmt0
 !
 server 10.253.4.113 vrf mgmt0
 !
 server 10.253.16.232 vrf mgmt0
 !
 server 10.255.250.35 vrf mgmt0
 !
 server 100.64.2.252 vrf mgmt0
 admin-state enabled
 port 1234
 protocol tcp
 severity info
 !
 !
 !
 login
 ncm
 user dnroot
 password enc-gAAAAABgUIs_1U7_1nOEseJ75fkb7UVQSFOD5hi_wyLlHHb41toEsyUgw15QaCqpU2RYwe0CDuNkp_oLw27NNTC1oOk9jWRJXQ==
 role admin
 !
 !
 user dnroot
 password $6$m7b5lYW.PtBTZ0lr$ibN6YvwGYpPauDUIbNu7XMYRkKOnG3H.gyCQVweYv/lgsXMj65ClO8yeGTYX/Rmy80Bx23AeXOWM47nDXr9sN.
 !
 user dntechsupport
 password $6$m0NpQOgH8CusogMK$7JYzbJt55BxEOTkCG6z.9DNTPhDEbX8TN4Po1/mUJzlu/KaywxSoIRLvNYDdfE/WlVZVS2nfXd5J6S2Vin9D3/
 !
 ipmi
 user admin
 password enc-gAAAAABeANLvFw0fQC2rmQLMAF4UtTDyWIjWR1mDWrZmyNoLByecxKV7XxqJVBwoZUB0Y_3NmicHiFlT__YZX_48pt2ghenEnA==
 !
 user dnroot
 password enc-gAAAAABeANIYYvDD_EFSSRwAsxfNj9gHnUlqQnzhE4-GvvCjRuEb21-WA1LQk3V4YvQPabuWjwnrje9jwWNYsKj10czQ7G7-8A==
 !
 !
 !
 ncp 3
 model NCP-40C hw-model S9700-53DX
 admin-state enabled
 description dn-ncp-3
 hardware
 usb 0
 !
 !
 !
 netconf
 max-sessions 12
 port 830
 session-timeout 60
 vrf default
 admin-state enabled
 !
 vrf mgmt0
 admin-state enabled
 !
 !
 ntp
 server 100.64.0.2 vrf mgmt0
 admin-state enabled
 !
 !
 snmp
 community 32946t8pyrehfdlb vrf mgmt0
 !
 !
 ssh
 server
 vrf mgmt0
 client-list type deny
 client-list 100.64.2.252/32
 !
 !
 !
 telnet
 server
 admin-state disabled
 !
 !
 !
 
 services
 performance-monitoring
 profiles
 endpoint-delay STW_P1
 test-duration probes probe-count 10 probe-interval 1 repeat-interval 60
 !
 !
 simple-twamp session DAL00_DAL01
 admin-state enabled
 destination-address 20.0.2.0
 profile STW_P1
 source-address 20.0.1.0
 !
 !
 simple-twamp
 session-reflector
 admin-state enabled
 !
 session-sender
 admin-state enabled
 !
 !
 twamp
 admin-state enabled
 !
 !
 
 interfaces
 bundle-111
 admin-state enabled
 description bundle_with_member
 mtu 9000
 !
 bundle-666
 admin-state enabled
 !
 bundle-60002
 admin-state enabled
 description SA-5
 mtu 9222
 qos policy INGRESS direction in
 qos policy EGRESS direction out
 !
 bundle-60002.1
 admin-state enabled
 description "TOPO-1 MPLS PEERING"
 ipv4-address 10.1.5.1/30
 qos ip-marking trusted
 qos policy INGRESS direction in
 qos policy EGRESS direction out
 vlan-id 1
 !
 bundle-60002.2
 admin-state enabled
 description "TOPO-1 INTERNET PEERING"
 ipv4-address 10.1.5.5/30
 ipv6-address 10:1:5::5/126
 qos ip-marking trusted
 qos policy INGRESS direction in
 qos policy EGRESS direction out
 vlan-id 2
 !
 bundle-60007
 admin-state enabled
 description "IXIA PORT_3_2"
 mtu 9222
 qos policy INGRESS direction in
 qos policy EGRESS direction out
 !
 bundle-60007.3
 admin-state enabled
 ipv4-address 20.1.26.0/31
 ipv6-address 20:1:26::/127
 qos ip-marking trusted
 qos policy INGRESS direction in
 qos policy EGRESS direction out
 vlan-tags outer-tag 7 inner-tag 3
 !
 bundle-60007.10
 admin-state enabled
 ipv4-address 10.1.20.1/24
 qos ip-marking trusted
 qos policy INGRESS direction in
 qos policy EGRESS direction out
 vlan-tags outer-tag 7 inner-tag 20
 !
 bundle-60008
 admin-state enabled
 description ATL00
 mtu 9222
 qos policy INGRESS direction in
 qos policy EGRESS direction out
 !
 bundle-60008.1
 admin-state enabled
 ipv4-address 10.1.3.1/30
 mpls enabled
 qos ip-marking trusted
 qos policy INGRESS direction in
 qos policy EGRESS direction out
 vlan-id 1
 !
 bundle-60008.3
 admin-state enabled
 ipv4-address 20.1.3.0/31
 ipv6-address 20:1:3::/127
 mpls enabled
 qos ip-marking trusted
 qos policy INGRESS direction in
 qos policy EGRESS direction out
 vlan-id 3
 !
 bundle-60009
 admin-state enabled
 description ATL01
 mtu 9222
 qos policy INGRESS direction in
 qos policy EGRESS direction out
 !
 bundle-60009.1
 admin-state enabled
 ipv4-address 10.1.4.1/30
 mpls enabled
 qos ip-marking trusted
 qos policy INGRESS direction in
 qos policy EGRESS direction out
 vlan-id 1
 !
 bundle-60009.3
 admin-state enabled
 ipv4-address 20.1.4.0/31
 ipv6-address 20:1:4::/127
 mpls enabled
 qos ip-marking trusted
 qos policy INGRESS direction in
 qos policy EGRESS direction out
 vlan-id 3
 !
 bundle-60010
 admin-state enabled
 description "IXIA PORT_3_1"
 mtu 9222
 qos policy INGRESS direction in
 qos policy EGRESS direction out
 !
 bundle-60010.3
 admin-state enabled
 ipv4-address 20.1.20.0/31
 ipv6-address 20:1:20::/127
 mpls enabled
 qos ip-marking trusted
 qos policy INGRESS direction in
 qos policy EGRESS direction out
 vlan-tags outer-tag 3 inner-tag 3
 !
 bundle-60010.10
 admin-state enabled
 ipv4-address 10.1.10.1/24
 qos ip-marking trusted
 qos policy INGRESS direction in
 qos policy EGRESS direction out
 vlan-tags outer-tag 3 inner-tag 10
 !
 bundle-60010.100
 admin-state enabled
 ipv4-address 31.100.0.1/30
 vlan-tags outer-tag 3 inner-tag 100
 !
 bundle-60010.101
 admin-state enabled
 ipv4-address 31.101.0.1/30
 vlan-tags outer-tag 3 inner-tag 101
 !
 bundle-60010.102
 admin-state enabled
 ipv4-address 31.102.0.1/30
 vlan-tags outer-tag 3 inner-tag 102
 !
 bundle-60010.103
 admin-state enabled
 ipv4-address 31.103.0.1/30
 vlan-tags outer-tag 3 inner-tag 103
 !
 bundle-60010.104
 admin-state enabled
 ipv4-address 31.104.0.1/30
 vlan-tags outer-tag 3 inner-tag 104
 !
 bundle-60010.105
 admin-state enabled
 ipv4-address 31.105.0.1/30
 vlan-tags outer-tag 3 inner-tag 105
 !
 bundle-60010.106
 admin-state enabled
 ipv4-address 31.106.0.1/30
 vlan-tags outer-tag 3 inner-tag 106
 !
 bundle-60010.107
 admin-state enabled
 ipv4-address 31.107.0.1/30
 vlan-tags outer-tag 3 inner-tag 107
 !
 bundle-60010.108
 admin-state enabled
 ipv4-address 31.108.0.1/30
 vlan-tags outer-tag 3 inner-tag 108
 !
 bundle-60010.109
 admin-state enabled
 ipv4-address 31.109.0.1/30
 vlan-tags outer-tag 3 inner-tag 109
 !
 bundle-60010.110
 admin-state enabled
 ipv4-address 31.110.0.1/30
 vlan-tags outer-tag 3 inner-tag 110
 !
 bundle-60010.111
 admin-state enabled
 ipv4-address 31.111.0.1/30
 vlan-tags outer-tag 3 inner-tag 111
 !
 bundle-60010.112
 admin-state enabled
 ipv4-address 31.112.0.1/30
 vlan-tags outer-tag 3 inner-tag 112
 !
 bundle-60010.113
 admin-state enabled
 ipv4-address 31.113.0.1/30
 vlan-tags outer-tag 3 inner-tag 113
 !
 bundle-60010.114
 admin-state enabled
 ipv4-address 31.114.0.1/30
 vlan-tags outer-tag 3 inner-tag 114
 !
 bundle-60010.115
 admin-state enabled
 ipv4-address 31.115.0.1/30
 vlan-tags outer-tag 3 inner-tag 115
 !
 bundle-60010.116
 admin-state enabled
 ipv4-address 31.116.0.1/30
 vlan-tags outer-tag 3 inner-tag 116
 !
 bundle-60010.117
 admin-state enabled
 ipv4-address 31.117.0.1/30
 vlan-tags outer-tag 3 inner-tag 117
 !
 bundle-60010.118
 admin-state enabled
 ipv4-address 31.118.0.1/30
 vlan-tags outer-tag 3 inner-tag 118
 !
 bundle-60010.119
 admin-state enabled
 ipv4-address 31.119.0.1/30
 vlan-tags outer-tag 3 inner-tag 119
 !
 bundle-60010.120
 admin-state enabled
 ipv4-address 31.120.0.1/30
 vlan-tags outer-tag 3 inner-tag 120
 !
 bundle-60010.121
 admin-state enabled
 ipv4-address 31.121.0.1/30
 vlan-tags outer-tag 3 inner-tag 121
 !
 bundle-60010.122
 admin-state enabled
 ipv4-address 31.122.0.1/30
 vlan-tags outer-tag 3 inner-tag 122
 !
 bundle-60010.123
 admin-state enabled
 ipv4-address 31.123.0.1/30
 vlan-tags outer-tag 3 inner-tag 123
 !
 bundle-60010.124
 admin-state enabled
 ipv4-address 31.124.0.1/30
 vlan-tags outer-tag 3 inner-tag 124
 !
 bundle-60010.125
 admin-state enabled
 ipv4-address 31.125.0.1/30
 vlan-tags outer-tag 3 inner-tag 125
 !
 bundle-60010.126
 admin-state enabled
 ipv4-address 31.126.0.1/30
 vlan-tags outer-tag 3 inner-tag 126
 !
 bundle-60010.127
 admin-state enabled
 ipv4-address 31.127.0.1/30
 vlan-tags outer-tag 3 inner-tag 127
 !
 bundle-60010.128
 admin-state enabled
 ipv4-address 31.128.0.1/30
 vlan-tags outer-tag 3 inner-tag 128
 !
 bundle-60010.129
 admin-state enabled
 ipv4-address 31.129.0.1/30
 vlan-tags outer-tag 3 inner-tag 129
 !
 bundle-60010.130
 admin-state enabled
 ipv4-address 31.130.0.1/30
 vlan-tags outer-tag 3 inner-tag 130
 !
 bundle-60010.131
 admin-state enabled
 ipv4-address 31.131.0.1/30
 vlan-tags outer-tag 3 inner-tag 131
 !
 bundle-60010.132
 admin-state enabled
 ipv4-address 31.132.0.1/30
 vlan-tags outer-tag 3 inner-tag 132
 !
 bundle-60010.133
 admin-state enabled
 ipv4-address 31.133.0.1/30
 vlan-tags outer-tag 3 inner-tag 133
 !
 bundle-60010.134
 admin-state enabled
 ipv4-address 31.134.0.1/30
 vlan-tags outer-tag 3 inner-tag 134
 !
 bundle-60010.135
 admin-state enabled
 ipv4-address 31.135.0.1/30
 vlan-tags outer-tag 3 inner-tag 135
 !
 bundle-60010.136
 admin-state enabled
 ipv4-address 31.136.0.1/30
 vlan-tags outer-tag 3 inner-tag 136
 !
 bundle-60010.137
 admin-state enabled
 ipv4-address 31.137.0.1/30
 vlan-tags outer-tag 3 inner-tag 137
 !
 bundle-60010.138
 admin-state enabled
 ipv4-address 31.138.0.1/30
 vlan-tags outer-tag 3 inner-tag 138
 !
 bundle-60010.139
 admin-state enabled
 ipv4-address 31.139.0.1/30
 vlan-tags outer-tag 3 inner-tag 139
 !
 bundle-60010.140
 admin-state enabled
 ipv4-address 31.140.0.1/30
 vlan-tags outer-tag 3 inner-tag 140
 !
 bundle-60010.141
 admin-state enabled
 ipv4-address 31.141.0.1/30
 vlan-tags outer-tag 3 inner-tag 141
 !
 bundle-60010.142
 admin-state enabled
 ipv4-address 31.142.0.1/30
 vlan-tags outer-tag 3 inner-tag 142
 !
 bundle-60010.143
 admin-state enabled
 ipv4-address 31.143.0.1/30
 vlan-tags outer-tag 3 inner-tag 143
 !
 bundle-60010.144
 admin-state enabled
 ipv4-address 31.144.0.1/30
 vlan-tags outer-tag 3 inner-tag 144
 !
 bundle-60010.145
 admin-state enabled
 ipv4-address 31.145.0.1/30
 vlan-tags outer-tag 3 inner-tag 145
 !
 bundle-60010.146
 admin-state enabled
 ipv4-address 31.146.0.1/30
 vlan-tags outer-tag 3 inner-tag 146
 !
 bundle-60010.147
 admin-state enabled
 ipv4-address 31.147.0.1/30
 vlan-tags outer-tag 3 inner-tag 147
 !
 bundle-60010.148
 admin-state enabled
 ipv4-address 31.148.0.1/30
 vlan-tags outer-tag 3 inner-tag 148
 !
 bundle-60010.149
 admin-state enabled
 ipv4-address 31.149.0.1/30
 vlan-tags outer-tag 3 inner-tag 149
 !
 bundle-60010.150
 admin-state enabled
 ipv4-address 31.150.0.1/30
 vlan-tags outer-tag 3 inner-tag 150
 !
 bundle-60011
 admin-state enabled
 description DAL00
 ipv4-address 192.168.0.1/30
 mpls disabled
 mtu 9222
 qos policy INGRESS direction in
 qos policy EGRESS direction out
 !
 bundle-60011.1
 admin-state enabled
 ipv4-address 10.1.2.1/30
 mpls enabled
 qos ip-marking trusted
 qos policy INGRESS direction in
 qos policy EGRESS direction out
 vlan-id 1
 !
 bundle-60011.2
 admin-state enabled
 ipv4-address 10.1.2.5/30
 qos ip-marking trusted
 qos policy INGRESS direction in
 qos policy EGRESS direction out
 vlan-id 2
 !
 bundle-60011.3
 admin-state enabled
 ipv4-address 20.1.2.0/31
 ipv6-address 20:1:2::/127
 mpls enabled
 qos ip-marking trusted
 qos policy INGRESS direction in
 qos policy EGRESS direction out
 vlan-id 3
 !
 bundle-60012
 admin-state enabled
 description SA-5
 mtu 9222
 qos policy INGRESS direction in
 qos policy EGRESS direction out
 !
 bundle-60012.1
 admin-state enabled
 ipv4-address 10.100.200.0/31
 ipv6-address 10:100:200::/127
 qos ip-marking trusted
 qos policy INGRESS direction in
 qos policy EGRESS direction out
 vlan-id 1
 !
 bundle-60012.2
 admin-state enabled
 ipv4-address 10.101.201.0/31
 ipv6-address 10:101:201::/127
 qos ip-marking trusted
 qos policy INGRESS direction in
 qos policy EGRESS direction out
 vlan-id 2
 !
 bundle-60012.3
 admin-state enabled
 ipv4-address 20.103.203.0/31
 ipv6-address 20:103:203::/127
 qos ip-marking trusted
 qos policy INGRESS direction in
 qos policy EGRESS direction out
 vlan-id 3
 !
 bundle-60038
 admin-state enabled
 description "TOPO-1 LOOP_BUNDLE_60038"
 ipv4-address 10.1.1.0/31
 mtu 9222
 qos policy INGRESS direction in
 qos policy EGRESS direction out
 !
 bundle-60038.100
 admin-state enabled
 ipv4-address 10.2.100.0/31
 ipv6-address 10:2:100::/127
 qos ip-marking trusted
 qos policy INGRESS direction in
 qos policy EGRESS direction out
 vlan-id 100
 !
 bundle-60038.101
 admin-state enabled
 ipv4-address 10.2.101.0/31
 ipv6-address 10:2:101::/127
 qos ip-marking trusted
 qos policy INGRESS direction in
 qos policy EGRESS direction out
 vlan-id 101
 !
 bundle-60038.999
 admin-state enabled
 ipv4-address 169.254.0.0/31
 qos ip-marking trusted
 qos policy INGRESS direction in
 qos policy EGRESS direction out
 vlan-id 999
 !
 bundle-60039
 admin-state enabled
 description "TOPO-1 LOOP_BUNDLE_60039"
 ipv4-address 10.1.1.1/31
 mtu 9222
 qos policy INGRESS direction in
 qos policy EGRESS direction out
 !
 bundle-60039.100
 admin-state enabled
 ipv4-address 10.2.100.1/31
 ipv6-address 10:2:100::1/127
 qos ip-marking trusted
 qos policy INGRESS direction in
 qos policy EGRESS direction out
 vlan-id 100
 !
 bundle-60039.101
 admin-state enabled
 ipv4-address 10.2.101.1/31
 ipv6-address 10:2:101::1/127
 qos ip-marking trusted
 qos policy INGRESS direction in
 qos policy EGRESS direction out
 vlan-id 101
 !
 bundle-60039.999
 admin-state enabled
 ipv4-address 169.254.0.1/31
 qos ip-marking trusted
 qos policy INGRESS direction in
 qos policy EGRESS direction out
 vlan-id 999
 !
 console-ncp-3/0
 admin-state enabled
 !
 ctrl-ncp-3/0
 admin-state enabled
 !
 ctrl-ncp-3/1
 admin-state enabled
 !
 fab-ncp400-3/0/0
 admin-state enabled
 !
 fab-ncp400-3/0/1
 admin-state enabled
 !
 fab-ncp400-3/0/2
 admin-state enabled
 !
 fab-ncp400-3/0/3
 admin-state enabled
 !
 fab-ncp400-3/0/4
 admin-state enabled
 !
 fab-ncp400-3/0/5
 admin-state enabled
 !
 fab-ncp400-3/0/6
 admin-state enabled
 !
 fab-ncp400-3/0/7
 admin-state enabled
 !
 fab-ncp400-3/0/8
 admin-state enabled
 !
 fab-ncp400-3/0/9
 admin-state enabled
 !
 fab-ncp400-3/0/10
 admin-state enabled
 !
 fab-ncp400-3/0/11
 admin-state enabled
 !
 fab-ncp400-3/0/12
 admin-state enabled
 !
 ge100-3/0/0
 admin-state enabled
 bundle-id 60010
 fec none
 !
 ge100-3/0/1
 admin-state enabled
 bundle-id 60007
 fec none
 !
 ge100-3/0/2
 admin-state enabled
 fec none
 !
 ge100-3/0/3
 admin-state enabled
 fec none
 !
 ge100-3/0/4
 admin-state enabled
 bundle-id 60008
 fec rs-fec-528-514
 !
 ge100-3/0/5
 admin-state enabled
 bundle-id 60009
 fec none
 !
 ge100-3/0/6
 admin-state enabled
 bundle-id 60002
 fec none
 !
 ge100-3/0/7
 admin-state enabled
 bundle-id 60012
 fec none
 !
 ge100-3/0/8
 admin-state enabled
 fec none
 !
 ge100-3/0/9
 admin-state enabled
 bundle-id 60011
 fec none
 !
 ge100-3/0/10
 admin-state enabled
 fec none
 !
 ge100-3/0/11
 admin-state enabled
 bundle-id 60011
 fec none
 !
 ge100-3/0/12
 admin-state enabled
 fec none
 !
 ge100-3/0/13
 admin-state enabled
 fec none
 !
 ge100-3/0/14
 admin-state enabled
 fec none
 !
 ge100-3/0/15
 admin-state enabled
 fec none
 !
 ge100-3/0/16
 admin-state enabled
 fec none
 !
 ge100-3/0/17
 admin-state enabled
 fec none
 !
 ge100-3/0/18
 admin-state enabled
 fec none
 !
 ge100-3/0/19
 admin-state enabled
 fec none
 !
 ge100-3/0/20
 admin-state enabled
 fec none
 !
 ge100-3/0/21
 admin-state enabled
 fec none
 !
 ge100-3/0/22
 admin-state enabled
 fec none
 !
 ge100-3/0/23
 admin-state enabled
 fec none
 !
 ge100-3/0/24
 admin-state enabled
 fec none
 !
 ge100-3/0/25
 admin-state enabled
 fec none
 !
 ge100-3/0/26
 admin-state enabled
 fec none
 !
 ge100-3/0/27
 admin-state enabled
 fec none
 !
 ge100-3/0/28
 admin-state disabled
 fec none
 util-rate-threshold 100
 !
 ge100-3/0/29
 admin-state enabled
 fec none
 !
 ge100-3/0/30
 admin-state enabled
 fec none
 !
 ge100-3/0/31
 admin-state enabled
 fec none
 !
 ge100-3/0/32
 admin-state enabled
 fec none
 !
 ge100-3/0/33
 admin-state enabled
 fec none
 !
 ge100-3/0/34
 admin-state enabled
 fec none
 !
 ge100-3/0/35
 admin-state enabled
 fec none
 !
 ge100-3/0/36
 admin-state enabled
 fec none
 !
 ge100-3/0/37
 admin-state enabled
 fec none
 !
 ge100-3/0/38
 admin-state enabled
 bundle-id 60038
 fec none
 dampening
 admin-state enabled
 !
 !
 ge100-3/0/39
 admin-state disabled
 bundle-id 60039
 fec none
 !
 ipmi-ncp-0/0
 admin-state enabled
 !
 ipmi-ncp-3/0
 admin-state enabled
 !
 lo0
 admin-state enabled
 description TOPO-1
 ipv4-address 10.0.1.0/32
 !
 lo2
 admin-state enabled
 description TOPO-2
 ipv4-address 20.0.1.0/32
 ipv6-address 20:0:1::/128
 !
 lo10
 admin-state enabled
 ipv4-address 100.0.0.1/32
 ipv6-address 100::1/128
 !
 lo900
 admin-state enabled
 description new_loopback
 !
 lo999
 admin-state enabled
 ipv4-address 9.9.9.9/32
 !
 lo1000
 admin-state enabled
 ipv4-address 1.1.1.1/32
 !
 lo12345
 admin-state enabled
 !
 mgmt0
 admin-state enabled
 !
 !
 
 routing-options
 maximum-routes 100000 threshold 20
 !
 
 protocols
 bgp nsr disabled
 bfd
 interface bundle-60011
 local-address 192.168.0.1
 neighbor 192.168.0.2
 !
 interface bundle-60038
 local-address 10.1.1.0
 neighbor 10.1.1.1
 !
 interface bundle-60039
 local-address 10.1.1.1
 neighbor 10.1.1.0
 !
 !
 bgp 65000
 network import-check disabled
 route-reflection policy-out attribute-change enabled
 router-id 10.0.1.0
 address-family ipv4-vpn
 fast-reroute enabled
 label-allocation per-nexthop
 !
 address-family ipv6-vpn
 fast-reroute enabled
 label-allocation per-nexthop
 !
 neighbor 20.1.20.1
 remote-as 65000
 update-source bundle-60010.3
 address-family ipv4-unicast
 labeled-unicast
 !
 !
 neighbor 31.100.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.100
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.101.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.101
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.102.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.102
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.103.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.103
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.104.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.104
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.105.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.105
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.106.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.106
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.107.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.107
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.108.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.108
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.109.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.109
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.110.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.110
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.111.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.111
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.112.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.112
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.113.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.113
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.114.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.114
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.115.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.115
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.116.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.116
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.117.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.117
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.118.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.118
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.119.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.119
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.120.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.120
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.121.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.121
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.122.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.122
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.123.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.123
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.124.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.124
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.125.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.125
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.126.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.126
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.127.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.127
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.128.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.128
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.129.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.129
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.130.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.130
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.131.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.131
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.132.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.132
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.133.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.133
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.134.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.134
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.135.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.135
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.136.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.136
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.137.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.137
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.138.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.138
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.139.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.139
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.140.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.140
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.141.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.141
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.142.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.142
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.143.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.143
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.144.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.144
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.145.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.145
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.146.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.146
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.147.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.147
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.148.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.148
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.149.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.149
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor 31.150.0.2
 remote-as 65000
 admin-state disabled
 update-source bundle-60010.150
 bfd
 admin-state enabled
 min-rx 500
 min-tx 500
 multiplier 5
 !
 address-family ipv4-unicast
 !
 !
 neighbor-group NC_CLIENTS
 remote-as 65000
 update-source lo0
 address-family ipv4-rt-constrains
 route-reflector-client
 !
 address-family ipv4-vpn
 policy TOPO-1-OUT out
 route-reflector-client
 send-community
 send-large-community enabled
 !
 address-family ipv6-vpn
 policy TOPO-1-OUT out
 route-reflector-client
 send-community
 send-large-community enabled
 !
 neighbor 10.0.2.0
 description DAL01
 !
 neighbor 10.0.4.0
 description ATL01
 !
 !
 neighbor-group RR_FULL_MESH
 remote-as 65000
 update-source lo0
 address-family ipv4-rt-constrains
 !
 address-family ipv4-vpn
 policy TOPO-1-OUT out
 send-community
 send-large-community enabled
 !
 address-family ipv6-vpn
 policy TOPO-1-IN out
 send-community
 send-large-community enabled
 !
 neighbor 10.0.3.0
 description ATL00
 !
 !
 neighbor-group TOPO-2-EBGP
 remote-as 65500
 local-as 65499 type no-prepend-replace-as
 address-family ipv4-unicast
 dampening enabled
 !
 neighbor 20.1.2.1
 !
 neighbor 20.1.3.1
 !
 neighbor 20.1.4.1
 !
 !
 neighbor-group TOPO-2-FULL-MESH
 remote-as 65000
 update-source lo2
 address-family ipv4-vpn
 policy TOPO-2-OUT out
 send-community
 !
 address-family ipv6-vpn
 policy TOPO-2-OUT out
 send-community
 !
 neighbor 20.0.2.0
 !
 neighbor 20.0.3.0
 !
 neighbor 20.0.4.0
 !
 !
 !
 isis
 maximum-routes 300 threshold 50
 nsr enabled
 instance CORE
 iso-network 49.0000.0100.0000.1000.00
 level level-2
 lsp-mtu 1400
 max-metric 1000000
 address-family ipv4-unicast
 traffic-engineering enabled
 segment-routing
 admin-state enabled
 !
 !
 interface lo0
 address-family ipv4-unicast
 prefix-sid index 1 prefix-type node
 !
 !
 interface bundle-60007.10
 address-family ipv4-unicast
 !
 !
 interface bundle-60008.1
 network-type point-to-point
 address-family ipv4-unicast
 metric 5
 !
 bfd
 admin-state enabled
 min-rx 100
 min-tx 100
 multiplier 3
 !
 !
 interface bundle-60009.1
 network-type point-to-point
 address-family ipv4-unicast
 metric 100
 !
 bfd
 admin-state enabled
 min-rx 100
 min-tx 100
 multiplier 3
 !
 !
 interface bundle-60010.10
 address-family ipv4-unicast
 !
 !
 interface bundle-60010.100
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.101
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.102
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.103
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.104
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.105
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.106
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.107
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.108
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.109
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.110
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.111
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.112
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.113
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.114
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.115
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.116
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.117
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.118
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.119
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.120
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.121
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.122
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.123
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.124
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.125
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.126
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.127
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.128
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.129
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.130
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.131
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.132
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.133
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.134
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.135
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.136
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.137
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.138
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.139
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.140
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.141
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.142
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.143
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.144
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.145
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.146
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.147
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.148
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.149
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60010.150
 address-family ipv4-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 5
 min-tx 5
 multiplier 5
 !
 !
 interface bundle-60011.1
 network-type point-to-point
 address-family ipv4-unicast
 metric 10
 !
 address-family ipv6-unicast
 metric 10
 !
 bfd
 admin-state enabled
 min-rx 100
 min-tx 100
 multiplier 3
 !
 !
 interface bundle-60011.2
 network-type point-to-point
 address-family ipv4-unicast
 metric 800
 !
 !
 timers
 lsp-lifetime 65535
 lsp-refresh 65000
 !
 !
 instance TOPO-2
 administrative-distance 117
 iso-network 49.0000.0100.0000.1000.00
 level level-2
 lsp-mtu 1400
 address-family ipv4-unicast
 traffic-engineering enabled
 !
 interface lo2
 address-family ipv4-unicast
 !
 address-family ipv6-unicast
 !
 !
 interface bundle-60007.3
 address-family ipv4-unicast
 metric 10
 !
 !
 interface bundle-60008.3
 address-family ipv4-unicast
 metric 10
 !
 address-family ipv6-unicast
 metric 10
 !
 !
 interface bundle-60009.3
 address-family ipv4-unicast
 metric 100
 !
 address-family ipv6-unicast
 metric 100
 !
 bfd
 admin-state enabled
 !
 !
 interface bundle-60010.3
 address-family ipv4-unicast
 metric 10
 !
 address-family ipv6-unicast
 metric 10
 !
 !
 interface bundle-60011.3
 address-family ipv4-unicast
 metric 10
 !
 address-family ipv6-unicast
 metric 10
 !
 bfd
 admin-state enabled
 !
 !
 overload on-startup
 admin-state enabled
 advertisement-type max-metric
 interval 600
 wait-for-bgp bgp-delay 60
 !
 !
 !
 lacp
 interface bundle-60002
 mode active
 period short
 !
 interface bundle-60008
 mode active
 period short
 !
 interface bundle-60009
 mode active
 period short
 !
 interface bundle-60011
 mode active
 period short
 !
 interface bundle-60012
 mode active
 period short
 !
 !
 lldp
 interface ge100-3/0/2
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/3
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/4
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/5
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/6
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/7
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/8
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/9
 !
 interface ge100-3/0/10
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/11
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/12
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/13
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/14
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/15
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/16
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/17
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/18
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/19
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/20
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/21
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/22
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/23
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/24
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/25
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/26
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/27
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/29
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/32
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/33
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/34
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/35
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/36
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/37
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/38
 receive enabled
 transmit enabled
 !
 interface ge100-3/0/39
 receive enabled
 transmit enabled
 !
 !
 mpls
 traffic-engineering
 router-id 20.0.1.0
 interface bundle-60008.3
 !
 interface bundle-60009.3
 !
 interface bundle-60010.3
 !
 interface bundle-60011.3
 !
 !
 !
 pim
 address-family ipv4
 interface bundle-60007.3
 admin-state enabled
 !
 interface bundle-60008.3
 admin-state enabled
 !
 interface bundle-60010.3
 admin-state enabled
 !
 interface bundle-60011.3
 admin-state enabled
 !
 !
 static-rp 20.0.22.0
 !
 !
 rsvp
 auto-mesh
 tunnel-template FULL_MESH_TOPO_2
 admin-state enabled
 destination-address TOPO-2-LO0
 igp-instance isis TOPO-2
 rib-unicast-install enabled
 source-address 20.0.1.0
 bfd
 admin-state enabled
 !
 primary
 protection node-protection
 path-options 1
 path dynamic
 !
 !
 !
 tunnel-template FULL_MESH_TOPO_2_IXIA
 admin-state enabled
 destination-address TOPO-2-IXIA
 igp-instance isis TOPO-2
 rib-unicast-install enabled
 source-address 20.0.1.0
 primary
 path-options 1
 path dynamic
 !
 !
 !
 !
 interface bundle-60008.3
 protection
 auto-bypass enabled
 !
 !
 interface bundle-60009.3
 protection
 auto-bypass enabled
 !
 !
 interface bundle-60010.3
 !
 interface bundle-60011.3
 protection
 auto-bypass enabled
 !
 !
 tunnel BBB bypass
 admin-state enabled
 destination-address 10.0.3.0
 primary
 path-options 1
 path dynamic
 !
 !
 !
 tunnel AAA
 admin-state enabled
 destination-address 20.0.2.0
 igp-instance isis TOPO-2
 primary
 path-options 1
 path dynamic
 !
 !
 !
 tunnel DAL00-ATL01
 admin-state disabled
 destination-address 20.0.4.0
 igp-instance isis TOPO-2
 source-address 20.0.1.0
 primary
 path-options 1
 path dynamic
 !
 !
 !
 !
 static
 maximum-routes 6 threshold 4
 !
 !
 
 multicast
 rpf-intact admin-state enabled
 !
 
 routing-policy
 community-list TOPO-1
 rule 1 allow value 65000:10
 !
 community-list TOPO-2
 rule 1 allow value 65000:20
 !
 policy TOPO-1-IN
 rule 1 allow
 set community 65000:10
 !
 !
 policy TOPO-1-OUT
 rule 10 allow
 match community TOPO-1
 !
 rule 100 allow
 !
 !
 policy TOPO-2-IN
 rule 1 allow
 set community 65000:20
 !
 !
 policy TOPO-2-OUT
 rule 10 allow
 match community TOPO-2
 !
 !
 prefix-list ipv4 TOPO-2-IXIA
 rule 1 deny 20.0.1.0/32,
 rule 2 deny 20.0.2.0/32,
 rule 3 deny 20.0.3.0/32,
 rule 4 deny 20.0.4.0/32,
 rule 5 allow 20.0.0.0/16 matching-len eq 32
 !
 prefix-list ipv4 TOPO-2-LO0
 rule 1 allow 20.0.2.0/32,
 rule 2 allow 20.0.4.0/32
 !
 !
 
 qos
 hw-mapping
 queue-size
 speed-ranges
 upto 50 mbps use 50 mbps
 upto 100 mbps use 100 mbps
 upto 250 mbps use 250 mbps
 upto 500 mbps use 500 mbps
 upto 750 mbps use 750 mbps
 upto 1 gbps use 1 gbps
 !
 !
 !
 wred-profile WRED-20K-DUAL
 curve green
 min 16 milliseconds max 20 milliseconds
 !
 curve yellow
 min 8 milliseconds max 12 milliseconds
 !
 !
 wred-profile WRED-25K-SINGLE
 curve green
 min 20 milliseconds max 25 milliseconds
 !
 !
 policy EGRESS
 rule 6
 match traffic-class CLASS-CONTROL
 action
 queue
 forwarding-class af
 bandwidth 4 percent
 size 5 milliseconds
 !
 !
 !
 !
 rule default
 action
 queue
 forwarding-class df
 bandwidth 30 percent
 size 25 milliseconds
 wred-profile WRED-25K-SINGLE
 !
 !
 !
 !
 !
 policy INGRESS
 rule 6
 match traffic-class CLASS6
 action
 set
 qos-tag 6
 !
 !
 !
 rule default
 !
 !
 traffic-class-map CLASS-CONTROL
 qos-tag 6
 !
 traffic-class-map CLASS6
 mpls-exp 6
 precedence 6
 !
 !
 
 network-services
 vrf
 instance DAL00_A
 interface bundle-60038.999
 protocols
 bgp 65000
 router-id 169.254.0.0
 neighbor 169.254.0.1
 remote-as 4290000002
 local-as 4290000001 type no-prepend-replace-as
 address-family ipv4-unicast
 !
 !
 !
 !
 !
 instance DAL00_B
 interface bundle-60039.999
 !
 instance INTERNET
 interface bundle-60002.2
 interface bundle-60038.100
 interface bundle-60038.101
 protocols
 bgp 65000
 route-distinguisher 10.0.1.0:2
 router-id 10.0.1.2
 address-family ipv4-unicast
 export-vpn route-target 65002:2
 fast-reroute enabled
 import-vpn route-target 65002:2
 !
 address-family ipv6-unicast
 export-vpn route-target 65002:2
 fast-reroute enabled
 import-vpn route-target 65002:2
 !
 neighbor 10.1.5.6
 remote-as 64512
 advertisement-interval 30
 address-family ipv4-unicast
 as-loop-check disabled
 policy TOPO-1-IN in
 !
 !
 neighbor 10.2.100.1
 remote-as 65100
 local-as 65002 type no-prepend-replace-as
 address-family ipv4-unicast
 policy TOPO-1-IN in
 !
 address-family ipv6-unicast
 policy TOPO-1-IN in
 !
 !
 neighbor 10.2.101.1
 remote-as 65101
 local-as 65002 type no-prepend-replace-as
 address-family ipv4-unicast
 policy TOPO-1-IN in
 !
 address-family ipv6-unicast
 policy TOPO-1-IN in
 !
 !
 neighbor 10:1:5::6
 remote-as 64512
 address-family ipv6-unicast
 policy TOPO-1-IN in
 !
 !
 !
 static
 address-family ipv4-unicast
 route 200.0.0.0/24
 next-hop 10.2.100.1
 !
 !
 !
 !
 !
 instance LEFT
 interface bundle-60038
 !
 instance RIGHT
 interface bundle-60039
 !
 instance VRF_100
 interface bundle-60012.1
 interface bundle-60039.100
 protocols
 bgp 65000
 route-distinguisher 10.0.1.0:100
 router-id 10.0.1.100
 address-family ipv4-unicast
 export-vpn route-target 65100:100
 fast-reroute enabled
 import-vpn route-target 65100:100
 !
 address-family ipv6-unicast
 export-vpn route-target 65100:100
 fast-reroute enabled
 import-vpn route-target 65100:100
 !
 neighbor 10.2.100.0
 remote-as 65002
 local-as 65100 type no-prepend-replace-as
 address-family ipv4-unicast
 policy TOPO-1-IN in
 !
 address-family ipv6-unicast
 policy TOPO-1-IN in
 !
 !
 neighbor 10.100.200.1
 remote-as 65200
 local-as 65100 type no-prepend-replace-as
 address-family ipv4-unicast
 policy TOPO-1-IN in
 !
 !
 neighbor 10:100:200::1
 remote-as 65200
 local-as 65100 type no-prepend-replace-as
 address-family ipv6-unicast
 policy TOPO-1-IN in
 !
 !
 !
 static
 maximum-routes 6 threshold 4
 address-family ipv4-unicast
 route 172.16.16.16/32
 null0
 !
 route 172.17.17.17/32
 null0
 !
 route 172.18.18.18/32
 null0
 !
 route 172.19.19.19/32
 null0
 !
 route 172.20.20.20/32
 null0
 !
 !
 !
 !
 !
 instance VRF_101
 interface bundle-60012.2
 interface bundle-60039.101
 protocols
 bgp 65000
 route-distinguisher 10.0.1.0:101
 router-id 10.0.1.101
 address-family ipv4-unicast
 export-vpn route-target 65101:101
 fast-reroute enabled
 import-vpn route-target 65101:101
 !
 address-family ipv6-unicast
 export-vpn route-target 65101:101
 fast-reroute enabled
 import-vpn route-target 65101:101
 !
 neighbor 10.2.101.0
 remote-as 65002
 local-as 65101 type no-prepend-replace-as
 address-family ipv4-unicast
 policy TOPO-1-IN in
 !
 address-family ipv6-unicast
 policy TOPO-1-IN in
 !
 !
 neighbor 10.101.201.1
 remote-as 65201
 local-as 65101 type no-prepend-replace-as
 address-family ipv4-unicast
 policy TOPO-1-IN in
 !
 !
 neighbor 10:101:201::1
 remote-as 65201
 local-as 65101 type no-prepend-replace-as
 address-family ipv6-unicast
 policy TOPO-1-IN in
 !
 !
 !
 !
 !
 instance VRF_103
 interface bundle-60012.3
 protocols
 bgp 65000
 route-distinguisher 20.0.1.0:103
 router-id 20.0.1.103
 address-family ipv4-unicast
 export-vpn route-target 65103:103
 fast-reroute enabled
 import-vpn route-target 65103:103
 !
 address-family ipv6-unicast
 export-vpn route-target 65103:103
 fast-reroute enabled
 import-vpn route-target 65103:103
 !
 neighbor 20.103.203.1
 remote-as 65203
 local-as 65103 type no-prepend-replace-as
 address-family ipv4-unicast
 policy TOPO-2-IN in
 !
 !
 neighbor 20:103:203::1
 remote-as 65203
 local-as 65103 type no-prepend-replace-as
 address-family ipv6-unicast
 policy TOPO-2-IN in
 !
 !
 !
 !
 !
 !
 !
 		`).
		Push(t)
}

func Test_TC_10_1(t *testing.T) {
	//8.1
	S1 := "./commands/command"
	S2 := ".txt"
	path := ""
	for i := 8001; i <= 8001; i++ {
		path = S1 + strconv.Itoa(i) + S2
		dut := ondatra.DUT(t, "dal00")
		// updates DUT config with replace config below
		dut.Config().New().
			WithDrivenetsFile(path).
			Push(t)
		//fmt.Println("test ended")
	}
}

func Test_TC_10_2(t *testing.T) {
	//8.1
	path := "./commands/last_good"
	dut := ondatra.DUT(t, "dal00")
	dut.Config().New().
		WithDrivenetsFile(path).
		Push(t)
}

func Test_TC_11_1(t *testing.T) {
	S := ""
	dut := ondatra.DUT(t, "dal00")
	for i := 4001; i <= 4002; i++ {
		S = strconv.Itoa(i)
		// updates DUT config with replace config below
		dut.Config().New().
			WithDrivenetsText(
				`interfaces
 					 lo{{var "lid"}}
 						 description "int lo nuber{{ var "description" }}"
 				 !`).
			WithVarMap(map[string]string{
				"lid":         S,
				"description": "",
			}).
			Append(t)
		//
	}
	cli := ondatra.DUT(t, "dal00").CLI()
	command := "show system name"
	fmt.Printf("Executing command %s\n", command)
	result := cli.RunResult(t, command)
	fmt.Println(result.Output())
	fmt.Println("++++++++++++++++++++++++++++++++++")
	fmt.Println(result.Error())
	//
	command = "show interfaces | include lo4 | count"
	fmt.Printf("Executing command %s\n", command)
	result = cli.RunResult(t, command)
	fmt.Println(result.Output())
	fmt.Println("++++++++++++++++++++++++++++++++++")
	fmt.Println(result.Error())
	//

}

func Test_TC_11_2(t *testing.T) {
	S := ""
	dut := ondatra.DUT(t, "dal00")
	for i := 4001; i <= 4002; i++ {
		S = strconv.Itoa(i)
		// updates DUT config with replace config below
		dut.Config().New().
			WithDrivenetsText(
				`interfaces
 					 lo{{var "lid"}}
 						 description "{{ var "lid" }}"
 				 !`).
			WithVarMap(map[string]string{
				"lid":         S,
				"description": "",
			}).
			Append(t)
		//
	}
	cli := ondatra.DUT(t, "dal00").CLI()
	command := "show system name"
	fmt.Printf("Executing command %s\n", command)
	result := cli.RunResult(t, command)
	fmt.Println(result.Output())
	fmt.Println("++++++++++++++++++++++++++++++++++")
	fmt.Println(result.Error())
	//
	command = "show interfaces | include lo4 | count"
	fmt.Printf("Executing command %s\n", command)
	result = cli.RunResult(t, command)
	fmt.Println(result.Output())
	fmt.Println("++++++++++++++++++++++++++++++++++")
	fmt.Println(result.Error())
	//

}

func Test_TC_12_1(t *testing.T) {
	command := ""
	// TC 12.1
	cli := ondatra.DUT(t, "dal00").CLI()
	cli2 := ondatra.DUT(t, "dal01").CLI()

	//
	command = "show system name"
	fmt.Printf("Executing command %s\n", command)
	fmt.Println("Command output is : ")
	result := cli.Run(t, command)
	fmt.Println(result)

	//
	command = "show system name"
	fmt.Printf("Executing command %s\n", command)
	fmt.Println("Command output is : ")
	result = cli2.Run(t, command)
	fmt.Println(result)

}

func Test_TC_12_2(t *testing.T) {
	command := ""
	// TC 12.1
	cli := ondatra.DUT(t, "dal00").CLI()

	//
	command = "show system name"
	fmt.Printf("Executing command %s\n", command)
	fmt.Println("Command output is : ")
	result := cli.Run(t, command)
	fmt.Println(result)

	//
	cli = ondatra.DUT(t, "dal01").CLI()
	command = "show system name"
	fmt.Printf("Executing command %s\n", command)
	fmt.Println("Command output is : ")
	result = cli.Run(t, command)
	fmt.Println(result)

}

func Test_TC_12_3(t *testing.T) {
	dut := ondatra.DUT(t, "dal00")
	dut2 := ondatra.DUT(t, "dal01")
	dut.Config().New().
		WithDrivenetsText(
			`system name DALLAS00
 			`).
		Append(t)
		//
	dut2.Config().New().
		WithDrivenetsText(
			`system name DALLAS01
 			`).
		Append(t)
	//
}

func Test_TC_12_4(t *testing.T) {
	dut := ondatra.DUT(t, "dal00")

	dut.Config().New().
		WithDrivenetsText(
			`system name DAL00
 			`).
		Append(t)
		//
	dut = ondatra.DUT(t, "dal01")
	dut.Config().New().
		WithDrivenetsText(
			`system name DAL01
 			`).
		Append(t)
	//
}

func Test_TC_14_1(t *testing.T) {
	// Test session timeout
	command := ""
	cli := ondatra.DUT(t, "dal01").CLI()
	//
	command = "set cli-terminal-length 0"
	fmt.Printf("Executing command %s\n", command)
	result := cli.Run(t, command)
	fmt.Println(result)
	// show system uptime
	command = "show system uptime"
	fmt.Printf("Executing command %s\n", command)
	result = cli.Run(t, command)
	fmt.Println(result)
	// go idle for 10 hours
	fmt.Printf("Going to wait for 10 seconds\n")
	time.Sleep(10 * time.Second)
	fmt.Println("Done waiting!")
	// show system uptime after sleep
	command = "show system uptime"
	fmt.Printf("Executing command %s\n", command)
	result = cli.Run(t, command)
	fmt.Println(result)
}

func Test_TC_14_2(t *testing.T) {
	// Test session timeout
	command := ""
	cli := ondatra.DUT(t, "dal01").CLI()
	//
	command = "set cli-terminal-length 0"
	fmt.Printf("Executing command %s\n", command)
	result := cli.Run(t, command)
	fmt.Println(result)
	// show system uptime
	command = "show system uptime"
	fmt.Printf("Executing command %s\n", command)
	result = cli.Run(t, command)
	fmt.Println(result)
	// go idle for 10 hours
	fmt.Printf("Going to wait for 10 hours\n")
	time.Sleep(10 * time.Hour)
	fmt.Println("Done waiting!")
	// show system uptime after sleep
	command = "show system uptime"
	fmt.Printf("Executing command %s\n", command)
	result = cli.Run(t, command)
	fmt.Println(result)
}

func Test_TC_14_3(t *testing.T) {
	// Test session timeout reconfigure and run CLI
	dut := ondatra.DUT(t, "dal01")
	// set session timeout to infinity
	dut.Config().New().
		WithDrivenetsText(
			`system login session-timeout 0
 			`).
		Append(t)
		//
	command := ""
	cli := ondatra.DUT(t, "dal01").CLI()
	//
	command = "set cli-terminal-length 0"
	fmt.Printf("Executing command %s\n", command)
	result := cli.Run(t, command)
	fmt.Println(result)
	// show system uptime
	command = "show system uptime"
	fmt.Printf("Executing command %s\n", command)
	result = cli.Run(t, command)
	fmt.Println(result)
	// go idle for 1 hour
	fmt.Printf("Going to wait for 1 hours\n")
	time.Sleep(1 * time.Hour)
	fmt.Println("Done waiting!")
	// show system uptime after sleep
	command = "show system uptime"
	fmt.Printf("Executing command %s\n", command)
	result = cli.Run(t, command)
	fmt.Println(result)
}

func Test_TC_14_4(t *testing.T) {
	// Test session timeout reconfigure and run CLI
	dut := ondatra.DUT(t, "dal01")
	// set session timeout to infinity
	dut.Config().New().
		WithDrivenetsText(
			`system login session-timeout 0
 			`).
		Append(t)
	//
}

func Test_TC_14_5(t *testing.T) {
	command := ""
	cli := ondatra.DUT(t, "dal01").CLI()
	//
	command = "set cli-terminal-length 0"
	fmt.Printf("Executing command %s\n", command)
	result := cli.Run(t, command)
	fmt.Println(result)
	// show system uptime
	command = "show system uptime"
	fmt.Printf("Executing command %s\n", command)
	result = cli.Run(t, command)
	fmt.Println(result)
	// go idle for 1 hour
	fmt.Printf("Going to wait for 1 hours\n")
	time.Sleep(1 * time.Hour)
	fmt.Println("Done waiting!")
	// show system uptime after sleep
	command = "show system uptime"
	fmt.Printf("Executing command %s\n", command)
	result = cli.Run(t, command)
	fmt.Println(result)
}
