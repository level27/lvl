# lvl app component create

Create a new appcomponent.

```
lvl app component create [flags]
```

### Examples

```
lvl app component create --name myComponentName --type docker
```

### Options

```
  -h, --help                 help for create
      --name string          
      --param stringArray    
  -f, --params-file string   JSON file to read params from. Pass '-' to read from stdin.
      --system string        
      --systemgroup string   
      --systemprovider int   
      --type string          
```

### Options inherited from parent commands

```
      --apikey string   API key
      --config string   config file (default is $HOME/.lvl.yaml)
  -o, --output string   Specifies output mode for commands. Accepted values are text or json. (default "text")
      --trace           Do detailed network request logging
```

### SEE ALSO

* [lvl app component](lvl_app_component.md)	 - Commands for managing appcomponents.

