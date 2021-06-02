package main

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/stuartleeks/devcontainer-cli/internal/pkg/devcontainers"
)

func createListCommand() *cobra.Command {
	var listIncludeContainerNames bool
	var listVerbose bool
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "List devcontainers",
		Long:  "Lists running devcontainers",
		RunE: func(cmd *cobra.Command, args []string) error {
			if listIncludeContainerNames && listVerbose {
				fmt.Println("Can't use both verbose and include-container-names")
				os.Exit(1)
			}
			devcontainers, err := devcontainers.ListDevcontainers()
			if err != nil {
				return err
			}
			if listVerbose {
				sort.Slice(devcontainers, func(i, j int) bool { return devcontainers[i].DevcontainerName < devcontainers[j].DevcontainerName })

				w := new(tabwriter.Writer)
				// minwidth, tabwidth, padding, padchar, flags
				w.Init(os.Stdout, 8, 8, 0, '\t', 0)
				defer w.Flush()

				fmt.Fprintf(w, "%s\t%s\n", "DEVCONTAINER NAME", "CONTAINER NAME")
				fmt.Fprintf(w, "%s\t%s\n", "-----------------", "--------------")

				for _, devcontainer := range devcontainers {
					fmt.Fprintf(w, "%s\t%s\n", devcontainer.DevcontainerName, devcontainer.ContainerName)
				}
				return nil
			}
			names := []string{}
			for _, devcontainer := range devcontainers {
				names = append(names, devcontainer.DevcontainerName)
				if listIncludeContainerNames {
					names = append(names, devcontainer.ContainerName)
				}
			}
			sort.Strings(names)
			for _, name := range names {
				fmt.Println(name)
			}
			return nil
		},
	}
	cmdList.Flags().BoolVar(&listIncludeContainerNames, "include-container-names", false, "Also include container names in the list")
	cmdList.Flags().BoolVarP(&listVerbose, "verbose", "v", false, "Verbose output")
	return cmdList
}

func countBooleans(values ...bool) int {
	count := 0
	for _, v := range values {
		if v {
			count++
		}
	}
	return count
}

func createExecCommand() *cobra.Command {
	var argDevcontainerName string
	var argDevcontainerPath string
	var argPromptForDevcontainer bool
	var argWorkDir string

	cmd := &cobra.Command{
		Use:   "exec [--name <name>| --path <path> | --prompt ] [--work-dir <work-dir>] [<command> [<args...>]] (command will default to /bin/bash if none provided)",
		Short: "Execute a command in a devcontainer",
		Long:  "Execute a command in a devcontainer, similar to `docker exec`",
		RunE: func(cmd *cobra.Command, args []string) error {

			// Default to executing /bin/bash
			if len(args) == 0 {
				args = []string{"/bin/bash"}
			}

			sourceCount := countBooleans(
				argDevcontainerName != "",
				argDevcontainerPath != "",
				argPromptForDevcontainer,
			)
			if sourceCount > 1 {
				fmt.Println("Can specify at most one of --name/--path/--prompt")
				return cmd.Usage()
			}

			// workDir default:
			// - devcontainer mount path if name or prompt specified (ExecInDevContainer defaults to this if workDir is "")
			// - path if path set
			// - current directory if path == "" and neither name or prompt set
			workDir := argWorkDir

			containerID := ""
			devcontainerList, err := devcontainers.ListDevcontainers()
			if err != nil {
				return err
			}
			if argDevcontainerName != "" {
				containerIDOrName := argDevcontainerName
				devcontainerList, err := devcontainers.ListDevcontainers()
				if err != nil {
					return err
				}

				// Get container ID
				for _, devcontainer := range devcontainerList {
					if devcontainer.ContainerName == containerIDOrName ||
						devcontainer.DevcontainerName == containerIDOrName ||
						devcontainer.ContainerID == containerIDOrName {
						containerID = devcontainer.ContainerID
						break
					}
				}

				if containerID == "" {
					return fmt.Errorf("Failed to find a matching (running) dev container for %q", containerIDOrName)
				}
			} else if argPromptForDevcontainer {
				// prompt user
				fmt.Println("Specify the devcontainer to use:")
				for index, devcontainer := range devcontainerList {
					fmt.Printf("%4d: %s (%s)\n", index, devcontainer.DevcontainerName, devcontainer.ContainerName)
				}
				selection := -1
				_, _ = fmt.Scanf("%d", &selection)
				if selection < 0 || selection >= len(devcontainerList) {
					return fmt.Errorf("Invalid option")
				}
				containerID = devcontainerList[selection].ContainerID
			} else {
				devcontainerPath := argDevcontainerPath
				// TODO - update to check for devcontainers in the path ancestry
				// Can't just check up the path for a .devcontainer folder as the container might
				// have been created via repository containers (https://github.com/microsoft/vscode-dev-containers/tree/main/repository-containers)
				containerID, err = devcontainers.GetContainerIDForPath(devcontainerPath)
				if err != nil {
					return err
				}
				if workDir == "" {
					if devcontainerPath == "" {
						workDir = "."
					} else {
						workDir = devcontainerPath
					}
				}
			}

			return devcontainers.ExecInDevContainer(containerID, workDir, args)
		},
		Args:                  cobra.ArbitraryArgs,
		DisableFlagsInUseLine: true,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
	}
	cmd.Flags().StringVarP(&argDevcontainerName, "name", "n", "", "name of dev container to exec into")
	cmd.Flags().StringVarP(&argDevcontainerPath, "path", "", "", "path containing the dev container to exec into")
	cmd.Flags().BoolVarP(&argPromptForDevcontainer, "prompt", "", false, "prompt for the dev container to exec into")
	cmd.Flags().StringVarP(&argWorkDir, "work-dir", "", "", "working directory to use in the dev container")

	_ = cmd.RegisterFlagCompletionFunc("name", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		devcontainers, err := devcontainers.ListDevcontainers()
		if err != nil {
			os.Exit(1)
		}
		names := []string{}
		for _, devcontainer := range devcontainers {
			names = append(names, devcontainer.DevcontainerName)
		}
		sort.Strings(names)
		return names, cobra.ShellCompDirectiveNoFileComp

	})
	return cmd
}
