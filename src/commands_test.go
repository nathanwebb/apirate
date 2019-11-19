package main

import (
	"fmt"
	"testing"
)

type testCases struct {
	args     map[string][]string
	cmd      command
	expected []string
}

var cases = []testCases{
	{
		args: map[string][]string{
			"name": []string{"local ping"},
			"ip":   []string{"127.0.0.1"},
		},
		cmd: command{
			Name:   "local ping",
			Cmd:    "ping",
			Params: "-c 4 {{.ip}}",
		},
		expected: []string{"-c", "4", "127.0.0.1"},
	}, {
		args: map[string][]string{
			"name": []string{"local ping"},
			"ip":   []string{"127.0.0.1"},
		},
		cmd: command{
			Name:   "local ping",
			Cmd:    "ping",
			Params: " -c 4",
		},
		expected: []string{"-c", "4"},
	}, {
		args: map[string][]string{
			"name": []string{"local ping"},
		},
		cmd: command{
			Name:   "local ping",
			Cmd:    "ping",
			Params: "{{.ip}} -c 4",
		},
		expected: []string{},
	},
}

func TestCommandsParsing(t *testing.T) {
	commands := []struct {
		template string
		output   string
		url      string
	}{{
		template: "ssh connect@192.168.1.9 ssh {{dgexname}} ping {{deviceip}}",
		output:   "ssh connect@192.168.1.9 ssh dgex1 ping 127.0.0.1",
		url:      "GET /results/?command=ping&dgexname=dgex1&ping=127.0.0.1",
	}, {
		template: "ssh connect@192.168.1.9 ssh {{dgexname}} snmp.pl --host={{deviceIP}} --community={{commstring}} --version={{snmpversion}} {{if eq snmpversion 3}}--authproto{{authproto}} --privprot{{privproto}}{{end}}",
		output:   "ssh connect@192.168.1.9 ssh dgex1 snmp.pl --host=127.0.0.1 --community=public --version=2",
		url:      "GET /results/?command=snmp.pl&dgexname=dgex1&deviceip=127.0.0.1&commstring=public&snmpversion=2",
	}, {
		template: "ssh connect@192.168.1.9 ssh {{dgexname}} snmp.pl --host={{deviceIP}} --community={{commstring}} --version={{snmpversion}} {{if eq snmpversion 3}}--authproto{{authproto}} --privprot{{privproto}}{{end}}",
		output:   "ssh connect@192.168.1.9 ssh dgex1 snmp.pl --host=127.0.0.1 --community=public --version=3 --authproto=sha --privproto=aes",
		url:      "GET /results/?command=snmp.pl&dgexname=dgex1&deviceip=127.0.0.1&commstring=apirate:pub:lic&snmpversion=3&authproto=sha&privprot=aes",
	}}
	fmt.Println(commands)
}

func TestParseArgs(t *testing.T) {
	for _, c := range cases {
		args, err := parseArgs(c.cmd, c.args)
		if len(c.expected) > 0 && err != nil {
			t.Errorf("error parsing query args: %s", err.Error())
		}
		if len(args) != len(c.expected) {
			t.Errorf("wrong args from parseArgs. Got %s (length: %d), expected %s (length: %d)", args, len(args), c.expected, len(c.expected))
		}
		for k, v := range args {
			if v != c.expected[k] {
				t.Errorf("wrong arg. Expected %s, got %s", v, c.expected[k])
			}
		}
	}
}

func TestFlatten(t *testing.T) {
	inputs := []map[string][]string{
		map[string][]string{
			"name": []string{"test1"},
			"args": []string{"just one"},
		},
		map[string][]string{
			"name": []string{"test2"},
			"args": []string{"this one", "and another"},
		},
	}
	for idx, i := range inputs {
		flat := flatten(i)
		if len(flat) != 2 {
			t.Errorf("wrong length back. Expected 2, got %d", len(flat))
		}
		fmt.Println(flat["args"])
		if idx == 0 && flat["args"] != "just one" {
			t.Errorf("expected string, got %+v\n", flat["args"])
		}
		if idx == 1 && len(flat["args"].([]string)) != 2 {
			t.Errorf("failed")
		}
	}
}

func TestExecCommand(t *testing.T) {
	for idx, c := range cases {
		cmd, err := execCommand(c.cmd, c.args)
		if idx == 0 && err != nil {
			t.Errorf(err.Error())
		}
		if idx == 0 && cmd.Stderr != "" {
			t.Errorf("got error from command: %s", cmd.Stderr)
		}
		if idx > 0 && cmd.Stderr == "" {
			t.Errorf("failed to get error from command number: %d", idx)
		}
	}
}

func TestLoadCommands(t *testing.T) {
	uris := []string{"commands_config_test.json", "commands_config_test.json"}
	for _, u := range uris {
		cmds, err := loadCommands(u)
		if err != nil {
			t.Error(err.Error())
		}
		if len(cmds) == 0 {
			t.Error("failed to load any commands")
		}
		if cmds[0].Name != "local ping" {
			t.Errorf("wrong name. expected 'local ping', got %s", cmds[0].Name)
		}
	}
}
