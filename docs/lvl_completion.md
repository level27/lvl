# lvl completion

A brief description of your command

### Synopsis

To load completions:

Bash:

  $ source <(lvl completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ lvl completion bash > /etc/bash_completion.d/lvl
  # macOS:
  $ lvl completion bash > /usr/local/etc/bash_completion.d/lvl

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ lvl completion zsh > "${fpath[1]}/_lvl"

  # You will need to start a new shell for this setup to take effect.

fish:

  $ lvl completion fish | source

  # To load completions for each session, execute once:
  $ lvl completion fish > ~/.config/fish/completions/lvl.fish

PowerShell:

  PS> lvl completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> lvl completion powershell > lvl.ps1
  # and source this file from your PowerShell profile.


```
lvl completion
```

### Options

```
  -h, --help   help for completion
```

### Options inherited from parent commands

```
      --apikey string   API key
      --config string   config file (default is $HOME/.lvl.yaml)
  -o, --output string   Specifies output mode for commands. Accepted values are text or json. (default "text")
      --trace           Do detailed network request logging
```

### SEE ALSO

* [lvl](lvl.md)	 - CLI tool to manage Level27 entities

