package devcontainers

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/stuartleeks/devcontainer-cli/internal/pkg/wsl"
)

// GetDevContainerURI gets the devcontainer URI for a folder to launch using the VS Code --folder-uri switch
func GetDevContainerURI(folderPath string) (string, error) {

	absPath, err := filepath.Abs(folderPath)
	if err != nil {
		return "", fmt.Errorf("Error handling path %q: %s", folderPath, err)
	}

	launchPath := absPath
	if wsl.IsWsl() {
		var err error
		launchPath, err = wsl.ConvertWslPathToWindowsPath(launchPath)
		if err != nil {
			return "", err
		}
	}

	launchPathHex := convertToHexString(launchPath)
	workspaceMountPath, err := GetWorkspaceMountPath(absPath)
	if err != nil {
		return "", err
	}
	uri := fmt.Sprintf("vscode-remote://dev-container+%s%s", launchPathHex, workspaceMountPath)

	return uri, nil
}

func convertToHexString(input string) string {
	return hex.EncodeToString([]byte(input))
}

// GetWorkspaceMountPath returns the devcontainer mount path for the devcontainer in the specified folder
func GetWorkspaceMountPath(folderPath string) (string, error) {
	// TODO - consider how to support repository-containers (https://github.com/microsoft/vscode-remote-release/issues/3218)

	// If we're called from WSL we want a WSL Path but will also handle a Windows Path
	if wsl.IsWsl() {
		if strings.HasPrefix(folderPath, "\\\\wsl$\\") {
			convertedPath, err := wsl.ConvertWindowsPathToWslPath(folderPath)
			if err != nil {
				return "", err
			}
			folderPath = convertedPath
		}
	}

	devcontainerDefinitionPath := filepath.Join(folderPath, ".devcontainer/devcontainer.json")
	buf, err := ioutil.ReadFile(devcontainerDefinitionPath)
	if err != nil {
		return "", fmt.Errorf("Error loading devcontainer definition: %s", err)
	}

	workspaceMountPath, err := getWorkspaceMountPathFromDevcontainerDefinition(buf)
	if err != nil {
		return "", fmt.Errorf("Error parsing devcontainer definition: %s", err)
	}
	if workspaceMountPath != "" {
		return workspaceMountPath, nil
	}

	// No `workspaceFolder` found in devcontainer.json - use default
	_, folderName := filepath.Split(folderPath)
	return fmt.Sprintf("/workspaces/%s", folderName), nil
}

// TODO: add tests (and implementation) to handle JSON parsing with comments
// Current implementation doesn't handle
//  - block comments
//  - the value split on a new line from the property name

func getWorkspaceMountPathFromDevcontainerDefinition(definition []byte) (string, error) {
	r, err := regexp.Compile("(?m)^\\s*\"workspaceFolder\"\\s*:\\s*\"(.*)\"")
	if err != nil {
		return "", fmt.Errorf("Error compiling regex: %s", err)
	}
	matches := r.FindSubmatch(definition)
	if len(matches) == 2 {
		return string(matches[1]), nil
	}
	return "", nil
}
