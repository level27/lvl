# lvl system network ip add

Add IP address to a system network

### Synopsis

Adds an IP address to a system network. Address can be either IPv4 or IPv6. The special values 'auto' and 'auto-v6' automatically fetch an unused address to use.

```
lvl system network ip add [system] [network] [address] [flags]
```

### Options

```
  -h, --help              help for add
      --hostname string   Hostname for the IP address. If not specified the system hostname is used.
```

### Options inherited from parent commands

```
      --apikey string   API key
      --config string   config file (default is $HOME/.lvl.yaml)
  -o, --output string   Specifies output mode for commands. Accepted values are text or json. (default "text")
      --trace           Do detailed network request logging
```

### SEE ALSO

* [lvl system network ip](lvl_system_network_ip.md)	 - Manage IP addresses on network connections

