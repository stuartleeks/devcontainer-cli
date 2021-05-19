package devcontainers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/stuartleeks/devcontainer-cli/internal/pkg/config"
	"github.com/stuartleeks/devcontainer-cli/internal/pkg/errors"
	ioutil2 "github.com/stuartleeks/devcontainer-cli/internal/pkg/ioutil"

	dora_ast "github.com/bradford-hamilton/dora/pkg/ast"
	dora_lexer "github.com/bradford-hamilton/dora/pkg/lexer"
	dora_merge "github.com/bradford-hamilton/dora/pkg/merge"
	dora_parser "github.com/bradford-hamilton/dora/pkg/parser"
)

type DevcontainerSnippetType string

const (
	DevcontainerSnippetTypeSingleFile = "Snippet:SingleFile"
	DevcontainerSnippetTypeFolder     = "Snippet:Folder"
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
			snippetJSONPath := filepath.Join(folder, entry.Name(), "snippet.json")
			snippetJSONInfo, err := os.Stat(snippetJSONPath)
			if err != nil || snippetJSONInfo.IsDir() {
				continue
			}
			snippet := DevcontainerSnippet{
				Name: entry.Name(),
				Type: DevcontainerSnippetTypeFolder,
				Path: filepath.Join(folder, entry.Name()),
			}
			snippets = append(snippets, snippet)
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

func AddSnippetToDevcontainer(projectFolder string, snippetName string) error {

	snippet, err := GetSnippetByName(snippetName)
	if err != nil {
		return err
	}
	if snippet == nil {
		return fmt.Errorf("Snippet '%s' not found\n", snippetName)
	}
	return addSnippetToDevcontainer(projectFolder, snippet)
}
func addSnippetToDevcontainer(projectFolder string, snippet *DevcontainerSnippet) error {
	switch snippet.Type {
	case DevcontainerSnippetTypeSingleFile:
		return addSingleFileSnippetToDevContainer(projectFolder, snippet)
	case DevcontainerSnippetTypeFolder:
		return addFolderSnippetToDevContainer(projectFolder, snippet)
	default:
		return fmt.Errorf("Unhandled snippet type: %q", snippet.Type)
	}
}

func addSingleFileSnippetToDevContainer(projectFolder string, snippet *DevcontainerSnippet) error {

	if snippet.Type != DevcontainerSnippetTypeSingleFile {
		return fmt.Errorf("Expected single file snippet")
	}

	scriptFolderPath := filepath.Join(projectFolder, ".devcontainer", "scripts")
	if err := os.MkdirAll(scriptFolderPath, 0755); err != nil {
		return err
	}
	_, scriptFilename := filepath.Split(snippet.Path)
	if err := ioutil2.CopyFile(snippet.Path, filepath.Join(scriptFolderPath, scriptFilename), 0755); err != nil {
		return err
	}

	dockerfileFilename := filepath.Join(projectFolder, ".devcontainer", "Dockerfile")
	buf, err := ioutil.ReadFile(dockerfileFilename)
	if err != nil {
		return fmt.Errorf("Error reading Dockerfile: %s", err)
	}

	snippetContent := fmt.Sprintf(`
# %[1]s
COPY scripts/%[2]s /tmp/
RUN /tmp/%[2]s
`, snippet.Name, scriptFilename)

	dockerfileContent := string(buf)
	dockerFileLines := strings.Split(dockerfileContent, "\n")
	addSeparator := false
	addedSnippetContent := false
	var newContent bytes.Buffer
	for _, line := range dockerFileLines {
		if addSeparator {
			if _, err = newContent.WriteString("\n"); err != nil {
				return err
			}
		}
		addSeparator = true
		if _, err = newContent.WriteString(line); err != nil {
			return err
		}

		if strings.Contains(line, "__DEVCONTAINER_SNIPPET_INSERT__") {
			if _, err = newContent.WriteString("\n"); err != nil {
				return err
			}
			if _, err = newContent.WriteString(snippetContent); err != nil {
				return err
			}
			addedSnippetContent = true
			addSeparator = false // avoid extra separator
		}
	}

	if !addedSnippetContent {
		if _, err = newContent.WriteString(snippetContent); err != nil {
			return err
		}
	}

	err = ioutil.WriteFile(dockerfileFilename, newContent.Bytes(), 0)

	return err
}

func addFolderSnippetToDevContainer(projectFolder string, snippet *DevcontainerSnippet) error {
	if snippet.Type != DevcontainerSnippetTypeFolder {
		return fmt.Errorf("Expected folder snippet")
	}

	snippetDevcontainerJSONPath := filepath.Join(snippet.Path, "devcontainer.json")
	snippetDevcontainerJSONInfo, err := os.Stat(snippetDevcontainerJSONPath)
	if err == nil && !snippetDevcontainerJSONInfo.IsDir() {
		projectDevcontainerFilename := filepath.Join(projectFolder, ".devcontainer", "devcontainer.json")
		baseDocument, err := loadJSONDocument(projectDevcontainerFilename)
		if err != nil {
			return err
		}

		snippetDocument, err := loadJSONDocument(snippetDevcontainerJSONPath)
		if err != nil {
			return err
		}

		resultDocument, err := dora_merge.MergeJSON(*baseDocument, *snippetDocument)
		if err != nil {
			return err
		}

		resultJSON, err := dora_ast.WriteJSONString(resultDocument)

		ioutil.WriteFile(projectDevcontainerFilename, []byte(resultJSON), 0666)
	}

	// TODO handle snippet.json actions

	return nil
}

func loadJSONDocument(path string) (*dora_ast.RootNode, error) {

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	l := dora_lexer.New(string(buf))
	p := dora_parser.New(l)
	baseDocument, err := p.ParseJSON()
	if err != nil {
		return nil, err
	}
	return &baseDocument, nil
}
