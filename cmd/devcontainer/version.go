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
		RunE: func(cmd *cobra.Command, args []string) error {
			if verbose {
				fmt.Printf("devcontainer version %s\nBuilt %s (commit %s)\n%s\n\n", version, date, commit, goversion)
				return nil
			}
			fmt.Println(version)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

	return cmd
}
