package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

func createTemplateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "work with templates",
		Long:  "Use subcommands to work with devcontainer templates",
	}
	cmd.AddCommand(createTemplateListCommand())
	return cmd
}

func createTemplateListCommand() *cobra.Command {
	isDevcontainerFolder := func(parentPath string, fi os.FileInfo) bool {
		if !fi.IsDir() {
			return false
		}
		devcontainerJsonPath := fmt.Sprintf("%s/%s/.devcontainer/devcontainer.json", parentPath, fi.Name())
		devContainerJsonInfo, err := os.Stat(devcontainerJsonPath)
		return err == nil && !devContainerJsonInfo.IsDir()
	}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list templates",
		Long:  "List devcontainer templates",
		Run: func(cmd *cobra.Command, args []string) {
			const containerFolder string = "$HOME/source/vscode-dev-containers/containers" // TODO - make configurable!

			folder := os.ExpandEnv(containerFolder)
			c, err := ioutil.ReadDir(folder)
			if err != nil {
				fmt.Printf("Error reading devcontainer definitions: %s\n", err)
				os.Exit(1)
			}

			for _, entry := range c {
				if isDevcontainerFolder(folder, entry) {
					fmt.Println(entry.Name())
				}
			}

		},
	}
	return cmd
}
