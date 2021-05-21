package main

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/stuartleeks/devcontainer-cli/internal/pkg/wsl"
	"github.com/stuartleeks/devcontainer-cli/pkg/devcontainers"
)

func createOpenInCodeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "open-in-code <path>",
		Short: "open the specified path devcontainer project in VS Code",
		Long:  "Open the specified path (containing a .devcontainer folder in VS Code",
		RunE: func(cmd *cobra.Command, args []string) error {
			return launchDevContainer(cmd, "code", args)
		},
	}
	return cmd
}
func createOpenInCodeInsidersCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "open-in-code-insiders <path>",
		Short: "open the specified path devcontainer project in VS Code Insiders",
		Long:  "Open the specified path (containing a .devcontainer folder in VS Code Insiders",
		RunE: func(cmd *cobra.Command, args []string) error {
			return launchDevContainer(cmd, "code-insiders", args)
		},
	}
	return cmd
}

func launchDevContainer(cmd *cobra.Command, appBase string, args []string) error {
	if len(args) != 1 {
		return cmd.Usage()
	}
	path := args[0]

	launchURI, err := devcontainers.GetDevContainerURI(path)
	if err != nil {
		return err
	}
	var execCmd *exec.Cmd
	if wsl.IsWsl() {
		execCmd = exec.Command("cmd.exe", "/C", appBase+".cmd", "--folder-uri="+launchURI)
	} else {
		execCmd = exec.Command(appBase, "--folder-uri="+launchURI)
	}
	output, err := execCmd.Output()
	fmt.Println(string(output))
	if err != nil {
		return err
	}
	return nil
}
