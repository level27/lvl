# lvl mail forwarder update

Update settings on a mail forwarder

```
lvl mail forwarder update [mailgroup] [mail forwarder] [flags]
```

### Options

```
      --destination string     Comma-separated list of all destinations for this forwarder
  -h, --help                   help for update
  -s, --settings-file string   JSON file to read settings from. Pass '-' to read from stdin.
```

### Options inherited from parent commands

```
      --apikey string   API key
      --config string   config file (default is $HOME/.lvl.yaml)
  -o, --output string   Specifies output mode for commands. Accepted values are text or json. (default "text")
      --trace           Do detailed network request logging
```

### SEE ALSO

* [lvl mail forwarder](lvl_mail_forwarder.md)	 - Commands for managing mail forwarders

