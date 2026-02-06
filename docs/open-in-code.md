# devcontainer open-in-code

The `devcontainer open-in-code` command opens the current folder in VS Code as a dev container, i.e. it skips the normal step of opening in VS Code and then clicking # on the "Re-open in container" prompt to reload the window as a dev container. With `devcontainer open-in-code` you get straight to the dev container!

You can also use `devcontainer open-in-code <path>` to open a the devcontainer from a different folder. 

The command also supports the `--sub-folder` option. This allows you to specify a sub-folder within the folder that contains the devcontainer configuration (i.e. with the `.devcontainer` folder or `.devcontainer.json` file). 

The `open-in-code` command will walk up the tree from the current/specified folder until it finds a folder with a devcontainer configuration, and then open that folder in VS Code as a dev container.
If `--sub-folder` is not specified and the devcontainer configuration is found in a parent folder, relative path to the current folder is used as the sub-folder.

If you want to use the VS Code Insiders release, you can use `devcontainer open-in-code-insiders`.