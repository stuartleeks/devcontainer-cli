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
	source <(devcontainer completion)

	# if you want to alias the CLI:
	alias dc=devcontainer
    source <(devcontainer completion bash | sed s/devcontainer/dc/g)

	`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				_ = cmd.Usage()
				os.Exit(1)
			}
			shell := args[0]
			var err error
			switch strings.ToLower(shell) {
			case "bash":
				err = rootCmd.GenBashCompletion(os.Stdout)
			case "fish":
				err = rootCmd.GenFishCompletion(os.Stdout, true)
			case "powershell":
				err = rootCmd.GenPowerShellCompletion(os.Stdout)
			case "zsh":
				err = rootCmd.GenZshCompletion(os.Stdout)
			default:
				fmt.Printf("Unsupported SHELL value: '%s'\n", shell)
				return cmd.Usage()
			}

			return err
		},
	}
	return cmd
}
