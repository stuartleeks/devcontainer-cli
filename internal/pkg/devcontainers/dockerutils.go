package devcontainers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/stuartleeks/devcontainer-cli/internal/pkg/terminal"
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
		localPath := parts[listPartLocalFolder]
		if localPath == "" {
			// not a dev container
			continue
		}
		if wsl.HasWslPathPrefix(localPath) && wsl.IsWsl() {
			localPath, err = wsl.ConvertWindowsPathToWslPath(localPath)
			if err != nil {
				return []DevcontainerInfo{}, fmt.Errorf("error converting path: %s", err)
			}
		}
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
			LocalFolderPath:  localPath,
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

// DockerMount represents mount info from Docker output
type DockerMount struct {
	Source      string `json:"Source"`
	Destination string `json:"Destination"`
}

// SourceInfo holds properties about the source mounted in a dev container
type SourceInfo struct {
	DevcontainerFolder string
	DockerMount        DockerMount
}

// getMountFolderFromFolder walks up the path hierarchy checking for a git repo. If found that is returned, if not the original folder is used
func getMountFolderFromFolder(folder string) (string, error) {

	for currentFolder := folder; currentFolder != ""; currentFolder = path.Dir(currentFolder) {
		info, err := os.Stat(path.Join(currentFolder, ".git"))
		if os.IsNotExist(err) {
			continue // walk up the hierarchy
		}
		if err != nil {
			// other error - e.g. permissions issue
			// give up the search
			break
		}
		if info.IsDir() {
			// We've found a .git folder - use it
			return currentFolder, nil
		}

	}
	// didn't find anything - use original folder
	return folder, nil
}

// GetSourceInfoFromDevContainer inspects the specified container and returns the SourceInfo
func GetSourceInfoFromDevContainer(containerIDOrName string) (SourceInfo, error) {
	localPath, err := GetLocalFolderFromDevContainer(containerIDOrName)
	if err != nil {
		return SourceInfo{}, err
	}

	if wsl.HasWslPathPrefix(localPath) && wsl.IsWsl() {
		localPath, err = wsl.ConvertWindowsPathToWslPath(localPath)
		if err != nil {
			return SourceInfo{}, fmt.Errorf("error converting path: %s", err)
		}
	}

	mountFolder, err := getMountFolderFromFolder(localPath)
	if err != nil {
		return SourceInfo{}, fmt.Errorf("search for mount folder failed: %s", err)
	}

	cmd := exec.Command("docker", "inspect", containerIDOrName, "--format", fmt.Sprintf("{{ range .Mounts }}{{if eq .Source \"%s\"}}{{json .}}{{end}}{{end}}", mountFolder))

	output, err := cmd.Output()
	if err != nil {
		return SourceInfo{}, fmt.Errorf("Failed to read docker stdout: %v", err)
	}

	var mount DockerMount
	err = json.Unmarshal(output, &mount)
	if err != nil {
		return SourceInfo{}, fmt.Errorf("failed to parse JSON getting mount folder for container %q (path=%q): %s", containerIDOrName, mountFolder, err)
	}

	return SourceInfo{
		DevcontainerFolder: localPath,
		DockerMount:        mount,
	}, nil
}

type byLocalPathLength []DevcontainerInfo

func (s byLocalPathLength) Len() int {
	return len(s)
}
func (s byLocalPathLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byLocalPathLength) Less(i, j int) bool {
	return len(s[i].LocalFolderPath) < len(s[j].LocalFolderPath)
}

// GetClosestPathMatchForPath returns the dev container with the closes match to the specified path
func GetClosestPathMatchForPath(devContainers []DevcontainerInfo, devcontainerPath string) (DevcontainerInfo, error) {
	if devcontainerPath == "" {
		devcontainerPath = "."
	}
	absPath, err := filepath.Abs(devcontainerPath)
	if err != nil {
		return DevcontainerInfo{}, fmt.Errorf("Error handling path %q: %s", devcontainerPath, err)
	}

	matchingPaths := byLocalPathLength{}
	for _, devcontainer := range devContainers {
		// Treat as match if the specified path is within the devcontainer path
		testPath := devcontainer.LocalFolderPath
		if wsl.IsWsl() && wsl.HasWslPathPrefix(testPath) {
			testPath, err = wsl.ConvertWindowsPathToWslPath(testPath)
			fmt.Println("Converted to..")
			if err != nil {
				return DevcontainerInfo{}, fmt.Errorf("Error converting path from dev container list (%q): %s", testPath, err)
			}
		}
		if strings.HasPrefix(absPath, testPath) {
			matchingPaths = append(matchingPaths, devcontainer)
		}
	}
	if len(matchingPaths) == 0 {
		return DevcontainerInfo{}, fmt.Errorf("Could not find running container for path %q", devcontainerPath)
	}

	// return longest prefix match
	sort.Sort(matchingPaths)
	return matchingPaths[len(matchingPaths)-1], nil
}

