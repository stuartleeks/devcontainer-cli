package main

import (
	"github.com/spf13/cobra"
)

// Overridden via ldflags
var (
	version   = "99.0.1-devbuild"
	commit    = "unknown"
	date      = "unknown"
	goversion = "unknown"
)

func main() {

	rootCmd := &cobra.Command{Use: "devcontainer"}

	rootCmd.AddCommand(createListCommand())
	rootCmd.AddCommand(createExecCommand())
	rootCmd.AddCommand(createTemplateCommand())
	rootCmd.AddCommand(createCompleteCommand(rootCmd))
	rootCmd.AddCommand(createConfigCommand())
	rootCmd.AddCommand(createVersionCommand())

	rootCmd.Execute()
}
