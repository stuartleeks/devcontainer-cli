# devcontainer-cli

devcontainer-cli is the start of a CLI to improve the experience of working with [Visual Studio Code devcontainers](https://code.visualstudio.com/docs/remote/containers)

**Status: this is a pet project that I've been experimenting with. It is not supported and you should expect bugs :-)**

## Installation

Head to the [latest release page](https://github.com/stuartleeks/devcontainer-cli/releases/latest) and download the archive for your platform.

Extract `devcontainer` from the archive and place in a folder in your `PATH`.

Or if you just don't care and are happy to run random scripts from the internet:

```bash
export OS=linux # also darwin
export ARCH=amd64 # also 386
wget https://raw.githubusercontent.com/stuartleeks/devcontainer-cli/main/scripts/install.sh
chmod +x install.sh
sudo -E ./install.sh
```

## Enabling bash completion

To enable bash completion, add the following to you `~/.bashrc` file:

```bash
source <(devcontainer completion bash)
```

Or to alias `devcontainer` (to `dc` in this example):

```bash
alias dc=devcontainer
complete -F __start_devcontainer dc
```

## Usage

### Working with devcontainers

#### Listing devcontainers

To see which running devcontainers the CLI detects you can run the `list` command.

#### Running commands inside a devcontainer

`devcontainer` allows you to run commands in devcontainers. This is similar to `docker exec` but works with devcontainer names (rather than requiring container names/IDs).

For example:

```bash
# Run an interactive bash shell in the vscode-remote-test-dockerfile devcontainer
devcontainer exec --name vscode-remote-test-dockerfile bash

# Run a command with args in the vscode-remote-test-dockercompose_devcontainer/mongo devcontainer
devcontainer exec --name vscode-remote-test-dockercompose_devcontainer/mongo ls -a /workspaces/vscode-remote-test-dockerfile

# Run `bash` in the dev container for the project at `~/ source/my-proj`
devcontainer exec --path ~/source/my-proj bash

# If none of --name/--path/--prompt are specified then `--path .` is assumed (i.e. use the dev container for the current directory)
devcontainer exec bash

# If command/args not set, `bash` is assumed
devcontainer exec --name vscode-remote-test-dockerfile

# Combining these to launch bash in the dev container for the project in the current directory:
devcontainer exec
```

You can use `--prompt` instead of `--name` or `--path` and the CLI will prompt you to pick a devcontainer to run the `exec` command against, e.g.:

```bash
$ ./devcontainer exec ? bash
Specify the devcontainer to use:
   0: devcontainer-cli (festive_saha)
   1: vscode-remote-test-dockerfile (fervent_gopher)
0
```

You can use this with Windows Terminal profiles:

```json
{
    "guid": "{4b304185-99d2-493c-940c-ae74e0f14bba}",
    "hidden": false,
    "name": "devcontainer exec",
    "commandline": "wsl bash -c \"path/to/devcontainer exec --prompt bash\"",
},
```

By default, `devcontainer exec` will set the working directory to be the mount path for the dev container. This can be overridden using `--work-dir`.

### Working with devcontainer templates

To work with devcontainer templates `devcontainer` needs to know where you have the templates stored.

As a quickstart, clone the VS Code devcontainers repo: `git clone https://github.com/microsoft/vscode-dev-containers`

Next, run `devcontainer config write` to save a config file and then open `~/.devcontainer-cli/devcontainer-cli.json` in your favourite editor.

The starting configuration will look something like:

```json
{
  "templatepaths": []
}
```

Update to include the path to the `containers` folder in the `vscode-dev-containers` repo you just cloned:

```json
{
  "templatepaths": ["$HOME/source/vscode-dev-containers/containers"]
}
```

See [Template Paths](#template-paths) for more details of the structure of template folders.

#### Listing templates

Running `devcontainer template list` will show the templates that `devcontainer` discovered

#### Adding a devcontainer

To add the files for a devcontainer definition to your project, change directory to the folder you want to add the devcontainer to and then run:

```bash
# Add the go template
devcontainer template add go
```

This will copy in the template files for you to modify as you wish.

#### Adding a link to a devcontainer

If you are working with a codebase that you don't want to commit the devcontainer definition to (e.g. an OSS project that doesn't want a devcontainer definition), you can use the `template add-link` command. Instead of copying template files it creates symlinks to the template files and adds a `.gitignore` file to avoid accidental git commits.

As with `template add`, run this from the folder you want to add the devcontainer to:

```bash
# Symlink to the go template
devcontainer template add-link go
```

## Template paths

`devcontainer` can be [configured to scan multiple folders](#working-with-devcontainer-templates) to find templates. It is designed to work with folders structured in the same was as the [containers from in github.com/microsoft/vscode-dev-containers](https://github.com/microsoft/vscode-dev-containers/tree/main/containers).

Assuming you cloned [github.com/microsoft/vscode-dev-containers/](https://github.com/microsoft/vscode-dev-containers/) into your `~/source/` folder and set up a custom devcontainer folder in `~/source/devcontainers` then you can configure your template paths as shown below. The sub-folder names are used as the template name and when duplicates are found the first matching folder is taken, so in the example below the `~/source/devcontainers` templates take precedence.

```json
{
  "templatepaths": [
    "$HOME/source/devcontainers",
    "$HOME/source/vscode-dev-containers/containers"
  ]
}
```

The structure for these template paths is shown in the following tree structure:

```misc
template-collection-folder
 |-template1
 |  |-.devcontainer
 |  |  |-devcontainer.json
 |  |  |-Dockerfile
 |  |  |-<other content for the template>
 |-misc-folder
 |-<misc content that is ignored as there is no .devcontainer folder>
 |-<README or other files that are ignore>
```
