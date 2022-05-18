# lvl domain delete

Delete a domain

### Synopsis

use LVL DOMAIN DELETE <ID or ID's>. You can give multiple ID's to this command by seperating them trough whitespaces.

```
lvl domain delete [domainId] [flags]
```

### Options

```
  -h, --help   help for delete
  -y, --yes    Confirmation flag. Set this flag to delete the domain without confirmation question.
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

