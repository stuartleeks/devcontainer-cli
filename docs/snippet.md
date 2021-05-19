# devcontainer snippet ...

***WARNING: This feature is experimental!***

{:toc}

## Setting up snippets

To use the `devcontainer snippet` commands you need to enable experimental feaaures and configure some snippet folders.

A snippet collection can be as simple as a set of `.sh` scripts. Good starting points for snippets are the `snippets` folder of [stuartleeks/devcontainers](https://github.com/stuartleeks/devcontainers) and [benc-uk/tools-install](https://github.com/benc-uk/tools-install/).

Choose a directory to put the snippets in, and clone the repos:

```bash
git clone https://github.com/stuartleeks/devcontainers
git clone https://github.com/benc-uk/tools-install/
```

Next, we need to tell the `devcontainer` CLI to use these folders. If you haven't previously created a config file, run `devcontainer config write` to save a config file and then open `~/.devcontainer-cli/devcontainer-cli.json` in your favourite editor.

The starting configuration will look something like:

```json
{
    "templatepaths": []
}
```

Add a `snippetpaths` setting with an array value containing the path to your snippet folders. For example, if you cloned the repos into your `~/source` folder add the following to your config file:

```json
{
    "experimental" : true,
    "snippetpaths": ["$HOME/source/sl-devcontainers/snippets", "$HOME/source/tools-install"]
}
```

NOTE: You also need to add the `experimental` setting with the value `true` as snippets are currently an experimental features

## Listing snippets

Running `devcontainer snippet list` will show the snippets that `devcontainer` discovered

## Adding a snippet

To add a snippet to the dev container definition to your project, change directory to the project folder (i.e. the one containing the `.devcontainer` folder) and then run:

```bash
# Add the azbrowse
devcontainer snippet add azbrowse
```

This will copy in the snippet files for you to modify as you wish.

## Creating your own snippets

`devcontainer` can be configured to scan multiple folders to find snippets. For each folder configured in the `snippetpaths` setting it searches for snippets. There are currently two types of snippet supported: single file snippets and folder-based snippets.

### Single file snippets

Single file snippets are `.sh` files that are copied to the `.devcontainer/scripts` folder and added to the `Dockerfile`. This is the simplest place to start when creating a snippet.

### Folder-based snippets

Folder-based snippets are folders containing a `snippet.json` file. The `snippet.json` describes the actions to take when applying the snippet, for example:

```json
{
    "actions": [
        {
            "type" : "copyAndRun",
            "source": "my-script.sh"
        }
    ]
}
```

The `actions` property can contain multiple actions and they are applied in order.

The following action types are supported:

- `copyAndRun`
- `mergeJSON` 
- `dockerfileSnippet`

#### copyAndRun action

The `copyAndRun` action provides the same capability as the single file snippet, i.e. the source file is copied and added to the `Dockerfile`.

The following properties are supported for a `copyAndRun` action:

| Property | Description                                                        |
|----------|--------------------------------------------------------------------|
| source   | The path to the script file to copy (relative to the snippet.json) |

For example:

```json
{
    "actions": [
        {
            "type": "copyAndRun",
            "source": "golang.sh"
        }
    ]
}
```

#### mergeJSON action

The `mergeJSON` action provides the ability to merge changes into a JSON file (e.g. `devcontainer.json`).

The following properties are supported for a `mergeJSON` action:

| Property | Description                                                                   |
|----------|-------------------------------------------------------------------------------|
| source   | The path to the JSON file containing the properties to merge in to the target |
| target   | The path to the JSON file to merge the changes into                           |

For example:

```json
{
    "actions": [
        {
            "type": "mergeJSON",
            "source": "devcontainer.json",
            "target": ".devcontainer/devcontainer.json"
        }
    ]
}
```

#### dockerfileSnippet action

The `dockerfileSnippet` action provides a way to add custom steps to the `Dockerfile` for a dev container.

The following properties are supported for a `dockerfileSnippet` action:

| Property | Description                          |
|----------|--------------------------------------|
| content  | The content to add to the Dockerfile |

For example:

```json
{
    "actions": [
        {
            "type": "dockerfileSnippet",
            "content": "# Add go to PATH\nENV PATH /usr/local/go/bin:$PATH"
        }
    ]
}
```
