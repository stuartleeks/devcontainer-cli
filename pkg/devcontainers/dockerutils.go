package devcontainers

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/stuartleeks/devcontainer-cli/internal/pkg/wsl"
)

// Devcontainer names are a derived property.
// For Dockerfile-based devcontainers this is derived from the image name
// E.g. vsc-vscode-remote-test-dockerfile-16020a1c793898f64e7c9cff44437369 => vscode-remote-test-dockerfile
// For dockercompose-based devcontainers this is derived from the com.docker.compose.project
// and com.docker.compose.service labels with a slash separator. E.g. vscode-remote-test-dockercompose_devcontainer\web

// DevcontainerInfo holds details about a devcontainer
type DevcontainerInfo struct {
	ContainerID      string
	ContainerName    string
	DevcontainerName string
	LocalFolderPath  string
}

const (
	listPartID                     int = 0
	listPartLocalFolder            int = 1
	listPartComposeProject         int = 2
	listPartComposeService         int = 3
	listPartComposeContainerNumber int = 4
	listPartContainerName          int = 5
)

var _ = listPartComposeContainerNumber

// ListDevcontainers returns a list of devcontainers
func ListDevcontainers() ([]DevcontainerInfo, error) {
	cmd := exec.Command("docker", "ps", "--format", "{{.ID}}|{{.Label \"vsch.local.folder\"}}|{{.Label \"com.docker.compose.project\"}}|{{.Label \"com.docker.compose.service\"}}|{{.Label \"com.docker.compose.container-number\"}}|{{.Names}}")

	output, err := cmd.Output()
	if err != nil {
		return []DevcontainerInfo{}, fmt.Errorf("Failed to read docker stdout: %v", err)
	}

	reader := bytes.NewReader(output)
	scanner := bufio.NewScanner(reader)
	if scanner == nil {
		return []DevcontainerInfo{}, fmt.Errorf("Failed to parse stdout")
	}
	devcontainers := []DevcontainerInfo{}
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "|")
		name := parts[listPartLocalFolder]
		if name == "" {
			// No local folder => use dockercompose parts
			name = fmt.Sprintf("%s/%s", parts[listPartComposeProject], parts[listPartComposeService])
		} else {
			// get the last path segment for the name
			if index := strings.LastIndexAny(name, "/\\"); index >= 0 {
				name = name[index+1:]
			}
		}
		devcontainer := DevcontainerInfo{
			ContainerID:      parts[listPartID],
			ContainerName:    parts[listPartContainerName],
			LocalFolderPath:  parts[listPartLocalFolder],
			DevcontainerName: name,
		}
		devcontainers = append(devcontainers, devcontainer)
	}
	return devcontainers, nil
}

