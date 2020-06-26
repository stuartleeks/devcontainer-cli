package main

import (
	"fmt"
	"os"
	"os/exec"
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

func createExecCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exec DEVCONTAINER_NAME COMMAND [args...]",
		Short: "Execute a command in a devcontainer",
		Long:  "Execute a command in a devcontainer, similar to `docker exec`. Pass `?` as DEVCONTAINER_NAME to be prompted.",
		RunE: func(cmd *cobra.Command, args []string) error {

			if len(args) < 2 {
				return cmd.Usage()
			}

			devcontainerName := args[0]
			devcontainers, err := devcontainers.ListDevcontainers()
			if err != nil {
				return err
			}

			containerID := ""
			if devcontainerName == "?" {
				// prompt user
				fmt.Println("Specify the devcontainer to use:")
				for index, devcontainer := range devcontainers {
					fmt.Printf("%4d: %s (%s)\n", index, devcontainer.DevcontainerName, devcontainer.ContainerName)
				}
				selection := -1
				_, _ = fmt.Scanf("%d", &selection)
				if selection < 0 || selection >= len(devcontainers) {
					return fmt.Errorf("Invalid option")
				}
				containerID = devcontainers[selection].ContainerID
			} else {
				for _, devcontainer := range devcontainers {
					if devcontainer.ContainerName == devcontainerName || devcontainer.DevcontainerName == devcontainerName {
						containerID = devcontainer.ContainerID
						break
					}
				}
				if containerID == "" {
					return cmd.Usage()
				}
			}

			dockerArgs := []string{"exec", "-it", containerID}
			dockerArgs = append(dockerArgs, args[1:]...)

			dockerCmd := exec.Command("docker", dockerArgs...)
			dockerCmd.Stdin = os.Stdin
			dockerCmd.Stdout = os.Stdout

			err = dockerCmd.Start()
			if err != nil {
				return fmt.Errorf("Exec: start error: %s\n", err)
			}
			err = dockerCmd.Wait()
			if err != nil {
				return fmt.Errorf("Exec: wait error: %s\n", err)
			}
			return nil
		},
		Args:                  cobra.ArbitraryArgs,
		DisableFlagParsing:    true,
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
				names = append(names, devcontainer.ContainerName)
			}
			sort.Strings(names)
			return names, cobra.ShellCompDirectiveNoFileComp
		},
	}
	return cmd
}
