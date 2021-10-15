package devcontainers

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"

	"github.com/stuartleeks/devcontainer-cli/internal/pkg/git"
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
		if wsl.HasWslPathPrefix(folderPath) {
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
	devcontainerPath, err := getDefaultWorkspaceFolderForPath(folderPath)
	if err != nil {
		return "", fmt.Errorf("Error getting default workspace path: %s", err)
	}
	return fmt.Sprintf("/workspaces/%s", devcontainerPath), nil
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

func getDefaultWorkspaceFolderForPath(path string) (string, error) {

	// get the git repo-root
	rootPath, err := git.GetTopLevelPath(path)
	if err != nil {
		return "", err
	}
	if rootPath == "" {
		// not a git repo, default to path
		rootPath = path
	}

	// get parent to root
	rootParent, _ := filepath.Split(rootPath)

	// return path relative to rootParent
	relativePath, err := filepath.Rel(rootParent, path)
	if err != nil {
		return "", err
	}
	return relativePath, nil
}
