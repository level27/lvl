# lvl app component restores download

Download the restore file.

```
lvl app component restores download [appname] [restoreID] [flags]
```

### Examples

```
lvl app component restore download MyAppName 4123
```

### Options

```
  -f, --filename string   The name of the downloaded file.
  -h, --help              help for download
```

### Options inherited from parent commands

```
      --apikey string   API key
      --config string   config file (default is $HOME/.lvl.yaml)
  -o, --output string   Specifies output mode for commands. Accepted values are text or json. (default "text")
      --trace           Do detailed network request logging
```

### SEE ALSO

* [lvl app component restores](lvl_app_component_restores.md)	 - Command to manage restores on an app.

