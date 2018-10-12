# mastparse

A (Golang learning) tool to parse the mast inventory file to print out
some useful command line commands, so that you can cut and paste them into
a terminal.

## Typical usage

```bash
martin@martin-Latitude-5590:~/go/src/github.com/martin-guthrie-docker/mastparse$ ./mastparse inspect mgag
[0000]  WARN       command.go: 852:cobra.(*Command).ExecuteC     | Config File ".mastp" Not Found in "[/home/martin]"
SSH:
----
linux-ucp-manager-primary : i-0dc78492b42702fa5  : ssh -i ~./mast/id_rsa docker@54.202.38.187
 linux-dtr-worker-primary : i-0a8d5c7224a6f4716  : ssh -i ~./mast/id_rsa docker@54.214.205.186
            linux-workers : i-04f15b776742630a8  : ssh -i ~./mast/id_rsa docker@34.221.156.173

Docker UCP CLI interface: (run these cmds from ssh terminal on UCP manager
-------------------------
sudo apt install unzip jq
AUTHTOKEN=$(curl -sk -d '{"username":"admin","password":"ucp12345"}' https://mgag-ucp-f7295da70e433929.elb.us-west-2.amazonaws.com/auth/login | jq -r .auth_token)
curl -k -H "Authorization: Bearer $AUTHTOKEN" https://mgag-ucp-f7295da70e433929.elb.us-west-2.amazonaws.com/api/clientbundle -o bundle.zip
unzip bundle.zip
eval "$(<env.sh)"
```

## Help

```bash
martin@martin-Latitude-5590:~/go/src/github.com/martin-guthrie-docker/mastparse$ ./mastparse
This tool provides insight into a mast deployment

Usage:
  mastparse [command]

Available Commands:
  help        Help about any command
  inspect     parse deployment 'name' to the console
  state       print state, configuration to the console/log
  version     Version and Release information

Flags:
      --config string     config file (default is $HOME/.mastp, ./.mastp.yaml)
  -h, --help              help for mastparse
      --mastpath string   path to mast data storage (default "/home/martin/.mast")
      --verbose           Set debug level

Use "mastparse [command] --help" for more information about a command.

```