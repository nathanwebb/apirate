package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func loadKeysFromFile(keyfile string) ([]key, error) {
	keys := []key{}
	keypath := filepath.Join(".", keyfile)
	f, err := os.OpenFile(keypath, os.O_RDONLY|os.O_CREATE, 0600)
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
		return keys, err
	}
	err = json.NewDecoder(f).Decode(&keys)
	if err == io.EOF {
		return keys, nil
	}
	return keys, err
}

func saveKeysToFile(keyfile string, keys []key) error {
	keypath := filepath.Join(".", keyfile)
	jsonfile, err := json.MarshalIndent(keys, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(keypath)
	return ioutil.WriteFile(keypath, jsonfile, 0644)
}
