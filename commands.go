package main

type command struct {
	Exec   string `json:"command"`
	Params string `json:"params"`
}

func loadCommands(source string) ([]command, error) {
	return "", nil
}
