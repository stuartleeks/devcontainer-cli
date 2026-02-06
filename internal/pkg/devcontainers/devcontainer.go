package devcontainers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func getDevContainerJsonPath(folderPath string) (string, error) {
	pathsToTest := []string{".devcontainer/devcontainer.json", ".devcontainer.json"}

	for _, path := range pathsToTest {
		devcontainerJsonPath := filepath.Join(folderPath, path)
		devContainerJsonInfo, err := os.Stat(devcontainerJsonPath)
		if err == nil && !devContainerJsonInfo.IsDir() {
			return devcontainerJsonPath, nil
		}
	}

	return "", fmt.Errorf("devcontainer.json not found. Looked for %s", strings.Join(pathsToTest, ","))
}

func FindDevContainerInAncestorPaths(folderPath string) (string, error) {
	currentPath, err := filepath.Abs(folderPath)
	if err != nil {
		return "", fmt.Errorf("error getting absolute path: %w", err)
	}

	for {
		// Check if devcontainer.json exists in current path
		_, err := getDevContainerJsonPath(currentPath)
		if err == nil {
			return currentPath, nil
		}

		// Check if this is a git repository root
		gitPath := filepath.Join(currentPath, ".git")
		gitInfo, gitErr := os.Stat(gitPath)
		isGitRoot := gitErr == nil && gitInfo.IsDir()

		// If we're at a git root, stop searching (we already checked this folder)
		if isGitRoot {
			return "", fmt.Errorf("devcontainer.json not found in ancestor paths")
		}

		// Move to parent directory
		parentPath := filepath.Dir(currentPath)

		// Check if we've reached the root (parent is same as current)
		if parentPath == currentPath {
			return "", fmt.Errorf("devcontainer.json not found in ancestor paths")
		}

		currentPath = parentPath
	}
}
