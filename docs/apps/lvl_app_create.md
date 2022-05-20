# lvl app create

Create a new app.

```
lvl app create [flags]
```

### Examples

```
lvl app create -n myNewApp --organisation level27
```

### Options

```
      --autoTeams ints        A csv list of team ID's.
      --externalInfo string   ExternalInfo (required when billableItemInfo entities for an organisation exist in DB.)
  -h, --help                  help for create
  -n, --name string           Name of the app.
      --organisation string   The name of the organisation/owner of the app.
```

### Options inherited from parent commands

```
      --apikey string   API key
      --config string   config file (default is $HOME/.lvl.yaml)
  -o, --output string   Specifies output mode for commands. Accepted values are text or json. (default "text")
      --trace           Do detailed network request logging
```

### SEE ALSO

* [lvl app](lvl_app.md)	 - Commands to manage apps

