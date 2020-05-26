package devcontainers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/stuartleeks/devcontainer-cli/internal/pkg/config"
	"github.com/stuartleeks/devcontainer-cli/internal/pkg/errors"
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
	templates := []DevcontainerTemplate{}
	folders := config.GetTemplateFolders()
	if len(folders) == 0 {
		return []DevcontainerTemplate{}, &errors.StatusError{Message: "No template folders configured - see https://github.com/stuartleeks/devcontainer-cli/#working-with-devcontainer-templates"}
	}
	for _, folder := range folders {
		folder := os.ExpandEnv(folder)
		newTemplates, err := getTemplatesFromFolder(folder)
		if err != nil {
			return []DevcontainerTemplate{}, err
		}
		templates = append(templates, newTemplates...)
	}
	return templates, nil
}

func getTemplatesFromFolder(folder string) ([]DevcontainerTemplate, error) {
	isDevcontainerFolder := func(parentPath string, fi os.FileInfo) bool {
		if !fi.IsDir() {
			return false
		}
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
