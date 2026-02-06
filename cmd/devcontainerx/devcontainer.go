package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"
	"github.com/stuartleeks/devcontainer-cli/internal/pkg/devcontainers"
	"github.com/stuartleeks/devcontainer-cli/internal/pkg/output"
)

func createListCommand() *cobra.Command {
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "List devcontainers",
		Long:  "Lists running devcontainers",
		RunE: func(cmd *cobra.Command, args []string) error {
			outputFormat, query, err := output.GetOutputAndQueryValues(cmd, `[].{name: devcontainerName, containerName: containerName, containerID: containerID, localFolderPath: localFolderPath}`)
			if err != nil {
				return err
			}
			devcontainerList, err := devcontainers.ListDevcontainers()
			if err != nil {
				return err
			}

			err = output.OutputResult(os.Stdout, devcontainerList, outputFormat, query, []string{"name", "containerName", "containerID", "localFolderPath"})
			if err != nil {
				return fmt.Errorf("error outputting result: %s", err)
			}
			return nil
		},
	}

	output.AddOutputAndQueryFlags(cmdList)
	return cmdList
}

func createShowCommand() *cobra.Command {
	var argDevcontainerName string
	cmd := &cobra.Command{
		Use:   "show --name <name>",
		Short: "Show devcontainer info",
		Long:  "Show information about a running dev container",
		RunE: func(cmd *cobra.Command, args []string) error {
			outputFormat, query, err := output.GetOutputAndQueryValues(cmd, `[].{name: devcontainerName, containerName: containerName, containerID: containerID, localFolderPath: localFolderPath}`)
			if err != nil {
				return err
			}
			devcontainerList, err := devcontainers.ListDevcontainers()
			if err != nil {
				return err
			}
			containerIDOrName := argDevcontainerName

			// Get container ID
			for _, devcontainer := range devcontainerList {
				if devcontainer.ContainerName == containerIDOrName ||
					devcontainer.DevcontainerName == containerIDOrName ||
					devcontainer.ContainerID == containerIDOrName {

					wrapped := []devcontainers.DevcontainerInfo{devcontainer}
					err = output.OutputResult(os.Stdout, wrapped, outputFormat, query, []string{"name", "containerName", "containerID", "localFolderPath"})
					if err != nil {
						return fmt.Errorf("error outputting result: %s", err)
					}
					return nil
				}
			}

			return fmt.Errorf("failed to find a matching (running) dev container for %q", containerIDOrName)
		},
	}
	cmd.Flags().StringVarP(&argDevcontainerName, "name", "n", "", "name of dev container to show")
	output.AddOutputAndQueryFlags(cmd)

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
					return fmt.Errorf("failed to find a matching (running) dev container for %q", containerIDOrName)
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
					return fmt.Errorf("invalid option")
				}
				containerID = devcontainerList[selection].ContainerID
			} else {
				devcontainerPath := argDevcontainerPath
				// TODO - update to check for devcontainers in the path ancestry
				// Can't just check up the path for a .devcontainer folder as the container might
				// have been created via repository containers (https://github.com/microsoft/vscode-dev-containers/tree/main/repository-containers)
				devcontainer, err := devcontainers.GetClosestPathMatchForPath(devcontainerList, devcontainerPath)
				if err != nil {
					return err
				}
				containerID = devcontainer.ContainerID
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
