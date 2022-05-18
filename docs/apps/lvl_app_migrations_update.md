# lvl app migrations update

Update an app migration.

```
lvl app migrations update [appID] [migrationID] [flags]
```

### Examples

```
lvl app migrations update MyAppName 3414
```

### Options

```
  -h, --help             help for update
      --planned string   DateTime - timestamp.
  -t, --type string      Migration type. (one of automatic (all migration steps are done automatically), confirmed (a user has to confirm each migration step)).
```

### Options inherited from parent commands

```
      --apikey string   API key
      --config string   config file (default is $HOME/.lvl.yaml)
  -o, --output string   Specifies output mode for commands. Accepted values are text or json. (default "text")
      --trace           Do detailed network request logging
```

### SEE ALSO

* [lvl app migrations](lvl_app_migrations.md)	 - Commands to manage app migrations.

