# lvl app update

Update an app.

```
lvl app update [appID] [flags]
```

### Examples

```
lvl app update 2067 --name myUpdatedName
```

### Options

```
      --autoTeams strings     A csv list of team ID's.
  -h, --help                  help for update
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

