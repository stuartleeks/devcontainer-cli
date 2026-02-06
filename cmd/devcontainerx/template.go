package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"
	"github.com/stuartleeks/devcontainer-cli/internal/pkg/devcontainers"
	ioutil2 "github.com/stuartleeks/devcontainer-cli/internal/pkg/ioutil"
	"github.com/stuartleeks/devcontainer-cli/internal/pkg/output"
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
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list templates",
		Long:  "List devcontainer templates",
		RunE: func(cmd *cobra.Command, args []string) error {
			outputFormat, query, err := output.GetOutputAndQueryValues(cmd, `[].{name: name, path: path}`)
			if err != nil {
				return err
			}

			templates, err := devcontainers.GetTemplates()
			if err != nil {
				return err
			}

			err = output.OutputResult(os.Stdout, templates, outputFormat, query, []string{"name", "path"})
			if err != nil {
				return fmt.Errorf("error outputting result: %s", err)
			}
			return nil
		},
	}
	output.AddOutputAndQueryFlags(cmd)
	return cmd
}

func createTemplateAddCommand() *cobra.Command {
	var devcontainerName string
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
				return fmt.Errorf("current folder already contains a .devcontainer folder - exiting")
			}

			currentDirectory, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("error reading current directory: %s", err)
			}

			err = devcontainers.CopyTemplateToFolder(template.Path, currentDirectory, devcontainerName)
			if err != nil {
				return err
			}

			return nil
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
	cmd.Flags().StringVar(&devcontainerName, "devcontainer-name", "", "Value to set the devcontainer.json name property to (default is folder name)")
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
				return fmt.Errorf("template '%s' not found", name)
			}

			info, err := os.Stat("./.devcontainer")
			if info != nil && err == nil {
				return fmt.Errorf("current folder already contains a .devcontainer folder - exiting")
			}

			currentDirectory, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("error reading current directory: %s", err)
			}
			if err = ioutil2.LinkFolder(template.Path, currentDirectory+"/.devcontainer"); err != nil {
				return fmt.Errorf("error linking folder: %s", err)
			}

			content := []byte("*\n")
			if err := os.WriteFile(currentDirectory+"/.devcontainer/.gitignore", content, 0644); err != nil { // -rw-r--r--
				return fmt.Errorf("error writing .gitignore: %s", err)
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
