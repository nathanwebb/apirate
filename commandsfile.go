package main

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

func loadCommandsFromFile(source string) ([]command, error) {
	commands := []command{}
	filepath := filepath.Join(".", source)
	f, err := os.OpenFile(filepath, os.O_RDONLY|os.O_CREATE, 0600)
	defer f.Close()
	if err != nil {
		return commands, err
	}
	err = json.NewDecoder(f).Decode(&commands)
	if err != nil && err != io.EOF {
		return commands, err
	}
	return commands, nil
}