func ExecInDevContainer(containerID string, workDir string, args []string) error {

	statusWriter := &terminal.UpdatingStatusWriter{}

	sourceInfo, err := GetSourceInfoFromDevContainer(containerID)
	if err != nil {
		return fmt.Errorf("failed to get source mount: %s", err)
	}
	localPath := sourceInfo.DevcontainerFolder

	statusWriter.Printf("Getting user name")
	devcontainerJSONPath := path.Join(localPath, ".devcontainer/devcontainer.json")
	userName, err := GetDevContainerUserName(devcontainerJSONPath)
	if err != nil {
		return err
	}

	statusWriter.Printf("Checking for SSH_AUTH_SOCK")
	sshAuthSockValue, err := getSshAuthSockValue(containerID)
	if err != nil {
		// output error and continue without SSH_AUTH_SOCK value
		sshAuthSockValue = ""
		fmt.Printf("Warning: Failed to get SSH_AUTH_SOCK value: %s\n", err)
		fmt.Println("Continuing without setting SSH_AUTH_SOCK...")
	}

	statusWriter.Printf("Getting container PATH")
	containerPath, err := getContainerEnvVar(containerID, "PATH")
	if err == nil {
		// Got the PATH
		statusWriter.Printf("Getting code server path")
		vscodeServerPath, err := getVscodeServerPath(containerID)
		if err == nil {
			// Got the VS Code server location - add bin subfolder to PATH
			containerPath = strings.TrimSpace(containerPath)
			containerPath = fmt.Sprintf("%s/bin:%s", vscodeServerPath, containerPath)
		} else {
			// output error and continue without adding to PATH value
			fmt.Printf("Warning: Failed to get VS Code server location: %s\n", err)
			fmt.Println("Continuing without adding VS Code Server to PATH...")
		}
	} else {
		// output error and continue without adding to PATH value
		containerPath = ""
		fmt.Printf("Warning: Failed to get PATH value for container: %s\n", err)
		fmt.Println("Continuing without overriding PATH...")
	}

	statusWriter.Printf("Getting code IPC SOCK")
	ipcSock, err := getVscodeIpcSock(containerID)
	if err != nil {
		ipcSock = ""
		fmt.Printf("Warning; Failed to get VS Code IPC SOCK: %s\n", err)
		fmt.Println("Continuing without setting VSCODE_IPC_HOOK_CLI...")
	}

	mountPath := sourceInfo.DockerMount.Destination
	if workDir == "" {
		workDir = sourceInfo.DevcontainerFolder
	} else if !filepath.IsAbs(workDir) {

		// Convert to absolute (local) path
		// This takes into account current directory (potentially within the dev container path)
		// We'll convert local to container path below
		workDir, err = filepath.Abs(workDir)
		if err != nil {
			return err
		}
	}

	statusWriter.Printf("Test container path")
	containerPathExists, err := testContainerPathExists(containerID, workDir)
	if err != nil {
		return fmt.Errorf("error checking container path: %s", err)
	}
	if !containerPathExists {
		// path not found - try converting from local path
		// ? Should we check here that the workDir has mountPath as a prefix?
		devContainerRelativePath, err := filepath.Rel(sourceInfo.DockerMount.Source, workDir)
		if err != nil {
			return fmt.Errorf("error getting path relative to mount dir: %s", err)
		}
		workDir = filepath.Join(mountPath, devContainerRelativePath)
	}

	statusWriter.Printf("Starting exec session\n") // newline to put container shell at start of line
	dockerArgs := []string{"exec", "-it", "--workdir", workDir}
	if userName != "" {
		dockerArgs = append(dockerArgs, "--user", userName)
	}
	if sshAuthSockValue != "" {
		dockerArgs = append(dockerArgs, "--env", "SSH_AUTH_SOCK="+sshAuthSockValue)
	}
	if containerPath != "" {
		dockerArgs = append(dockerArgs, "--env", "PATH="+containerPath)
	}
	if ipcSock != "" {
		dockerArgs = append(dockerArgs, "--env", "VSCODE_IPC_HOOK_CLI="+ipcSock)
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

	return getLatestFileMatch(containerID, "\"${TMPDIR:-/tmp}\"/vscode-ssh-auth-*")
}

func getVscodeServerPath(containerID string) (string, error) {
	return getLatestFileMatch(containerID, "/vscode/vscode-server/bin/x64/*")
}
func getVscodeIpcSock(containerID string) (string, error) {
	return getLatestFileMatch(containerID, "\"${TMPDIR:-/tmp}\"/vscode-ipc-*")
}

// getLatestFileMatch lists files matching `pattern` in the container and returns the latest filename
func getLatestFileMatch(containerID string, pattern string) (string, error) {

	dockerArgs := []string{"exec", containerID, "bash", "-c", fmt.Sprintf("ls -t -d -1 %s", pattern)}
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

func getContainerEnvVar(containerID string, varName string) (string, error) {

	// could inspect the docker container as an alternative approach
	dockerArgs := []string{"exec", containerID, "bash", "-c", fmt.Sprintf("echo $%s", varName)}
	dockerCmd := exec.Command("docker", dockerArgs...)
	buf, err := dockerCmd.CombinedOutput()
	if err != nil {
		errMessage := string(buf)
		return "", fmt.Errorf("Docker exec error: %s (%s)", err, strings.TrimSpace(errMessage))
	}

	return string(buf), nil
}

func testContainerPathExists(containerID string, path string) (bool, error) {
	dockerArgs := []string{"exec", containerID, "bash", "-c", fmt.Sprintf("[[ -d %s ]]; echo $?", path)}
	dockerCmd := exec.Command("docker", dockerArgs...)
	buf, err := dockerCmd.CombinedOutput()
	if err != nil {
		errMessage := string(buf)
		return false, fmt.Errorf("Docker exec error: %s (%s)", err, strings.TrimSpace(errMessage))
	}

	response := strings.TrimSpace(string(buf))
	return response == "0", nil
}
