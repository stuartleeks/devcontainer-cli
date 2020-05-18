# devcontainer-cli

devcontainer-cli is the start of a CLI to improve the experience of working with [Visual Studio Code devcontainers](https://code.visualstudio.com/docs/remote/containers)

## Installation

TODO - add this once releases are automated!

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
