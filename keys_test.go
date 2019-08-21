package main

import "testing"

func TestCreateSSHKey(t *testing.T) {
	sshkey := key{
		Name: "test",
	}
	sshkey, err := createSSHKey(sshkey)
	if err != nil {
		t.Errorf(err.Error())
	}
	if sshkey.Type != "ssh" {
		t.Errorf("failed to get key: %+v", sshkey)
	}
}
