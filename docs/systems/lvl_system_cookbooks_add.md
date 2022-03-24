## lvl system cookbooks add

add a cookbook to a system

```
lvl system cookbooks add [systemID] [flags]
```

### Options

```
  -h, --help                     help for add
  -p, --parameters stringArray   Add custom parameters for cookbook. SINGLE PAR: [ -p waf=true ], MULTIPLE PAR: [ -p waf=true -p timeout=200 ], MULTIPLE VALUES: [ -p versions=''7, 5.4'']
  -t, --type string              Cookbook type (non-editable). Cookbook types can't repeat for one system
```

### Options inherited from parent commands

```
      --apikey string   API key
      --config string   config file (default is $HOME/.lvl.yaml)
  -o, --output string   Specifies output mode for commands. Accepted values are text or json. (default "text")
      --trace           Do detailed network request logging
```

### SEE ALSO

* [lvl system cookbooks](lvl_system_cookbooks.md)	 - Manage systems cookbooks

