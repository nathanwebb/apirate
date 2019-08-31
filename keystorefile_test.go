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
	if len(keys) == 0 || keys[0].ID != "1" {
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
		ID:                 "2",
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

func TestDeleteAllKeys(t *testing.T) {
	keyfile := "keys_test.json"
	keys := []key{{
		ID:                 "2",
		Type:               "ssh",
		PublicKey:          "rsa-ssh...",
		PrivateKeyFilename: "id_rsa_test_3",
	}}
	err := saveKeysToFile(keyfile, keys)
	if err != nil {
		t.Error(err.Error())
	}
	err = deleteAllKeysFromFile(keyfile)
	if err != nil {
		t.Errorf("failed to delete keys: %s", err.Error())
	}
	_, err = os.Open(keyfile)
	if err == nil || !os.IsNotExist(err) {
		t.Errorf("able to remove keystore: %s", err.Error())
	}
}

func TestDeleteKey(t *testing.T) {
	keyfile := "keys_test_single_delete.json"
	keys := []key{{
		ID:                 "2",
		Type:               "ssh",
		PublicKey:          "rsa-ssh...",
		PrivateKeyFilename: "id_rsa_test_4",
	}, {
		ID:                 "3",
		Type:               "ssh",
		PublicKey:          "rsa-ssh...",
		PrivateKeyFilename: "id_rsa_test_5",
	}}
	for _, k := range keys {
		_, err := os.Create(k.PrivateKeyFilename)
		if err != nil {
			t.Errorf("error creating dummy private key file: %s", err.Error())
		}
	}
	err := saveKeysToFile(keyfile, keys)
	if err != nil {
		t.Error(err.Error())
	}
	err = deleteKeyFromFile(keyfile, "2")
	if err != nil {
		t.Errorf("failed to delete key: %s", err.Error())
	}
	keys, err = loadKeysFromFile(keyfile)
	if err != nil {
		t.Error(err.Error())
	}
	os.Remove(keyfile)
	for _, k := range keys {
		os.Remove(k.PrivateKeyFilename)
	}
	if len(keys) > 1 {
		t.Errorf("failed to delete key from file. Still %d keys", len(keys))
	}
	if len(keys) == 1 && keys[0].ID == "2" {
		t.Errorf("deleted wrong key")
	}

}
