# apirate
Add a REST API to any command.

Apirate can be used to add a REST API to either local commands, or remote commands via SSH. 

## Roadmap
* Authentication/authorization
* better responses, especially error messages and status codes in json responses
* https
* mongodb storage - to allow for shared credentials
* ability to set environment variables for commands

## Usage

To get started, you will need to create a commands file. This file is a JSON document listing the commands that are available to the API.

### Example: default_commandstore.json
```
[{
    "name": "remote ping",
    "cmd":  "hostname",
    "user": "nathan",
    "host": "192.168.2.6"
}, {
    "name": "local ping",
    "cmd": "ping",
    "params": "-c 4 {{.ip}}"
}, {
    "name": "quiet ping",
    "cmd": "ping",
    "params": "-q -c 4 {{.ip}}"
}]
```

For the remote (ssh-based) commands, apirate can generate a key and send you the public key. Copy the public key to the remote server. The private key will be stored securely by apirate.

Then just start apirate, and send it some commands:

```
curl "http://localhost:8080/api/v1/results?name=remote%20ping&ip=192.168.2.6"
```

## Installation

there are three ways that it can be used.

1. Inside a docker container.
```
docker run --rm -p 8080:8080 crasily/apirate
```

2. Stand-alone
This can be accomplished by downloading the repo from github, building with Go-lang, and then running the installation script.


## Environment variables
There are two json files that are managed by Apirate - the keystore and the commandstore. Environment variables can be used to set the locations.

* KEYSTORE - Default is "file:///var/local/apirate/default_keystore.json"
* COMMANDSTORE - Default is "file:///var/local/apirate/default_commandstore.json"

## Notes


* security - how to authenticate REST APIs, does this override command permissions
* tainted inputs - DONE. Inputs are escaped before passing onto commands


