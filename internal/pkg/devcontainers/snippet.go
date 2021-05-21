package devcontainers

import (
	"encoding/json"
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
	"github.com/bradford-hamilton/dora/pkg/dora"
	dora_lexer "github.com/bradford-hamilton/dora/pkg/lexer"
	dora_merge "github.com/bradford-hamilton/dora/pkg/merge"
	dora_parser "github.com/bradford-hamilton/dora/pkg/parser"
)

type SubstitutionValues struct {
	Name       string
	UserName   string
	HomeFolder string
}

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

type FolderSnippetActionType string

const (
	FolderSnippetActionMergeJSON         FolderSnippetActionType = "mergeJSON"         // merge JSON file from snippet with target JSON file
	FolderSnippetActionCopyAndRun        FolderSnippetActionType = "copyAndRun"        // COPY and RUN script from snippet in the Dockerfile (as with single-file snippet)
	FolderSnippetActionDockerfileSnippet FolderSnippetActionType = "dockerfileSnippet" // snippet to include as-is in the Dockerfile
)

type FolderSnippetAction struct {
	Type        FolderSnippetActionType `json:"type"`
	SourcePath  string                  `json:"source"`      // for mergeJSON this is snippet-relative path to JSON. for copyAndRun this is the script filename
	TargetPath  string                  `json:"target"`      // for mergeJSON this is project-relative path to JSON
	Content     string                  `json:"content"`     // for dockerfileSnippet this is the content to include
	ContentPath string                  `json:"contentPath"` // for dockerfileSnippet this is the path to content to include
}

// FolderSnippet maps to the content of the snippet.json file for folder-based snippets
type FolderSnippet struct {
	Actions []FolderSnippetAction `json:"actions"`
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
	snippetBasePath, scriptFilename := filepath.Split(snippet.Path)

	scriptFolderPath := filepath.Join(projectFolder, ".devcontainer", "scripts")
	err := copyAndRunScriptFile(projectFolder, snippet, snippetBasePath, scriptFolderPath, scriptFilename)
	return err
}

func addFolderSnippetToDevContainer(projectFolder string, snippet *DevcontainerSnippet) error {
	if snippet.Type != DevcontainerSnippetTypeFolder {
		return fmt.Errorf("Expected folder snippet")
	}

	snippetJSONPath := filepath.Join(snippet.Path, "snippet.json")
	buf, err := ioutil.ReadFile(snippetJSONPath)
	if err != nil {
		return err
	}
	var snippetJSON FolderSnippet
	err = json.Unmarshal(buf, &snippetJSON)
	if err != nil {
		return err
	}

	for _, action := range snippetJSON.Actions {
		switch action.Type {
		case FolderSnippetActionMergeJSON:
			if action.SourcePath == "" {
				return fmt.Errorf("source must be set for %s actions", action.Type)
			}
			if action.TargetPath == "" {
				return fmt.Errorf("target must be set for %s actions", action.Type)
			}
			err = mergeJSON(projectFolder, snippet, action.SourcePath, action.TargetPath)
			if err != nil {
				return err
			}
		case FolderSnippetActionCopyAndRun:
			if action.SourcePath == "" {
				return fmt.Errorf("source must be set for %s actions", action.Type)
			}
			targetPath := filepath.Join(projectFolder, ".devcontainer", "scripts")
			sourceParent, sourceFileName := filepath.Split(action.SourcePath)
			sourceBasePath := filepath.Join(snippet.Path, sourceParent)
			err = copyAndRunScriptFile(projectFolder, snippet, sourceBasePath, targetPath, sourceFileName)
			if err != nil {
				return err
			}
		case FolderSnippetActionDockerfileSnippet:
			var content string
			if action.Content != "" {
				if action.ContentPath != "" {
					return fmt.Errorf("can only set one of content and contentPath")
				}
				content = action.Content + "\n"
			} else if action.ContentPath != "" {
				buf, err = ioutil.ReadFile(filepath.Join(snippet.Path, action.ContentPath))
				if err != nil {
					return err
				}
				content = string(buf)
			} else {
				return fmt.Errorf("one of content and contentPath must be set for %s actions", action.Type)
			}
			dockerfileFilename := filepath.Join(projectFolder, ".devcontainer", "Dockerfile")
			err = insertDockerfileSnippet(projectFolder, dockerfileFilename, content)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unhandled action type: %q", action.Type)
		}
	}

	return nil
}

func copyAndRunScriptFile(projectFolder string, snippet *DevcontainerSnippet, snippetBasePath string, targetPath, scriptFilename string) error {
	if err := os.MkdirAll(targetPath, 0755); err != nil {
		return err
	}
	if err := ioutil2.CopyFile(filepath.Join(snippetBasePath, scriptFilename), filepath.Join(targetPath, scriptFilename), 0755); err != nil {
		return err
	}

	snippetContent := fmt.Sprintf(`# %[1]s
COPY scripts/%[2]s /tmp/
RUN /tmp/%[2]s
`, snippet.Name, scriptFilename)
	dockerfileFilename := filepath.Join(projectFolder, ".devcontainer", "Dockerfile")

	err := insertDockerfileSnippet(projectFolder, dockerfileFilename, snippetContent)
	return err
}

