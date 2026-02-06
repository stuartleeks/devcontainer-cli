package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stuartleeks/devcontainer-cli/internal/pkg/devcontainers"
	"github.com/stuartleeks/devcontainer-cli/internal/pkg/wsl"
)

func createOpenInCodeCommand() *cobra.Command {
	var subFolder string
	cmd := &cobra.Command{
		Use:   "open-in-code <path>",
		Short: "open the specified path devcontainer project in VS Code",
		Long:  "Open the specified path (containing a .devcontainer folder in VS Code",
		RunE: func(cmd *cobra.Command, args []string) error {
			return launchDevContainer(cmd, "code", args, subFolder)
		},
	}
	cmd.Flags().StringVarP(&subFolder, "sub-folder", "s", "", "sub-folder within the devcontainer to open")
	return cmd
}
func createOpenInCodeInsidersCommand() *cobra.Command {
	var subFolder string
	cmd := &cobra.Command{
		Use:   "open-in-code-insiders <path>",
		Short: "open the specified path devcontainer project in VS Code Insiders",
		Long:  "Open the specified path (containing a .devcontainer folder in VS Code Insiders",
		RunE: func(cmd *cobra.Command, args []string) error {
			return launchDevContainer(cmd, "code-insiders", args, subFolder)
		},
	}
	cmd.Flags().StringVarP(&subFolder, "sub-folder", "s", "", "sub-folder within the devcontainer to open")
	return cmd
}

func launchDevContainer(cmd *cobra.Command, appBase string, args []string, subFolder string) error {
	if len(args) > 1 {
		return cmd.Usage()
	}
	path := "." // default to current directory
	if len(args) >= 1 {
		path = args[0]
	}

	// allow the command to be invoked from a subfolder of the folder containing the devcontainer.json by searching ancestor paths for a devcontainer.json
	devcontainerPath, err := devcontainers.FindDevContainerInAncestorPaths(path)
	if err != nil {
		return fmt.Errorf("error finding devcontainer.json: %s", err)
	}
	fmt.Printf("Found devcontainer.json at %s\n", devcontainerPath)

	// If subFolder is not empty then use that as the path to open in VS Code (note it should be relative to the folder containing the devcontainer.json/.devcontainer folder)
	// If subFolder is empty and the current folder is contained within the folder containing the devcontainer.json, then open the current folder in VS Code, otherwise open the folder containing the devcontainer.json in VS Code
	if subFolder == "" {
		absDevContainerPath, err := filepath.Abs(devcontainerPath)
		if err != nil {
			return fmt.Errorf("error getting absolute path: %s", err)
		}
		absCurrentPath, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("error getting absolute path: %s", err)
		}
		if absCurrentPath != absDevContainerPath && strings.HasPrefix(absCurrentPath, absDevContainerPath) {
			subFolder, err = filepath.Rel(absDevContainerPath, absCurrentPath)
			if err != nil {
				return fmt.Errorf("error getting relative path: %s", err)
			}
			fmt.Printf("Opening subfolder %s within devcontainer\n", subFolder)
		}
	}

	launchURI, err := devcontainers.GetDevContainerURI(devcontainerPath, subFolder)
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
