package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func createVersionCommand() *cobra.Command {

	var verbose bool

	cmd := &cobra.Command{
		Use:   "version",
		Short: "show version",
		Long:  "Show version",
		Run: func(cmd *cobra.Command, args []string) {
			if verbose {
				fmt.Printf("devcontainer version %s\nBuilt %s (commit %s)\n%s\n\n", version, date, commit, goversion)
				return
			}
			fmt.Println(version)
		},
	}
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

	return cmd
}
