If you find frequently find yourself working at the terminal and are working with [Visual Studio Code dev containers](https://code.visualstudio.com/docs/remote/containers) then the `devcontainer` CLI might be of interest for you!

Examples:

```bash

# The following command opens the current folder in 
# VS Code as a dev container # i.e. it skips the
# normal step of opening in VS Code and then 
# clicking # on the "Re-open in container" prompt.
$ devcontainer open-in-code

# If you don't have a dev container definition for 
# your folder then you can use 
#`devcontainer template add <name>` to add a 
# dev container definition.
$ devcontainer template add python-3

# You can use `devcontainer exec` to create a 
# shell (or run a process) # inside a dev container.
$ devcontainer exec
```

See the following topics for more information:

* [Installation](installation)
* Commands
  * [open-in-code](open-in-code) - open dev containers in VS Code from the terminal
  * [template](template) - add dev container definitions to a folder
  * [exec](exec) - launch a terminal or other command in a dev container

