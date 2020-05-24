package main

import (
	"github.com/spf13/cobra"
)

func main() {

	rootCmd := &cobra.Command{Use: "devcontainer"}

	rootCmd.AddCommand(createListCommand())
	rootCmd.AddCommand(createExecCommand())
	rootCmd.AddCommand(createTemplateCommand())
	rootCmd.AddCommand(createCompleteCommand(rootCmd))
	rootCmd.AddCommand(createConfigCommand())

	rootCmd.Execute()
}
