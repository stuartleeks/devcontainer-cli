package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func createCompleteCommand(rootCmd *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion SHELL",
		Short: "Generates bash completion scripts",
		Long: `To load completion run
	
	. <(devcontainer completion SHELL)

	Valid values for SHELL are : bash, fish, powershell, zsh
	
	For example, to configure your bash shell to load completions for each session add to your bashrc
	
	# ~/.bashrc or ~/.profile
	. <(devcontainer completion)
	`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				cmd.Usage()
				os.Exit(1)
			}
			shell := args[0]
			switch strings.ToLower(shell) {
			case "bash":
				rootCmd.GenBashCompletion(os.Stdout)
			case "fish":
				rootCmd.GenFishCompletion(os.Stdout, true)
			case "powershell":
				rootCmd.GenPowerShellCompletion(os.Stdout)
			case "zsh":
				rootCmd.GenPowerShellCompletion(os.Stdout)
			default:
				fmt.Printf("Unsupported SHELL value: '%s'\n", shell)
				cmd.Usage()
				os.Exit(1)
			}
		},
	}
	return cmd
}
