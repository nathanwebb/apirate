package main

import "testing"

func TestLoadCommandsFromFile(t *testing.T) {
	testCommandsFile := "commands_config_test.json"
	commands, err := loadCommandsFromFile(testCommandsFile)
	if err != nil {
		t.Errorf("failed to load commands: %s", err.Error())
	}
	if len(commands) != 2 {
		t.Errorf("failed to load commands: only %d loaded, expected 2.", len(commands))
	}
	if commands[1].Cmd != "ping" {
		t.Errorf("expected first command to be ping, got %s", commands[1].Cmd)
	}
}
