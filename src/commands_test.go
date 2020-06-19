package main

import (
	"fmt"
	"testing"
)

type testCases struct {
	args     map[string][]string
	cmd      command
	failing  bool
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
		failing:  false,
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
		failing:  true,
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
		failing:  true,
		expected: []string{},
	}, {
		args: map[string][]string{
			"name":   []string{"command with spaces"},
			"string": []string{"Hello, World!"},
		},
		cmd: command{
			Name:   "command with spaces",
			Cmd:    "echo",
			Params: "'{{.string}}'",
		},
		failing:  false,
		expected: []string{"Hello, World!"},
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
			t.Errorf("wrong args from parseArgs. input %s (length: %d), got %s (length: %d), expected %s (length: %d)", c.args, len(c.args), args, len(args), c.expected, len(c.expected))
		}
		for k, v := range args {
			if v != c.expected[k] {
				t.Errorf("wrong arg. Expected %s, got %s", c.expected[k], v)
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
		if !c.failing && cmd.Stderr != "" {
			t.Errorf("got error from command: %s", cmd.Stderr)
		}
		if c.failing && cmd.Stderr == "" {
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

func TestCheckForTaints(t *testing.T) {
	type testcase struct {
		in        map[string][]string
		untainted bool
	}
	cases := []testcase{
		{
			in: map[string][]string{
				"a": []string{"local ping", "b", "a.b"},
			},
			untainted: true,
		}, {
			in: map[string][]string{
				"a": []string{"a", "b", "a&b"},
			},
			untainted: false,
		},
	}
	for _, c := range cases {
		err := checkForTaints(c.in)
		if c.untainted && err != nil {
			t.Error(err.Error())
		}
	}
}

func TestSplitUnquotedSpace(t *testing.T) {
	cases := []string{"asd", "sdf sdf", "sdflj 'sdfkj sdklj' sdf"}
	expected := [][]string{{"asd"}, {"sdf", "sdf"}, {"sdflj", "sdfkj sdklj", "sdf"}}
	for i, c := range cases {
		result := splitUnquotedSpace(c)
		if len(expected[i]) != len(result) {
			t.Fatalf("wrong args from splitUnquoted. input %s, got %s (length: %d), expected %s (length: %d)", c, result, len(result), expected[i], len(expected[i]))
		}
		for k, v := range result {
			if v != expected[i][k] {
				t.Errorf("wrong arg. Expected %s, got %s", expected[i][k], v)
			}
		}
	}
}
