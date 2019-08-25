package main

import (
	"encoding/json"
	"errors"
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

func deleteAllKeysFromFile(keyfile string) error {
	existingKeys, err := loadKeysFromFile(keyfile)
	for _, k := range existingKeys {
		err = os.Remove(k.PrivateKeyFilename)
		if err != nil && !os.IsNotExist(err) {
			return errors.New(fmt.Sprintf("failed to remove private key %s: %s", k.PrivateKeyFilename, err.Error()))
		}
	}
	err = os.Remove(keyfile)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to remove keystore %s: %s", keyfile, err.Error()))
	}
	return nil
}

func deleteKeyFromFile(keyfile string, id int) error {
	existingKeys, err := loadKeysFromFile(keyfile)
	n := 0
	for _, k := range existingKeys {
		if k.ID == id {
			err = os.Remove(k.PrivateKeyFilename)
			if err != nil && !os.IsNotExist(err) {
				return errors.New(fmt.Sprintf("failed to remove private key %s: %s", k.PrivateKeyFilename, err.Error()))
			}
		} else {
			existingKeys[n] = k
			n++
		}
	}
	existingKeys = existingKeys[:n]
	return saveKeysToFile(keyfile, existingKeys)
}
