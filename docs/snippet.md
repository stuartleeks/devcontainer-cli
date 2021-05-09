# devcontainer snippet ...

***WARNING: This feature is experimental***

{:toc}


## Setting up snippets


To use the `devcontainer snippet` commands you need to enable experimental feaaures and configure some snippet folders.

A snippet folder can be as simple as a set of `.sh` scripts - examples can be found in the `snippets` folder of [stuartleeks/devcontainers](https://github.com/stuartleeks/devcontainers). 

Choose a directory, and clone the repo using  `git clone https://github.com/stuartleeks/devcontainers`

Next, we need to tell the `devcontainer` CLI to use this folder. If you haven't previously created a config file, run `devcontainer config write` to save a config file and then open `~/.devcontainer-cli/devcontainer-cli.json` in your favourite editor.

The starting configuration will look something like:

```json
{
  "templatepaths": []
}
```

Add an `experimental` setting with the value `true` to turn on experimental features, and add a `snippetpaths` setting with an array value containing the path to your snippet folder, e.g.:

```json
{
  "experimental" : true,
  "snippetpaths": ["$HOME/source/sl-devcontainers/snippets"]
}
```

## Listing snippets

Running `devcontainer snippet list` will show the snippets that `devcontainer` discovered

## Adding a snippet

To add a snippet to the dev container definition to your project, change directory to the project folder (i.e. the one with the `.devcontainer` folder) and then run:

```bash
# Add the azbrowse
devcontainer snippet add azbrowse
```

This will copy in the snippet files for you to modify as you wish.

## Creating your own snippets **TODO**

`devcontainer` can be configured to scan multiple folders to find snippets. For each folder configured in the `snippetpaths` setting it searches for `*.sh` files.
