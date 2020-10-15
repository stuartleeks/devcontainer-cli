package git

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"
)

// GetTopLevelPath returns the top-level folder for the git-repo that contains path, or empty string if not a repo
func GetTopLevelPath(path string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = path

	buf, err := cmd.Output()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				if status.ExitStatus() == 128 {
					// exit code 128 indictates not a git repo
					return "", nil
				}
			}
			return "", fmt.Errorf("Error git rev-parse --show-toplevel: %s", err)
		}
		return "", fmt.Errorf("Error git rev-parse --show-toplevel: %s", err)
	}
	return strings.TrimSpace(string(buf)), nil
}
