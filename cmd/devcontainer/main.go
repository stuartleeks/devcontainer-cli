package main

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/stuartleeks/devcontainer-cli/internal/pkg/devcontainers"
)

func main() {
	var listIncludeContainerNames bool
	var listVerbose bool

	var rootCmd = &cobra.Command{Use: "devcontainer"}

	var cmdList = &cobra.Command{
		Use:   "list",
		Short: "List devcontainers",
		Long:  "Lists devcontainers that are currently running",
		Run: func(cmd *cobra.Command, args []string) {
			if listIncludeContainerNames && listVerbose {
				fmt.Println("Can't use both verbose and include-container-names")
				os.Exit(1)
			}
			devcontainers, err := devcontainers.ListDevcontainers()
			if err != nil {
				fmt.Printf("Error: %v", err)
				os.Exit(1)
			}
			if listVerbose {
				sort.Slice(devcontainers, func(i, j int) bool { return devcontainers[i].DevcontainerName < devcontainers[j].DevcontainerName })

				w := new(tabwriter.Writer)
				// minwidth, tabwidth, padding, padchar, flags
				w.Init(os.Stdout, 8, 8, 0, '\t', 0)
				defer w.Flush()

				fmt.Fprintf(w, "\n %s\t%s\t", "Devcontainer name", "Container name")
				fmt.Fprintf(w, "\n %s\t%s\t", "----", "----")

				for _, devcontainer := range devcontainers {
					fmt.Fprintf(w, "\n %s\t%s\t", devcontainer.DevcontainerName, devcontainer.ContainerName)
				}
				fmt.Fprintln(w)
				return
			}
			names := []string{}
			for _, devcontainer := range devcontainers {
				names = append(names, devcontainer.DevcontainerName)
			}
			if listIncludeContainerNames {
				for _, devcontainer := range devcontainers {
					names = append(names, devcontainer.ContainerName)
				}
			}
			sort.Strings(names)
			for _, name := range names {
				fmt.Println(name)
			}
		},
	}
	cmdList.Flags().BoolVar(&listIncludeContainerNames, "include-container-names", false, "Also include container names in the list")
	cmdList.Flags().BoolVarP(&listVerbose, "verbose", "v", false, "Verbose output")

	rootCmd.AddCommand(cmdList)
	rootCmd.Execute()
}
