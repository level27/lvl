# lvl mail create

Create a new mail group

### Synopsis

Does not automatically link any domains to the mail group. Use separate commands after the mail group has been created.

```
lvl mail create [flags]
```

### Options

```
      --externalInfo string   
  -h, --help                  help for create
      --name string           Name of the new mailgroup
      --organisation string   Organisation owning the new mailgroup
```

### Options inherited from parent commands

```
      --apikey string   API key
      --config string   config file (default is $HOME/.lvl.yaml)
  -o, --output string   Specifies output mode for commands. Accepted values are text or json. (default "text")
      --trace           Do detailed network request logging
```

### SEE ALSO

* [lvl mail](lvl_mail.md)	 - Commands to manage mailgroups and mailboxes

