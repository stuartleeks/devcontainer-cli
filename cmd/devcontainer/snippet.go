package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/stuartleeks/devcontainer-cli/internal/pkg/devcontainers"
)

func createSnippetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "snippet",
		Short: "work with snippets (experimental)",
		Long:  "Use subcommands to work with devcontainer snippets (experimental)",
	}
	cmd.AddCommand(createSnippetListCommand())
	// cmd.AddCommand(createSnippetAddCommand())
	return cmd
}

func createSnippetListCommand() *cobra.Command {
	var listVerbose bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list snippets",
		Long:  "List devcontainer snippets",
		RunE: func(cmd *cobra.Command, args []string) error {

			snippets, err := devcontainers.GetSnippets()
			if err != nil {
				return err
			}

			if listVerbose {
				w := new(tabwriter.Writer)
				// minwidth, tabwidth, padding, padchar, flags
				w.Init(os.Stdout, 8, 8, 0, '\t', 0)
				defer w.Flush()

				fmt.Fprintf(w, "%s\t%s\n", "SNIPPET NAME", "PATH")
				fmt.Fprintf(w, "%s\t%s\n", "-------------", "----")

				for _, snippet := range snippets {
					fmt.Fprintf(w, "%s\t%s\n", snippet.Name, snippet.Path)
				}
				return nil
			}

			for _, snippet := range snippets {
				fmt.Println(snippet.Name)
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&listVerbose, "verbose", "v", false, "Verbose output")
	return cmd
}
