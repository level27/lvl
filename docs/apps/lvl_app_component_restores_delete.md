# lvl app component restores delete

Delete a specific restore from an app.

```
lvl app component restores delete [flags]
```

### Examples

```
lvl app component restore delete MyAppName 4532
```

### Options

```
  -h, --help   help for delete
  -y, --yes    Set this flag to skip confirmation when deleting a check
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