func insertDockerfileSnippet(projectFolder string, dockerfileFilename string, snippetContent string) error {

	buf, err := ioutil.ReadFile(dockerfileFilename)
	if err != nil {
		return fmt.Errorf("Error reading Dockerfile: %s", err)
	}

	dockerfileContent := string(buf)
	dockerFileLines := strings.Split(dockerfileContent, "\n")
	addSeparator := false
	addedSnippetContent := false
	var newContent strings.Builder
	for _, line := range dockerFileLines {
		if addSeparator {
			if _, err = newContent.WriteString("\n"); err != nil {
				return err
			}
		}
		addSeparator = true

		if strings.Contains(line, "__DEVCONTAINER_SNIPPET_INSERT__") {
			if _, err = newContent.WriteString(snippetContent); err != nil {
				return err
			}
			if _, err = newContent.WriteString("\n"); err != nil {
				return err
			}
			line += "\n"
			addedSnippetContent = true
			addSeparator = false // avoid extra separator
		}

		if _, err = newContent.WriteString(line); err != nil {
			return err
		}
	}

	if !addedSnippetContent {
		if _, err = newContent.WriteString("\n"); err != nil {
			return err
		}
		if _, err = newContent.WriteString(snippetContent); err != nil {
			return err
		}
	}

	content := newContent.String()
	values, err := getSubstitutionValuesFromFile(filepath.Join(projectFolder, ".devcontainer/devcontainer.json"))
	if err != nil {
		return fmt.Errorf("failed to get dev container values: %s", err)
	}
	content = performSubstitutionString(values, content)

	err = ioutil.WriteFile(dockerfileFilename, []byte(content), 0)

	return err

}
func mergeJSON(projectFolder string, snippet *DevcontainerSnippet, relativeMergePath string, relativeBasePath string) error {
	mergePath := filepath.Join(snippet.Path, relativeMergePath)
	_, err := os.Stat(mergePath)
	if err != nil {
		return err
	}
	basePath := filepath.Join(projectFolder, relativeBasePath)
	baseDocument, err := loadJSONDocument(basePath)
	if err != nil {
		return err
	}

	mergeDocument, err := loadJSONDocument(mergePath)
	if err != nil {
		return err
	}

	resultDocument, err := dora_merge.MergeJSON(*baseDocument, *mergeDocument)
	if err != nil {
		return err
	}

	resultJSON, err := dora_ast.WriteJSONString(resultDocument)
	if err != nil {
		return err
	}

	values, err := getSubstitutionValuesFromFile(filepath.Join(projectFolder, ".devcontainer/devcontainer.json"))
	if err != nil {
		return fmt.Errorf("failed to get dev container values: %s", err)
	}
	resultJSON = performSubstitutionString(values, resultJSON)

	err = ioutil.WriteFile(basePath, []byte(resultJSON), 0666)
	if err != nil {
		return fmt.Errorf("failed to write file: %s", err)
	}

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

func getSubstitutionValuesFromFile(devContainerJsonPath string) (*SubstitutionValues, error) {
	// This doesn't use standard `json` pkg as devcontainer.json permits comments (and the default templates include them!)

	buf, err := ioutil.ReadFile(devContainerJsonPath)
	if err != nil {
		return nil, err
	}

	c, err := dora.NewFromBytes(buf)
	if err != nil {
		return nil, err
	}

	name, err := c.GetString("$.name")
	if err != nil {
		name = ""
	}
	userName, err := c.GetString("$.remoteUser")
	if err != nil {
		userName = "root"
	}
	homeFolder := "/home/" + userName
	if userName == "root" {
		homeFolder = "/root"
	}

	return &SubstitutionValues{
		Name:       name,
		UserName:   userName,
		HomeFolder: homeFolder,
	}, nil
}

func performSubstitutionFile(substitutionValues *SubstitutionValues, filename string) error {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	content := string(buf)
	content = performSubstitutionString(substitutionValues, content)
	err = ioutil.WriteFile(filename, []byte(content), 0)
	return err
}

func performSubstitutionString(substitutionValues *SubstitutionValues, content string) string {
	// replace __DEVCONTAINER_NAME__ with name etc
	content = strings.ReplaceAll(content, "__DEVCONTAINER_NAME__", substitutionValues.Name)
	content = strings.ReplaceAll(content, "__DEVCONTAINER_USER_NAME__", substitutionValues.UserName)
	content = strings.ReplaceAll(content, "__DEVCONTAINER_HOME__", substitutionValues.HomeFolder)
	return content
}
