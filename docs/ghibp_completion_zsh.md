## ghibp completion zsh

Generate the autocompletion script for zsh

### Synopsis

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

	source <(ghibp completion zsh)

To load completions for every new session, execute once:

#### Linux:

	ghibp completion zsh > "${fpath[1]}/_ghibp"

#### macOS:

	ghibp completion zsh > $(brew --prefix)/share/zsh/site-functions/_ghibp

You will need to start a new shell for this setup to take effect.


```
ghibp completion zsh [flags]
```

### Options

```
  -h, --help              help for zsh
      --no-descriptions   disable completion descriptions
```

### SEE ALSO

* [ghibp completion](ghibp_completion.md)	 - Generate the autocompletion script for the specified shell

