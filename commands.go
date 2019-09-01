package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os/exec"
	"strings"
	"text/template"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/ssh"
	"gopkg.in/alessio/shellescape.v1"
)

type command struct {
	Name     string `json:"name"`
	Cmd      string `json:"cmd"`
	User     string `json:"user"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Params   string `json:"params"`
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	ExitCode int    `json:"exitcode"`
}

func loadCommands(source string) ([]command, error) {
	f, err := url.Parse(source)
	if err != nil {
		return []command{}, err
	}
	return loadCommandsFromFile(f.Path)
}

func getCommandForRequest(c *gin.Context, commands []command) (command, error) {
	requestedCommand := c.Query("name")
	if requestedCommand == "" {
		return command{}, errors.New("missing 'name' query argument")
	}
	for _, cmd := range commands {
		if cmd.Name == requestedCommand {
			return cmd, nil
		}
	}
	return command{}, errors.New("invalid command name")
}

func execCommand(cmd command, queryArgs map[string][]string) (command, error) {
	args, err := parseArgs(cmd, queryArgs)
	if err != nil {
		return command{}, err
	}
	if cmd.Host == "" {
		return execLocalCommand(cmd, args)
	} else {
		return execRemoteCommand(cmd, args)
	}
}

func execLocalCommand(cmd command, args []string) (command, error) {
	cmdToRun := exec.Command(cmd.Cmd, args...)
	fmt.Printf("%+v\n", cmdToRun)
	var stdout, stderr bytes.Buffer
	cmdToRun.Stdout = &stdout
	cmdToRun.Stderr = &stderr
	err := cmdToRun.Run()
	cmd.Stdout = stdout.String()
	cmd.Stderr = stderr.String()
	return cmd, err
}

func execRemoteCommand(cmd command, args []string) (command, error) {
	config := &ssh.ClientConfig{
		User:            cmd.User,
		Auth:            makeKeyring(),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	log.Printf("%+v\n", cmd)
	log.Printf("%+v\n", config)
	if cmd.Port == "" {
		cmd.Port = "22"
	}
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", cmd.Host, cmd.Port), config)
	if err != nil {
		return command{}, err
	}
	cmdToRun, err := conn.NewSession()
	defer cmdToRun.Close()
	if err != nil {
		return command{}, err
	}

	var stdout, stderr bytes.Buffer
	cmdToRun.Stdout = &stdout
	cmdToRun.Stderr = &stderr
	log.Println(cmd.Cmd)
	err = cmdToRun.Run(cmd.Cmd + " " + strings.Join(args, " "))
	cmd.Stdout = stdout.String()
	cmd.Stderr = stderr.String()
	log.Println(cmd.Stderr)
	return cmd, err
}

func makeKeyring() []ssh.AuthMethod {
	signers := []ssh.AuthMethod{}
	keys, err := getAllPrivateKeyFilenames()
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	for _, keyfile := range keys {
		signer, err := makeSigner(keyfile)
		if err == nil {
			signers = append(signers, signer)
		}
	}
	return signers
}

func makeSigner(path string) (ssh.AuthMethod, error) {
	key, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(signer), nil
}

func parseArgs(cmd command, queryArgs map[string][]string) ([]string, error) {
	buf := new(bytes.Buffer)
	flatArgs := flatten(queryArgs)
	t := template.Must(template.New("t2").Parse(cmd.Params))
	t.Option("missingkey=error")
	err := t.Execute(buf, flatArgs)
	args := strings.Fields(buf.String())
	response := []string{}
	for _, s := range args {
		response = append(response, shellescape.Quote(s))
	}
	return response, err
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
