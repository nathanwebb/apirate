# apirate
Add a REST API to any command

* security - how to authenticate REST APIs, does this override command permissions
* tainted inputs


Base API Design:

POST /keys/ - generate a new key. Request options: type (default: ssh), description: (default: none), name (default: type+' key ' + id)
 - creates a new key. For type=ssh, the private key will be stored securily and the public key will be returned by the API.
GET /keys - get all of the keys (not including the private keys!). This will return an array of key objects (incl. public key)
 - filter on type or name
GET /keys:id - get the key specified by the id




functions:
main.go - start the api

GET /results/?command=ping&dgexname=X&ping=deviceip
GET /results/?command=probeSnmpTests&

read in a config file, linking api with command
 - ssh connect@192.168.188.69 ssh {{dgexname}} ping {{deviceIP}}
 - ssh connect@192.168.188.69 ssh {{dgexname}} /traverse/utils/probeSnmpTests.pl --host={{deviceIP}} --community={{commstring}}


{
    name: 'ping',
    exec: "ssh connect@192.168.188.69 ssh {{dgexname}} ping {{deviceIP}}"
}, {
    name: 'probeSnmpTests',
    exec: ssh connect@192.168.188.69 ssh {{dgexname}} /traverse/utils/probeSnmpTests.pl --host={{deviceIP}} --community={{commstring}} --version={{snmpversion}} {{if eq snmpversion 3}}--authproto{{authproto}} --privprot{{privproto}}{{end}}
}


[nwebb@imcdgex39 ~]$ /traverse/utils/probeSnmpTests.pl
ERROR: no device name or input file name has been provided

  usage: probeSnmpTests.pl --host=<fqdn|ip_address>
         [ --community=<community_string> ]      [ --version=<1|2|3> ]
         [ --authproto=<none|md5|sha> ] [ --privproto=<none|des|aes> ]
         [ --type=<windows|unix|router|switch|firewall|slb|unknown>  ]
         [ --runtime=<seconds> ] [ --help ] [ remote_execution_options ]

  --host      = host name or ip address of device to probe
  --community = snmp community string (user:password:secret for version=3)
  --version   = snmp version supported by device
  --authproto = snmp agent authentication protocol        (only version=3)
  --privproto = snmp agent privary protocol               (only version=3)
  --type      = type of device being probed
  --runtime     = amount of time to run before terminating (def. 900 seconds)
  --help      = print this help message

  for remote execution via internal communication bus:

  --remote    = perform discovery on specified remote dge/dge extension
  --username  = login id with superuser privileges    (required parameter)
  --password  = password for specified login user     (required parameter)
  --endpoint  = fqdn or ip address of web application (default=localhost)
