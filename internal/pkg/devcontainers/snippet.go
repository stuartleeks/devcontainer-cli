package devcontainers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/stuartleeks/devcontainer-cli/internal/pkg/config"
	"github.com/stuartleeks/devcontainer-cli/internal/pkg/errors"
)

type DevcontainerSnippetType int

const (
	DevcontainerSnippetTypeSingleFile = 1
	DevcontainerSnippetTypeFolder     = 2
)

// DevcontainerSnippet holds info on snippets for list/add etc
// Snippets can be either single script files or a directory with a set of files
type DevcontainerSnippet struct {
	Name string
	Type DevcontainerSnippetType
	// Path is the path to either the path to the single script file or to the directory for multi-file snippets
	Path string
}

// GetSnippetByName returns the template with the specified name or nil if not found
func GetSnippetByName(name string) (*DevcontainerSnippet, error) {
	// TODO - could possibly make this quicker by searching using the name rather than listing all and filtering
	snippets, err := GetSnippets()
	if err != nil {
		return nil, err
	}
	for _, snippet := range snippets {
		if snippet.Name == name {
			return &snippet, nil
		}
	}
	return nil, nil
}

// GetSnippets returns a list of discovered templates
func GetSnippets() ([]DevcontainerSnippet, error) {

	folders := config.GetSnippetFolders()
	if len(folders) == 0 {
		return []DevcontainerSnippet{}, &errors.StatusError{Message: "No snippet folders configured - see https://github.com/stuartleeks/devcontainer-cli/#working-with-devcontainer-snippets"}
	}

	snippets, err := getSnippetsFromFolders(folders)
	if err != nil {
		return []DevcontainerSnippet{}, err
	}
	return snippets, nil
}

func getSnippetsFromFolders(folders []string) ([]DevcontainerSnippet, error) {
	snippets := []DevcontainerSnippet{}
	snippetNames := map[string]bool{}
	for _, folder := range folders {
		folder := os.ExpandEnv(folder)
		newSnippets, err := getSnippetsFromFolder(folder)
		if err != nil {
			return []DevcontainerSnippet{}, err
		}
		for _, snippet := range newSnippets {
			if !snippetNames[snippet.Name] {
				snippetNames[snippet.Name] = true
				snippets = append(snippets, snippet)
			}
		}
	}
	sort.Slice(snippets, func(i int, j int) bool { return snippets[i].Name < snippets[j].Name })
	return snippets, nil
}

func getSnippetsFromFolder(folder string) ([]DevcontainerSnippet, error) {
	c, err := ioutil.ReadDir(folder)

	if err != nil {
		return []DevcontainerSnippet{}, fmt.Errorf("Error reading snippet definitions: %s\n", err)
	}

	snippets := []DevcontainerSnippet{}
	for _, entry := range c {
		if strings.HasPrefix(entry.Name(), ".") || strings.HasPrefix(entry.Name(), "_") {
			// ignore files/directories starting with "_" or "."
			continue
		}
		if entry.IsDir() {
			// TODO!
		} else {
			if strings.HasSuffix(entry.Name(), ".sh") {
				snippet := DevcontainerSnippet{
					Name: strings.TrimSuffix(entry.Name(), ".sh"),
					Type: DevcontainerSnippetTypeSingleFile,
					Path: filepath.Join(folder, entry.Name()),
				}
				snippets = append(snippets, snippet)
			}
		}
	}
	return snippets, nil
}
