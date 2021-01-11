package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/stuartleeks/devcontainer-cli/internal/pkg/devcontainers"
	"github.com/stuartleeks/devcontainer-cli/internal/pkg/wsl"
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

	cmd := &cobra.Command{
		Use:   "exec [--name <name>| --path <path> | --prompt ] [<command> [<args...>]] (command will default to /bin/bash if none provided)",
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
			if sourceCount > 0 {
				fmt.Println("Can specify at most one of --name/--path/--prompt")
				return cmd.Usage()
			}

			containerID := ""
			devcontainerList, err := devcontainers.ListDevcontainers()
			if err != nil {
				return err
			}
			if argDevcontainerName != "" {
				var devcontainerName string
				devcontainerName = argDevcontainerName
				for _, devcontainer := range devcontainerList {
					if devcontainer.ContainerName == devcontainerName || devcontainer.DevcontainerName == devcontainerName {
						containerID = devcontainer.ContainerID
						break
					}
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
				if devcontainerPath == "" {
					devcontainerPath = "."
				}
				absPath, err := filepath.Abs(devcontainerPath)
				if err != nil {
					return fmt.Errorf("Error handling path %q: %s", devcontainerPath, err)
				}

				windowsPath := absPath
				if wsl.IsWsl() {
					var err error
					windowsPath, err = wsl.ConvertWslPathToWindowsPath(windowsPath)
					if err != nil {
						return err
					}
				}
				for _, devcontainer := range devcontainerList {
					if devcontainer.LocalFolderPath == windowsPath {
						containerID = devcontainer.ContainerID
						break
					}
				}
			}

			if containerID == "" {
				fmt.Println("Failed to find dev container")
				return cmd.Usage()
			}

			localPath, err := devcontainers.GetLocalFolderFromDevContainer(containerID)
			if err != nil {
				return err
			}

			mountPath, err := devcontainers.GetWorkspaceMountPath(localPath)
			if err != nil {
				return err
			}

			wslPath := localPath
			if strings.HasPrefix(wslPath, "\\\\wsl$") && wsl.IsWsl() {
				wslPath, err = wsl.ConvertWindowsPathToWslPath(wslPath)
				if err != nil {
					return fmt.Errorf("error converting path: %s", err)
				}
			}

			devcontainerJSONPath := path.Join(wslPath, ".devcontainer/devcontainer.json")
			userName, err := devcontainers.GetDevContainerUserName(devcontainerJSONPath)
			if err != nil {
				return err
			}

			dockerArgs := []string{"exec", "-it", "--workdir", mountPath}
			if userName != "" {
				dockerArgs = append(dockerArgs, "--user", userName)
			}
			dockerArgs = append(dockerArgs, containerID)
			dockerArgs = append(dockerArgs, args...)

			dockerCmd := exec.Command("docker", dockerArgs...)
			dockerCmd.Stdin = os.Stdin
			dockerCmd.Stdout = os.Stdout

			err = dockerCmd.Start()
			if err != nil {
				return fmt.Errorf("Exec: start error: %s", err)
			}
			err = dockerCmd.Wait()
			if err != nil {
				return fmt.Errorf("Exec: wait error: %s", err)
			}
			return nil
		},
		Args: cobra.ArbitraryArgs,
		// DisableFlagParsing:    true,
		DisableFlagsInUseLine: true,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			// only completing the first arg  (devcontainer name)
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
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
		},
	}
	cmd.Flags().StringVarP(&argDevcontainerName, "name", "n", "", "name of dev container to exec into")
	cmd.Flags().StringVarP(&argDevcontainerPath, "path", "", "", "path containing the dev container to exec into")
	cmd.Flags().BoolVarP(&argPromptForDevcontainer, "prompt", "", false, "prompt for the dev container to exec into")
	return cmd
}
