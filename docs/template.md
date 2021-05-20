# devcontainer template ...

{:toc}

## Setting up templates

To use the `devcontainer template` commands you need to configure some templates.

A good starting point is the the VS Code devcontainers repo. Choose a directory, and clone the repo using  `git clone https://github.com/microsoft/vscode-dev-containers`

Next, we need to tell the `devcontainer` CLI to use this folder. If you haven't previously created a config file, run `devcontainer config write` to save a config file and then open `~/.devcontainer-cli/devcontainer-cli.json` in your favourite editor.

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

## Listing templates

Running `devcontainer template list` will show the templates that `devcontainer` discovered

## Adding a devcontainer definition

To add the files for a devcontainer definition to your project, change directory to the folder you want to add the devcontainer to and then run:

```bash
# Add the go template
devcontainer template add go
```

This will copy in the template files for you to modify as you wish.

## Adding a link to a devcontainer

If you are working with a codebase that you don't want to commit the devcontainer definition to (e.g. an OSS project that doesn't want a devcontainer definition), you can use the `template add-link` command. Instead of copying template files it creates symlinks to the template files and adds a `.gitignore` file to avoid accidental git commits.

As with `template add`, run this from the folder you want to add the devcontainer to:

```bash
# Symlink to the go template
devcontainer template add-link go
```

See the [repository containers](#repository-containers) section for an alternative to template links.

## Creating your own templates

`devcontainer` can be configured to scan multiple folders to find templates. It is designed to work with folders structured in the same was as the [containers from in github.com/microsoft/vscode-dev-containers](https://github.com/microsoft/vscode-dev-containers/tree/master/containers), e.g.:


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

Assuming you cloned [github.com/microsoft/vscode-dev-containers/](https://github.com/microsoft/vscode-dev-containers/) into your `~/source/` folder and set up a custom devcontainer folder in `~/source/devcontainers` then you can configure your template paths as shown below. The sub-folder names are used as the template name and when duplicates are found the first matching folder is taken, so in the example below the `~/source/devcontainers` templates take precedence.

```json
{
  "templatepaths": [
    "$HOME/source/devcontainers",
    "$HOME/source/vscode-dev-containers/containers"
  ]
}
```

## Placeholder Values

After content has been copied to the project folder from a template, the following placeholder values are substituted:

| Placeholder                  | Value                                                                                                                |
|------------------------------|----------------------------------------------------------------------------------------------------------------------|
| `__DEVCONTAINER_NAME__`      | The name of the dev container (from the `name` property in `devcontainer.json`)                                      |
| `__DEVCONTAINER_USER_NAME__` | The name of the user for dev container (from the `remoteuser` property in `devcontainer.json`, or `root` if not set) |
| `__DEVCONTAINER_HOME__`      | The home folder for the dev container (e.g. `/home/vscode` or `/root`)                                               |

## Repository containers

VS Code dev containers have another feature called "Repository containers". These are a set of dev container definitions that VS Code will automatically apply to a project based on its git repo.

The default definitions are in the [microsoft/vscode-dev-containers](https://github.com/microsoft/vscode-dev-containers/tree/master/repository-containers) repo. If you look at the repo, you will see a `github.com` folder followed by paths for `<org>/<repo>`, e.g. `django/django`. The `https://github.com/django/django` repo doesn't contain a dev container definition, but VS Code will use the repository container definition from the `microsoft/vscode-dev-containers` repo.

You can also configure VS Code to look for additional local paths for repository containers by providing a value for the VS Code `remote.containers.repository-container-paths` setting (see [this issue](https://github.com/microsoft/vscode-remote-release/issues/3218) for more details).
