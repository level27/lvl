## lvl system volume update

Update settings on a volume

```
lvl system volume update [system] [volume] [flags]
```

### Options

```
      --autoResize             New autoResize setting
  -h, --help                   help for update
      --name string            New name for the volume
  -s, --settings-file string   JSON file to read settings from. Pass '-' to read from stdin.
      --space int              New volume space (in GB)
```

### Options inherited from parent commands

```
      --apikey string   API key
      --config string   config file (default is $HOME/.lvl.yaml)
  -o, --output string   Specifies output mode for commands. Accepted values are text or json. (default "text")
      --trace           Do detailed network request logging
```

### SEE ALSO

* [lvl system volume](lvl_system_volume.md)	 - Commands to manage volumes

