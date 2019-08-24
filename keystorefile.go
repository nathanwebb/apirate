package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

func loadKeysFromFile(keyfile string) ([]key, error) {
	keys := []key{}
	keypath := filepath.Join(".", keyfile)
	f, err := os.OpenFile(keypath, os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		return keys, err
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&keys)
	return keys, err
}

func saveKeysToFile(keyfile string, keys []key) error {
	keypath := filepath.Join(".", keyfile)
	file, err := json.MarshalIndent(keys, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(keypath, file, 0644)
}
