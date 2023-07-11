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
