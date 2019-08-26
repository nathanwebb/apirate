package main

import (
	"bytes"
	"errors"
	"html/template"
	"net/url"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
)

type command struct {
	Name     string `json:"name"`
	Cmd      string `json:"cmd"`
	User     string `json:"user"`
	Host     string `json:"host"`
	Params   string `json:"params"`
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	ExitCode int    `json:"exitcode"`
}

func loadCommands(source string) ([]command, error) {
	f, _ := url.ParseRequestURI(source)
	return loadCommandsFromFile(f.Path)
}

func getCommandForRequest(c *gin.Context, commands []command) (command, error) {
	requestedCommand := c.Query("name")
	for _, cmd := range commands {
		if cmd.Name == requestedCommand {
			return cmd, nil
		}
	}
	return command{}, errors.New("invalid command name")
}

func execCommand(cmd command, queryArgs map[string][]string) (command, error) {
	args, err := parseArgs(cmd, queryArgs)
	cmdToRun := exec.Command(cmd.Cmd, args...)
	var stdout, stderr bytes.Buffer
	cmdToRun.Stdout = &stdout
	cmdToRun.Stderr = &stderr
	err = cmdToRun.Run()
	cmd.Stdout = stdout.String()
	cmd.Stderr = stderr.String()
	return cmd, err
}

func parseArgs(cmd command, queryArgs map[string][]string) ([]string, error) {
	buf := new(bytes.Buffer)
	flatArgs := flatten(queryArgs)
	t := template.Must(template.New("t2").Parse(cmd.Params))
	err := t.Execute(buf, flatArgs)
	return strings.Fields(buf.String()), err
}

func flatten(queryArgs map[string][]string) map[string]interface{} {
	flat := make(map[string]interface{}, len(queryArgs))
	for k, v := range queryArgs {
		if len(v) == 1 {
			flat[k] = v[0]
		} else {
			flat[k] = v
		}
	}
	return flat
}
