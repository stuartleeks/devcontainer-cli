# devcontainer exec

The `devcontainer exec` command can be used to run commands inside a running dev container. When no options are passed, it will run `bash` in the dev container for the dev container for the current working directory.

Some examples of using `devcontainer exec` are shown below:

```bash
# Run an interactive bash shell in the 
# vscode-remote-test-dockerfile devcontainer
devcontainer exec --name vscode-remote-test-dockerfile bash

# Run a command with args in the 
# vscode-remote-test-dockercompose_devcontainer/mongo 
# devcontainer
devcontainer exec --name vscode-remote-test-dockercompose_devcontainer/mongo ls -a /workspaces/vscode-remote-test-dockerfile

# Run `bash` in the dev container for 
# the project at `~/ source/my-proj`
devcontainer exec --path ~/source/my-proj bash

# If none of --name/--path/--prompt 
# are specified then `--path .` is assumed 
# (i.e. use the dev container for the current directory)
devcontainer exec bash

# If command/args not set, `bash` is assumed
devcontainer exec --name vscode-remote-test-dockerfile

# Combining these to launch bash in the 
# dev container for the project in the current directory:
devcontainer exec
```

## Features of devcontainer exec

Under the covers, `devcontainer exec` launches `docker exec`, but it has a few features on top of this to try to increase productivity.

First, it sets the working directory to be the mount path for the dev container rather than just dropping you in at the root of the container flie system. This can be overridden using `--work-dir`.

Second, it checks whether you have [configured a user in the dev container](https://code.visualstudio.com/docs/remote/containers-advanced#_adding-a-nonroot-user-to-your-dev-container) and uses this user for the `docker exec`.

Lastly, it checks whether you have set up an SSH agent on your host. If you have and VS Code detects it then VS Code will [forward key requests from the container](https://code.visualstudio.com/docs/remote/containers#_using-ssh-keys). In this scenario, `devcontainer exec` configures the exec session to also forward key requests. This enables operations against git remotes secured with SSH keys to succeed.


## Prompting for the dev container


You can use `--prompt` with `devcontainer exec` instead of `--name` or `--path` and the CLI will prompt you to pick a devcontainer to run the `exec` command against, e.g.:

```bash
$ ./devcontainer exec ? bash
Specify the devcontainer to use:
   0: devcontainer-cli (festive_saha)
   1: vscode-remote-test-dockerfile (fervent_gopher)
0
```

This works well as a terminal profile. For example, you can use this with Windows Terminal profiles:

```json
{
    "guid": "{4b304185-99d2-493c-940c-ae74e0f14bba}",
    "hidden": false,
    "name": "devcontainer exec",
    "commandline": "wsl bash -c \"path/to/devcontainer exec --prompt bash\"",
},
```
