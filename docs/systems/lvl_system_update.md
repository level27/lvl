## lvl system update

Update settings on a system

```
lvl system update [flags]
```

### Options

```
      --cpu int                      Set amount of CPU cores of the system
  -h, --help                         help for update
      --installSecurityUpdates int   Set security updates mode index
      --limitRiops int               Set read IOPS limit
      --limitWiops int               Set write IOPS limit
      --managementType string        Set management type of the system
      --memory int                   Set amount of memory in GB of the system
      --name string                  New name for this system
      --organisation string          Set organisation that owns this system. Can be both a name or an ID
      --publicNetworking int         
      --remarks string               
  -s, --settings-file string         JSON file to read settings from. Pass '-' to read from stdin.
```

### Options inherited from parent commands

```
      --apikey string   API key
      --config string   config file (default is $HOME/.lvl.yaml)
  -o, --output string   Specifies output mode for commands. Accepted values are text or json. (default "text")
      --trace           Do detailed network request logging
```

### SEE ALSO

* [lvl system](lvl_system.md)	 - Commands for managing systems

