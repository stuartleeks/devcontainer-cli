package wsl

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// IsWsl returns true if running under WSL
func IsWsl() bool {
	_, exists := os.LookupEnv("WSL_DISTRO_NAME")
	return exists
}

// ConvertWslPathToWindowsPath converts a WSL path to the corresponding \\wsl$\... path for access from Windows
func ConvertWslPathToWindowsPath(path string) (string, error) {
	cmd := exec.Command("wslpath", "-w", path)

	buf, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("Error running wslpath (for %q): %s", path, err)
	}
	return strings.TrimSpace(string(buf)), nil
}

// ConvertWslPathToWindowsPath converts a WSL path to the corresponding \\wsl$\... path for access from Windows
func ConvertWindowsPathToWslPath(path string) (string, error) {
	cmd := exec.Command("wslpath", "-u", path)

	buf, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("Error running wslpath (for %q): %s", path, err)
	}
	return strings.TrimSpace(string(buf)), nil
}
