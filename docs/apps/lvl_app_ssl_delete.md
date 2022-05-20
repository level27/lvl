# lvl app ssl delete

Delete an SSL certificate from an app

```
lvl app ssl delete [app] [SSL cert] [flags]
```

### Examples

```
lvl app ssl delete forum forum.example.com
```

### Options

```
      --force   Do not ask for confirmation to delete the SSL certificate
  -h, --help    help for delete
```

### Options inherited from parent commands

```
      --apikey string   API key
      --config string   config file (default is $HOME/.lvl.yaml)
  -o, --output string   Specifies output mode for commands. Accepted values are text or json. (default "text")
      --trace           Do detailed network request logging
```

### SEE ALSO

* [lvl app ssl](lvl_app_ssl.md)	 - Commands for managing SSL certificates on apps

