package main

import (
	"os"
	"testing"
)

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
	os.Setenv("KEYSTORE", "file:///keystore_test.json")
	keys, err := loadKeys()
	if err != nil {
		t.Error(err.Error())
	}
	if len(keys) == 0 || keys[0].ID != "1" {
		t.Errorf("failed to load keys")
	}
}

func TestSaveKey(t *testing.T) {

}

func TestDeleteKeys(t *testing.T) {
	keyfile := "keystore_test.json"
	os.Setenv("KEYSTORE", "file:///"+keyfile)
	keys := []key{{
		ID:                 "2",
		Type:               "ssh",
		PublicKey:          "rsa-ssh...",
		PrivateKeyFilename: "id_rsa_test_3",
	}}
	err := saveKeys(keys)
	if err != nil {
		t.Error(err.Error())
	}
	err = deleteAllKeys()
	if err != nil {
		t.Errorf("failed to delete keys: %s", err.Error())
	}
	_, err = os.Open(keyfile)
	if err == nil || !os.IsNotExist(err) {
		t.Errorf("able to remove keystore: %s", err.Error())
	}
}
