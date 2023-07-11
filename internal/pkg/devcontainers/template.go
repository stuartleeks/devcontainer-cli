package devcontainers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/stuartleeks/devcontainer-cli/internal/pkg/config"
	"github.com/stuartleeks/devcontainer-cli/internal/pkg/errors"
	ioutil2 "github.com/stuartleeks/devcontainer-cli/internal/pkg/ioutil"
)

// DevcontainerTemplate holds info on templates for list/add etc
type DevcontainerTemplate struct {
	Name string
	// Path is the path including the .devcontainer folder
	Path string
}

// GetTemplateByName returns the template with the specified name or nil if not found
func GetTemplateByName(name string) (*DevcontainerTemplate, error) {
	// TODO - could possibly make this quicker by searching using the name rather than listing all and filtering
	templates, err := GetTemplates()
	if err != nil {
		return nil, err
	}
	for _, template := range templates {
		if template.Name == name {
			return &template, nil
		}
	}
	return nil, nil
}

// GetTemplates returns a list of discovered templates
func GetTemplates() ([]DevcontainerTemplate, error) {

	folders := config.GetTemplateFolders()
	if len(folders) == 0 {
		return []DevcontainerTemplate{}, &errors.StatusError{Message: "No template folders configured - see https://github.com/stuartleeks/devcontainer-cli/#working-with-devcontainer-templates"}
	}
	templates, err := getTemplatesFromFolders(folders)
	if err != nil {
		return []DevcontainerTemplate{}, err
	}
	return templates, nil
}

func getTemplatesFromFolders(folders []string) ([]DevcontainerTemplate, error) {
	templates := []DevcontainerTemplate{}
	templateNames := map[string]bool{}

	for _, folder := range folders {
		folder := os.ExpandEnv(folder)
		newTemplates, err := getTemplatesFromFolder(folder)
		if err != nil {
			return []DevcontainerTemplate{}, err
		}
		for _, template := range newTemplates {
			if !templateNames[template.Name] {
				templateNames[template.Name] = true
				templates = append(templates, template)
			}
		}
	}
	sort.Slice(templates, func(i int, j int) bool { return templates[i].Name < templates[j].Name })
	return templates, nil
}

func getTemplatesFromFolder(folder string) ([]DevcontainerTemplate, error) {
	isDevcontainerFolder := func(parentPath string, fi os.FileInfo) bool {
		if !fi.IsDir() {
			return false
		}
		// TOODO - add support for templates with .devcontainer.json rather than .devcontainer folder
		devcontainerJsonPath := filepath.Join(parentPath, fi.Name(), ".devcontainer/devcontainer.json")
		devContainerJsonInfo, err := os.Stat(devcontainerJsonPath)
		return err == nil && !devContainerJsonInfo.IsDir()
	}
	c, err := ioutil.ReadDir(folder)

	if err != nil {
		return []DevcontainerTemplate{}, fmt.Errorf("Error reading devcontainer definitions: %s\n", err)
	}

	templates := []DevcontainerTemplate{}
	for _, entry := range c {
		if isDevcontainerFolder(folder, entry) {
			template := DevcontainerTemplate{
				Name: entry.Name(),
				Path: filepath.Join(folder, entry.Name(), ".devcontainer"),
			}
			templates = append(templates, template)
		}
	}
	return templates, nil
}

func GetDefaultDevcontainerNameForFolder(folderPath string) (string, error) {

	absPath, err := filepath.Abs(folderPath)
	if err != nil {
		return "", err
	}

	_, folderName := filepath.Split(absPath)
	return folderName, nil
}

func CopyTemplateToFolder(templatePath string, targetFolder string, devcontainerName string) error {
	var err error

	if err = ioutil2.CopyFolder(templatePath, filepath.Join(targetFolder, ".devcontainer")); err != nil {
		return fmt.Errorf("Error copying folder: %s\n", err)
	}

	// by default the "name" in devcontainer.json is set to the name of the template
	// override it here with the value passed in as --devcontainer-name (or the containing folder if not set)
	if devcontainerName == "" {
		devcontainerName, err = GetDefaultDevcontainerNameForFolder(targetFolder)
		if err != nil {
			return fmt.Errorf("Error getting default devcontainer name: %s", err)
		}
	}
	devcontainerJsonPath := filepath.Join(targetFolder, ".devcontainer", "devcontainer.json")
	err = SetDevcontainerName(devcontainerJsonPath, devcontainerName)
	if err != nil {
		return fmt.Errorf("Error setting devcontainer name: %s", err)
	}

	values, err := getSubstitutionValuesFromFile(devcontainerJsonPath)
	if err != nil {
		return fmt.Errorf("Error getting substituion values: %s", err)
	}
	err = recursiveSubstituteValues(values, filepath.Join(targetFolder, ".devcontainer"))
	if err != nil {
		return fmt.Errorf("Error performing substitution: %s", err)
	}

	return nil
}

func recursiveSubstituteValues(values *SubstitutionValues, path string) error {
	_, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("Error reading folder: %s\n", err)
	}

	subItems, err := ioutil.ReadDir(path)
	if err != nil {
		return fmt.Errorf("Error reading source folder contents: %s\n", err)
	}

	for _, subItem := range subItems {
		if subItem.IsDir() {
			if err = recursiveSubstituteValues(values, filepath.Join(path, subItem.Name())); err != nil {
				return err
			}
		} else {
			if err = performSubstitutionFile(values, filepath.Join(path, subItem.Name())); err != nil {
				return err
			}
		}
	}
	return nil
}

func SetDevcontainerName(devContainerJsonPath string, name string) error {
	// This doesn't use `json` as devcontainer.json permits comments (and the default templates include them!)

	// TODO - update this to use dora to query
	// TODO - update this to replace __DEVCONTAINER_USER_NAME__ and __DEVCONTAINER_HOME__

	buf, err := ioutil.ReadFile(devContainerJsonPath)
	if err != nil {
		return fmt.Errorf("error reading file %q: %s", devContainerJsonPath, err)
	}
	content := string(buf)

	// replace `name` property in JSON
	r := regexp.MustCompile("(\"name\"\\s*:\\s*\")[^\"]*(\")")
	replacement := "${1}" + name + "${2}"
	content = r.ReplaceAllString(content, replacement)

	// replace __DEVCONTAINER_NAME__ with name
	content = strings.ReplaceAll(content, "__DEVCONTAINER_NAME__", name)

	buf = []byte(content)
	if err = ioutil.WriteFile(devContainerJsonPath, buf, 0777); err != nil {
		return fmt.Errorf("error writing file %q: %s", devContainerJsonPath, err)
	}

	return nil
}

// "remoteUser": "vscode"
func GetDevContainerUserName(devContainerJsonPath string) (string, error) {
	buf, err := ioutil.ReadFile(devContainerJsonPath)
	if err != nil {
		return "", fmt.Errorf("error reading file %q: %s", devContainerJsonPath, err)
	}

	r := regexp.MustCompile("\n[^/]*\"remoteUser\"\\s*:\\s*\"([^\"]*)\"")
	match := r.FindStringSubmatch(string(buf))

	if len(match) <= 0 {
		return "", nil
	}
	return match[1], nil
}
