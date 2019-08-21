package main

import "testing"

func TestCommandsParsing(t *testing.T) {
	commands := []struct {
		template string
		output   string
		url      string
	}{{
		template: "ssh connect@192.168.188.69 ssh {{dgexname}} ping {{deviceip}}",
		output:   "ssh connect@192.168.188.69 ssh imcdgex39 ping 127.0.0.1",
		url:      "GET /results/?command=ping&dgexname=imcdgex39&ping=127.0.0.1",
	}, {
		template: "ssh connect@192.168.188.69 ssh {{dgexname}} /traverse/utils/probeSnmpTests.pl --host={{deviceIP}} --community={{commstring}} --version={{snmpversion}} {{if eq snmpversion 3}}--authproto{{authproto}} --privprot{{privproto}}{{end}}",
		output:   "ssh connect@192.168.188.69 ssh imcdgex39 /traverse/utils/probeSnmpTests.pl --host=127.0.0.1 --community=public --version=2",
		url:      "GET /results/?command=ping&dgexname=imcdgex39&deviceip=127.0.0.1&commstring=public&snmpversion=2",
	}, {
		template: "ssh connect@192.168.188.69 ssh {{dgexname}} /traverse/utils/probeSnmpTests.pl --host={{deviceIP}} --community={{commstring}} --version={{snmpversion}} {{if eq snmpversion 3}}--authproto{{authproto}} --privprot{{privproto}}{{end}}",
		output:   "ssh connect@192.168.188.69 ssh imcdgex39 /traverse/utils/probeSnmpTests.pl --host=127.0.0.1 --community=public --version=3 --authproto=sha --privproto=aes",
		url:      "GET /results/?command=ping&dgexname=imcdgex39&deviceip=127.0.0.1&commstring=rits:pub:lic&snmpversion=3&authproto=sha&privprot=aes",
	}}

	for _, command := range commands {
		result, err := getCommandFromUrl(command.url)
		if err != nil {
			t.Errorf("error in getCommandFromUrl: %s", err.Error())
		}
		if result != command.output {
			t.Errorf("expected %s, got %s", command.output, result)
		}
	}
}
