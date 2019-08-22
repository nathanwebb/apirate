package main

import "testing"

func TestCreateSSHKey(t *testing.T) {
	sshkey := key{}
	sshkey, err := createSSHKey(sshkey)
	if err != nil {
		t.Errorf(err.Error())
	}
	if sshkey.Type != "ssh" {
		t.Errorf("failed to get key: %+v", sshkey)
	}
	t.Logf("%+v\n", sshkey)
}

func TestLoadKeys(t *testing.T) {
	keystore := "file:///keys_test.json"
	keys, err := loadKeys(keystore)
	if err != nil {
		t.Error(err.Error())
	}
	if len(keys) == 0 || keys[0].ID != 1 {
		t.Errorf("failed to load keys")
	}
}

func TestSaveKey(t *testing.T) {

}