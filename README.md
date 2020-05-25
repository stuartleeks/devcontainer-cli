# devcontainer-cli

devcontainer-cli is the start of a CLI to improve the experience of working with [Visual Studio Code devcontainers](https://code.visualstudio.com/docs/remote/containers)

**Status: this is a pet project that I've been experimenting with. It is not supported and you should expect bugs :-)**

## Installation

Head to the [latest release page](https://github.com/stuartleeks/devcontainer-cli/releases/latest) and download the archive for your platform.

Extract `devcontainer` from the archive and place in a folder in your `PATH`.

Or if you just don't care and are happy to run random scripts from the internet:

```bash
OS=linux # also darwin
ARCH=amd64 # also 386
wget https://raw.githubusercontent.com/stuartleeks/devcontainer-cli/master/scripts/install.sh
chmod +x install.sh
sudo -E ./install.sh
```

## Enabling bash completion

To enable bash completion, add the following to you `~/.bashrc` file:

```bash
. <(devcontainer completion)
```

## Usage

### Listing devcontainers

To see which running devcontainers the CLI detects you can run the `list` command.

### Running commands inside a devcontainer

`devcontainer` allows you to run commands in devcontainers. This is similar to `docker exec` but works with devcontainer names (rather than requiring container names/IDs). 

For example:

```bash
# Run an interactive bash shell in the vscode-remote-test-dockerfile devcontainer
devcontainer exec vscode-remote-test-dockerfile bash

# Run a command with args in the vscode-remote-test-dockercompose_devcontainer/mongo devcontainer
devcontainer exec vscode-remote-test-dockercompose_devcontainer/mongo ls -a /workspaces/vscode-remote-test-dockerfile
```
