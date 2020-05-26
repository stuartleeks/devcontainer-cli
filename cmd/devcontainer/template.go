package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/stuartleeks/devcontainer-cli/internal/pkg/devcontainers"
	ioutil2 "github.com/stuartleeks/devcontainer-cli/internal/pkg/ioutil"
)

func createTemplateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "work with templates",
		Long:  "Use subcommands to work with devcontainer templates",
	}
	cmd.AddCommand(createTemplateListCommand())
	cmd.AddCommand(createTemplateAddCommand())
	cmd.AddCommand(createTemplateAddLinkCommand())
	return cmd
}

func createTemplateListCommand() *cobra.Command {
	var listVerbose bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list templates",
		Long:  "List devcontainer templates",
		RunE: func(cmd *cobra.Command, args []string) error {

			templates, err := devcontainers.GetTemplates()
			if err != nil {
				return err
			}

			if listVerbose {
				w := new(tabwriter.Writer)
				// minwidth, tabwidth, padding, padchar, flags
				w.Init(os.Stdout, 8, 8, 0, '\t', 0)
				defer w.Flush()

				fmt.Fprintf(w, "%s\t%s\n", "TEMPLATE NAME", "PATH")
				fmt.Fprintf(w, "%s\t%s\n", "-------------", "----")

				for _, template := range templates {
					fmt.Fprintf(w, "%s\t%s\n", template.Name, template.Path)
				}
				return nil
			}

			for _, template := range templates {
				fmt.Println(template.Name)
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&listVerbose, "verbose", "v", false, "Verbose output")
	return cmd
}

func createTemplateAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add TEMPLATE_NAME",
		Short: "add devcontainer from template",
		Long:  "Add a devcontainer definition to the current folder using the specified template",
		RunE: func(cmd *cobra.Command, args []string) error {

			if len(args) != 1 {
				return cmd.Usage()
			}
			name := args[0]

			template, err := devcontainers.GetTemplateByName(name)
			if err != nil {
				return err
			}
			if template == nil {
				fmt.Printf("Template '%s' not found\n", name)
			}

			info, err := os.Stat("./.devcontainer")
			if info != nil && err == nil {
				return fmt.Errorf("Current folder already contains a .devcontainer folder - exiting")
			}

			currentDirectory, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("Error reading current directory: %s\n", err)
			}
			if err = ioutil2.CopyFolder(template.Path, currentDirectory+"/.devcontainer"); err != nil {
				return fmt.Errorf("Error copying folder: %s\n", err)
			}
			return err
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			// only completing the first arg  (template name)
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			templates, err := devcontainers.GetTemplates()
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
	return cmd
}

func createTemplateAddLinkCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-link TEMPLATE_NAME",
		Short: "add-link devcontainer from template",
		Long:  "Symlink a devcontainer definition to the current folder using the specified template",
		RunE: func(cmd *cobra.Command, args []string) error {

			if len(args) != 1 {
				return cmd.Usage()
			}
			name := args[0]

			template, err := devcontainers.GetTemplateByName(name)
			if err != nil {
				return err
			}
			if template == nil {
				return fmt.Errorf("Template '%s' not found\n", name)
			}

			info, err := os.Stat("./.devcontainer")
			if info != nil && err == nil {
				return fmt.Errorf("Current folder already contains a .devcontainer folder - exiting")
			}

			currentDirectory, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("Error reading current directory: %s\n", err)
			}
			if err = ioutil2.LinkFolder(template.Path, currentDirectory+"/.devcontainer"); err != nil {
				return fmt.Errorf("Error linking folder: %s\n", err)
			}

			content := []byte("*\n")
			if err := ioutil.WriteFile(currentDirectory+"/.devcontainer/.gitignore", content, 0644); err != nil { // -rw-r--r--
				return fmt.Errorf("Error writing .gitignore: %s\n", err)
			}
			return err
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			// only completing the first arg  (template name)
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			templates, err := devcontainers.GetTemplates()
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
	return cmd
}
