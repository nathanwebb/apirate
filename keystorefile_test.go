package main

import "testing"

func TestLoadKeysFromFile(t *testing.T) {
	keyfile := "keys_test.json"
	keys, err := loadKeysFromFile(keyfile)
	if err != nil {
		t.Error(err.Error())
	}
	if len(keys) == 0 || keys[0].ID != 1 {
		t.Errorf("failed to load keys")
	}
}

func TestSaveKeyToFile(t *testing.T) {
	keyfile := "keys_test.json"
	keys, err := loadKeysFromFile(keyfile)
	originalKeysLen := len(keys)
	err = saveKeysToFile(keyfile, keys)
	if err != nil {
		t.Error(err.Error())
	}
	keys, err = loadKeysFromFile(keyfile)
	if len(keys) <= originalKeysLen {
		t.Errorf("failed to save new key")
	}
}
