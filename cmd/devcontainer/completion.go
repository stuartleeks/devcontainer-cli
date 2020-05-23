package main

import (
	"os"

	"github.com/spf13/cobra"
)

func createCompleteCommand(rootCmd *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:   "completion",
		Short: "Generates bash completion scripts",
		Long: `To load completion run
	
	. <(devcontainer completion)
	
	To configure your bash shell to load completions for each session add to your bashrc
	
	# ~/.bashrc or ~/.profile
	. <(devcontainer completion)
	`,
		Run: func(cmd *cobra.Command, args []string) {
			rootCmd.GenBashCompletion(os.Stdout)
		},
	}
}
