# lvl app migrations action

Execute an action for a migration

### Examples

```
lvl app migrations action deny MyAppName 241
lvl app migrations action restart MyAppName 234
```

### Options

```
  -h, --help   help for action
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
* [lvl app migrations action confirm](lvl_app_migrations_action_confirm.md)	 - Execute confirm action on a migration
* [lvl app migrations action deny](lvl_app_migrations_action_deny.md)	 - Execute confirm action on a migration
* [lvl app migrations action retry](lvl_app_migrations_action_retry.md)	 - Execute confirm action on a migration

