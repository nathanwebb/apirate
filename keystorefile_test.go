package main

import (
	"os"
	"testing"
)

func TestLoadKeysFromFile(t *testing.T) {
	keyfile := "keystore_test.json"
	keys, err := loadKeysFromFile(keyfile)
	if err != nil {
		t.Error(err.Error())
	}
	if len(keys) == 0 || keys[0].ID != 1 {
		t.Errorf("failed to load keys")
	}
}

func TestLoadKeysFromMissingFile(t *testing.T) {
	keyfile := "keys_test_missing.json"
	keys, err := loadKeysFromFile(keyfile)
	if err != nil {
		t.Error(err.Error())
	}
	if len(keys) > 0 {
		t.Errorf("loaded keys from empty file")
	}
	if err = os.Remove(keyfile); err != nil {
		t.Errorf("failed to remove keystore: %s", err)
	}
}

func TestSaveKeyToFile(t *testing.T) {
	keyfile := "keys_test.json"
	keys := []key{{
		ID:                 2,
		Type:               "ssh",
		PublicKey:          "rsa-ssh...",
		PrivateKeyFilename: "id_rsa_test_2",
	}}
	originalKeysLen := len(keys)
	err := saveKeysToFile(keyfile, keys)
	if err != nil {
		t.Error(err.Error())
	}
	keys, err = loadKeysFromFile(keyfile)
	if len(keys) != originalKeysLen {
		t.Errorf("failed to save new key")
	}
	if err = os.Remove(keyfile); err != nil {
		t.Errorf("failed to remove keystore: %s", err)
	}
}