// GetLocalFolderFromDevContainer looks up the local (host) folder name from the container labels
func GetLocalFolderFromDevContainer(containerIDOrName string) (string, error) {
	//docker inspect cool_goldberg --format '{{ index .Config.Labels "vsch.local.folder" }}'

	cmd := exec.Command("docker", "inspect", containerIDOrName, "--format", "{{ index .Config.Labels \"vsch.local.folder\" }}")

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("Failed to read docker stdout: %v", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// GetContainerIDForPath returns the ID of the running container that matches the path
func GetContainerIDForPath(devcontainerPath string) (string, error) {
	if devcontainerPath == "" {
		devcontainerPath = "."
	}
	absPath, err := filepath.Abs(devcontainerPath)
	if err != nil {
		return "", fmt.Errorf("Error handling path %q: %s", devcontainerPath, err)
	}

	windowsPath := absPath
	if wsl.IsWsl() {
		var err error
		windowsPath, err = wsl.ConvertWslPathToWindowsPath(windowsPath)
		if err != nil {
			return "", err
		}
	}

	devcontainerList, err := ListDevcontainers()
	if err != nil {
		return "", fmt.Errorf("Error getting container list: %s", err)
	}

	for _, devcontainer := range devcontainerList {
		if devcontainer.LocalFolderPath == windowsPath {
			containerID := devcontainer.ContainerID
			return containerID, nil
		}
	}
	return "", fmt.Errorf("Could not find running container for path %q", devcontainerPath)
}

func ExecInDevContainer(containerIDOrName string, workDir string, args []string) error {

	containerID := ""
	devcontainerList, err := ListDevcontainers()
	if err != nil {
		return err
	}

	for _, devcontainer := range devcontainerList {
		if devcontainer.ContainerName == containerIDOrName ||
			devcontainer.DevcontainerName == containerIDOrName ||
			devcontainer.ContainerID == containerIDOrName {
			containerID = devcontainer.ContainerID
			break
		}
	}

	if containerID == "" {
		return fmt.Errorf("Failed to find a matching (running) dev container for %q", containerIDOrName)
	}

	localPath, err := GetLocalFolderFromDevContainer(containerID)
	if err != nil {
		return err
	}

	if workDir == "" {
		workDir, err = GetWorkspaceMountPath(localPath)
		if err != nil {
			return err
		}
	}

	wslPath := localPath
	if strings.HasPrefix(wslPath, "\\\\wsl$") && wsl.IsWsl() {
		wslPath, err = wsl.ConvertWindowsPathToWslPath(wslPath)
		if err != nil {
			return fmt.Errorf("error converting path: %s", err)
		}
	}

	devcontainerJSONPath := path.Join(wslPath, ".devcontainer/devcontainer.json")
	userName, err := GetDevContainerUserName(devcontainerJSONPath)
	if err != nil {
		return err
	}

	sshAuthSockValue, err := getSshAuthSockValue(containerID)
	if err != nil {
		// output error and continue without SSH_AUTH_SOCK value
		sshAuthSockValue = ""
		fmt.Printf("Warning: Failed to get SSH_AUTH_SOCK value: %s\n", err)
		fmt.Println("Continuing without setting SSH_AUTH_SOCK...")
	}

	dockerArgs := []string{"exec", "-it", "--workdir", workDir}
	if userName != "" {
		dockerArgs = append(dockerArgs, "--user", userName)
	}
	if sshAuthSockValue != "" {
		dockerArgs = append(dockerArgs, "--env", "SSH_AUTH_SOCK="+sshAuthSockValue)
	}
	dockerArgs = append(dockerArgs, containerID)
	dockerArgs = append(dockerArgs, args...)

	dockerCmd := exec.Command("docker", dockerArgs...)
	dockerCmd.Stdin = os.Stdin
	dockerCmd.Stdout = os.Stdout

	err = dockerCmd.Start()
	if err != nil {
		return fmt.Errorf("Exec: start error: %s", err)
	}
	err = dockerCmd.Wait()
	if err != nil {
		return fmt.Errorf("Exec: wait error: %s", err)
	}
	return nil
}

// getSshAuthSockValue returns the value to use for the SSH_AUTH_SOCK env var when exec'ing into the container, or empty string if no value is found
func getSshAuthSockValue(containerID string) (string, error) {

	// If the host has SSH_AUTH_SOCK set then VS Code spins up forwarding for key requests
	// inside the dev container to the SSH agent on the host.

	hostSshAuthSockValue := os.Getenv("SSH_AUTH_SOCK")
	if hostSshAuthSockValue == "" {
		// Nothing to see, move along
		return "", nil
	}

	// Host has SSH_AUTH_SOCK set, so expecting the dev container to have forwarding set up
	// Find the latest /tmp/vscode-ssh-auth-<...>.sock

	dockerArgs := []string{"exec", containerID, "bash", "-c", "ls -t -d -1 \"${TMPDIR:-/tmp}\"/vscode-ssh-auth-*"}

	dockerCmd := exec.Command("docker", dockerArgs...)
	buf, err := dockerCmd.CombinedOutput()
	if err != nil {
		errMessage := string(buf)
		return "", fmt.Errorf("Docker exec error: %s (%s)", err, strings.TrimSpace(errMessage))
	}

	output := string(buf)
	lines := strings.Split(output, "\n")
	if len(lines) <= 0 {
		return "", nil
	}
	return strings.TrimSpace(lines[0]), nil
}
