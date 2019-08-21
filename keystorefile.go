package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func loadKeysFromFile(keyfile string) ([]key, error) {
	keys := []key{}
	keypath := filepath.Join(".", keyfile)
	f, err := os.Open(keypath)
	if err != nil {
		return keys, err
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&keys)
	return keys, err
}

func saveKeysToFile(keystore string) error {
	return nil
}
