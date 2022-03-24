## lvl system checks update

update a specific check from a system

```
lvl system checks update [SystemID] [CheckID] [flags]
```

### Options

```
  -h, --help                     help for update
  -p, --parameters stringArray   Add custom parameters for a check. Usage -> SINGLE PAR: [ -p waf=true ], MULTIPLE PAR: [ -p waf=true -p timeout=200 ], MULTIPLE VALUES: [ -p versions=''7, 5.4'']
```

### Options inherited from parent commands

```
      --apikey string   API key
      --config string   config file (default is $HOME/.lvl.yaml)
  -o, --output string   Specifies output mode for commands. Accepted values are text or json. (default "text")
      --trace           Do detailed network request logging
```

### SEE ALSO

* [lvl system checks](lvl_system_checks.md)	 - Manage systems checks

