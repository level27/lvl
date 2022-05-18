# lvl system create

Create a new system

```
lvl system create [flags]
```

### Examples

```
lvl system create -n mySystemName --zone hasselt --organisation level27 --image 'Ubuntu 20.04 LTS' --config 'Level27 Small' --management professional_level27
```

### Options

```
      --Fqdn string            Valid hostname for the system
      --autoTeams string       A csv list of team ID's
      --cpu int                Cpu (Required for Level27 systems)
      --disk int               Disk (non-editable)
      --externalInfo string    ExternalInfo (required when billableItemInfo entities for an organisation exist in db)
  -h, --help                   help for create
      --image string           The ID of a systemimage. (must match selected configuration and zone. non-editable)
      --management string      Managament type (one of basic, professional, enterprise, professional_level27). (default "basic")
      --memory int             Memory (Required for Level27 systems)
  -n, --name string            The name you want to give the system
      --networks stringArray   Array of network IP's. (default: null)
      --organisation string    The unique ID of an organisation
      --parent int             The unique ID of a system (parent system)
      --publicNetworking       For digitalOcean servers always true. (non-editable) (default true)
      --remarks string         Remarks (Admin only)
      --type string            System type
      --version int            The unique ID of an OperatingsystemVersion (non-editable)
      --zone string            The unique ID of a zone
```

### Options inherited from parent commands

```
      --apikey string   API key
      --config string   config file (default is $HOME/.lvl.yaml)
  -o, --output string   Specifies output mode for commands. Accepted values are text or json. (default "text")
      --trace           Do detailed network request logging
```

### SEE ALSO

* [lvl system](lvl_system.md)	 - Commands for managing systems

