# Installation


## Download latest release

Head to the [latest release page](https://github.com/stuartleeks/devcontainer-cli/releases/latest) and download the archive for your platform.

Extract `devcontainer` from the archive and place in a folder in your `PATH`.

## Homebrew

You can also install using `homebrew` with `brew install stuartleeks/tap/devcontainer`

## Just give me a script

Or if you just don't care and are happy to run random scripts from the internet:

```bash
export OS=linux # also darwin
export ARCH=amd64 # also 386
wget https://raw.githubusercontent.com/stuartleeks/devcontainer-cli/main/scripts/install.sh
chmod +x install.sh
sudo -E ./install.sh
```

## Enabling bash completion

The `devcontainer completion <shell>` command generates a completion script for the specified shell. 

To enable bash completion, add the following to you `~/.bashrc` file:

```bash
source <(devcontainer completion bash)
```

Or to alias `devcontainer` (to `dc` in this example):

```bash
alias dc=devcontainer
complete -F __start_devcontainer dc
```

The `devcontainer completion <shell>` command accepts `bash`, `zsh`, and `powershell` for the `<shell>` parameter.
