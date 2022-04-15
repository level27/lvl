## lvl system integrity create

Create a new integrity report for a system.

```
lvl system integrity create [systemID] [flags]
```

### Options

```
      --doJobs      Create jobs (default: true) (default true)
      --forceJobs   Create jobs even if integrity check failed (default: false)
  -h, --help        help for create
```

### Options inherited from parent commands

```
      --apikey string   API key
      --config string   config file (default is $HOME/.lvl.yaml)
  -o, --output string   Specifies output mode for commands. Accepted values are text or json. (default "text")
      --trace           Do detailed network request logging
```

### SEE ALSO

* [lvl system integrity](lvl_system_integrity.md)	 - Manage integritychecks for a system
