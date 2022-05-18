# lvl app ssl

Commands for managing SSL certificates on apps

### Examples

```
lvl app ssl get forum
lvl app ssl describe forum forum.example.com
```

### Options

```
  -f, --filter string   How to filter API results?
  -h, --help            help for ssl
  -n, --number int      How many things should we retrieve from the API?
```

### Options inherited from parent commands

```
      --apikey string   API key
      --config string   config file (default is $HOME/.lvl.yaml)
  -o, --output string   Specifies output mode for commands. Accepted values are text or json. (default "text")
      --trace           Do detailed network request logging
```

### SEE ALSO

* [lvl app](lvl_app.md)	 - Commands to manage apps
* [lvl app ssl action](lvl_app_ssl_action.md)	 - 
* [lvl app ssl create](lvl_app_ssl_create.md)	 - Create a new SSL certificate on an app
* [lvl app ssl delete](lvl_app_ssl_delete.md)	 - Delete an SSL certificate from an app
* [lvl app ssl describe](lvl_app_ssl_describe.md)	 - Get detailed information of an SSL certificate
* [lvl app ssl fix](lvl_app_ssl_fix.md)	 - Fix an invalid SSL certificate
* [lvl app ssl get](lvl_app_ssl_get.md)	 - Get a list of SSL certificates for an app
* [lvl app ssl key](lvl_app_ssl_key.md)	 - Return a private key for type 'own' sslCertificate.
* [lvl app ssl update](lvl_app_ssl_update.md)	 - 

