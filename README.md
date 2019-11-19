# apirate
Add a REST API to any command.

Apirate can be used to add a REST API to either local commands, or remote commands via SSH. 

## Roadmap
* Authentication/authorization
* better responses, especially error messages and status codes in json responses
* https
* mongodb storage - to allow for shared credentials
* ability to set environment variables

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
Build your container using a Dockerfile

```



2. Stand-alone


* security - how to authenticate REST APIs, does this override command permissions
* tainted inputs - DONE. Inputs are escaped before passing onto commands


Base API Design:

POST /keys/ - generate a new key. Request options: type (default: ssh), description: (default: none), name (default: type+' key ' + id)
 - creates a new key. For type=ssh, the private key will be stored securily and the public key will be returned by the API.
GET /keys - get all of the keys (not including the private keys!). This will return an array of key objects (incl. public key)
 - filter on type or name
GET /keys:id - get the key specified by the id

GET /results/?command=ping&dgexname=X&ping=deviceip
GET /results/?command=probeSnmpTests&

{
    name: "ping",
    method: "ssh",
    exec: "ssh connect@192.168.188.69 ssh {{dgexname}} ping {{deviceIP}}"
}, {
    name: "probeSnmpTests",
    exec: "ssh connect@192.168.188.69 ssh {{dgexname}} /traverse/utils/probeSnmpTests.pl --host={{deviceIP}} --community={{commstring}} --version={{snmpversion}} {{if eq snmpversion 3}}--authproto{{authproto}} --privprot{{privproto}}{{end}}"
}, {
    name: "writeConnect",
    exec: "ssh connect@192.168.188.69 /scripts/connect/appendToConf.sh {{devicename}} {{deviceip}} {{password}} {{enable}}"
}
