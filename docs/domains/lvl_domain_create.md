# lvl domain create

Create a new domain

```
lvl domain create [flags]
```

### Options

```
  -a, --action string             Specify the action you want to commit
      --externalInfo string       Required when billableItemInfo for an organisation exist in db
  -n, --name string               the name of the domain (REQUIRED)
  -l, --licensee int              The unique identifier of a domaincontact with type licensee (REQUIRED)
      --organisation int          the organisation of the domain (REQUIRED)
      --nameserver1 string        Nameserver
      --nameserver2 string        Nameserver
      --nameserver3 string        Nameserver
      --nameserver4 string        Nameserver
      --nameserverIp1 string      IP address for nameserver
      --nameserverIp2 string      IP address for nameserver
      --nameserverIp3 string      IP address for nameserver
      --nameserverIp4 string      IP address for nameserver
      --nameserverIpv61 string    IPv6 address for nameserver
      --nameserverIpv62 string    IPv6 address for nameserver
      --nameserverIpv63 string    IPv6 address for nameserver
      --nameserverIpv64 string    IPv6 address for nameserver
      --ttl int                   Time to live: amount of time (in seconds) the DNS-records stay in the cache (default 28800)
      --eppCode string            eppCode
      --handleDns                 should dns be handled by lvl27 (default true)
      --extra fields string       extra fields (json, non-editable)
      --domaincontactOnsite int   the unique id of a domaincontact with type onsite
      --autoTeams string          a csv list of team id's
  -h, --help                      help for create
```

### Options inherited from parent commands

```
      --apikey string   API key
      --config string   config file (default is $HOME/.lvl.yaml)
  -o, --output string   Specifies output mode for commands. Accepted values are text or json. (default "text")
      --trace           Do detailed network request logging
```

### SEE ALSO

* [lvl domain](lvl_domain.md)	 - Commands for managing domains

