package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os/exec"
	"regexp"
	"strings"
	"text/template"

	"github.com/davecgh/go-spew/spew"
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
	Error    string
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
	log.Printf("%+v\n", cmd)
	spew.Dump(queryArgs)
	err := checkForTaints(queryArgs)
	if err != nil {
		return cmd, err
	}
	args, err := parseArgs(cmd, queryArgs)
	log.Printf("%+v\n", args)
	if err != nil && !strings.Contains(err.Error(), "map has no entry for key") {
		return command{}, err
	}
	if cmd.Host == "" {
		return execLocalCommand(cmd, args)
	}
	return execRemoteCommand(cmd, args)
}

func checkForTaints(queryArgs map[string][]string) error {
	var err error
	pattern := regexp.MustCompile(`[^\sa-zA-Z0-9,._+:@%/!-]`)
	for k, a := range queryArgs {
		for _, s := range a {
			if pattern.MatchString(s) {
				return errors.New(k + " has tainted variable " + s)
			}
		}
	}
	return err
}

func execLocalCommand(cmd command, args []string) (command, error) {
	cmdToRun := exec.Command(cmd.Cmd, args...)
	var stdout, stderr bytes.Buffer
	cmdToRun.Stdout = &stdout
	cmdToRun.Stderr = &stderr
	err := cmdToRun.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			cmd.ExitCode = exitError.ExitCode()
		}
	}
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
	defer conn.Close()
	session, err := conn.NewSession()
	defer session.Close()
	if err != nil {
		return command{}, err
	}
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	if err := session.RequestPty("xterm", 40, 80, modes); err != nil {
		log.Fatal("request for pseudo terminal failed: ", err)
	}

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr
	log.Println(cmd.Cmd + " " + strings.Join(args, " "))
	err = session.Run(cmd.Cmd + " " + strings.Join(args, " "))
	if err != nil {
		if exitError, ok := err.(*ssh.ExitError); ok {
			cmd.ExitCode = exitError.ExitStatus()
		}
	}
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
	log.Println(len(queryArgs))
	//quotedArgs := quoteArgs(queryArgs)
	//flatArgs := flatten(quotedArgs)
	flatArgs := flatten(queryArgs)
	spew.Dump(flatArgs)

	t := template.Must(template.New("t2").Parse(cmd.Params))
	t.Option("missingkey=error")
	err := t.Execute(buf, flatArgs)
	log.Printf("%v\n", buf)
	args := splitUnquotedSpace(buf.String())
	return args, err
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

func quoteArgs(queryArgs map[string][]string) map[string][]string {
	quoted := make(map[string][]string, len(queryArgs))
	for k, v := range queryArgs {
		quotedArgs := []string{}
		for _, s := range v {
			quotedArgs = append(quotedArgs, shellescape.Quote(s))
		}
		quoted[k] = quotedArgs
	}
	return quoted
}

func splitUnquotedSpace(s string) []string {
	var quoteRune rune
	result := []string{}
	currentWord := ""
	for _, r := range s {
		if quoteRune == 0 && isSpace(r) {
			result = appendWord(result, currentWord)
			currentWord = ""
			continue
		} else if r == quoteRune {
			quoteRune = 0
		} else if quoteRune == 0 && isQuote(r) {
			quoteRune = r
		}
		currentWord = currentWord + string(r)
	}
	result = appendWord(result, currentWord)
	return result
}

func appendWord(r []string, word string) []string {
	if len(word) > 0 {
		return append(r, word)
	}
	return r
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

func isQuote(r rune) bool {
	return r == '\'' || r == '"'
}
