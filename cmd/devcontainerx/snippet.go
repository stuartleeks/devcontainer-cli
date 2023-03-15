package main

import (
	"fmt"
	"os"
	"sort"
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
	cmd.AddCommand(createSnippetAddCommand())
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

func createSnippetAddCommand() *cobra.Command {
	var devcontainerName string
	cmd := &cobra.Command{
		Use:   "add SNIPPET_NAME",
		Short: "add snippet to devcontainer",
		Long:  "Add a snippet to the devcontainer definition for the current folder",
		RunE: func(cmd *cobra.Command, args []string) error {

			if len(args) != 1 {
				return cmd.Usage()
			}
			name := args[0]

			currentDirectory, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("Error reading current directory: %s\n", err)
			}

			err = devcontainers.AddSnippetToDevcontainer(currentDirectory, name)
			if err != nil {
				return fmt.Errorf("Error setting devcontainer name: %s", err)
			}

			return nil
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			// only completing the first arg  (template name)
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			templates, err := devcontainers.GetSnippets()
			if err != nil {
				os.Exit(1)
			}
			names := []string{}
			for _, template := range templates {
				names = append(names, template.Name)
			}
			sort.Strings(names)
			return names, cobra.ShellCompDirectiveNoFileComp
		},
	}
	cmd.Flags().StringVar(&devcontainerName, "devcontainer-name", "", "Value to set the devcontainer.json name property to (default is folder name)")
	return cmd
}
