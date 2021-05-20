package devcontainers

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSnippets_ListsSingleFileTemplates(t *testing.T) {

	root, err := ioutil.TempDir("", "devcontainer*")
	if !assert.NoError(t, err) {
		return
	}
	defer os.RemoveAll(root)

	folders := []string{root}

	_ = ioutil.WriteFile(filepath.Join(root, "test1.sh"), []byte{}, 0755)
	_ = ioutil.WriteFile(filepath.Join(root, "test2.sh"), []byte{}, 0755)

	snippets, err := getSnippetsFromFolders(folders)
	if !assert.NoError(t, err) {
		return
	}

	expectedTemplates := []DevcontainerSnippet{
		{
			Name: "test1",
			Type: DevcontainerSnippetTypeSingleFile,
			Path: filepath.Join(root, "test1.sh"),
		},
		{
			Name: "test2",
			Type: DevcontainerSnippetTypeSingleFile,
			Path: filepath.Join(root, "test2.sh"),
		},
	}

	assert.ElementsMatch(t, expectedTemplates, snippets)
}
func TestGetSnippets_IgnoresFilesWithIncorrectPrefix(t *testing.T) {

	root, err := ioutil.TempDir("", "devcontainer*")
	if !assert.NoError(t, err) {
		return
	}
	defer os.RemoveAll(root)

	folders := []string{root}

	_ = ioutil.WriteFile(filepath.Join(root, "_ignore.sh"), []byte{}, 0755)
	_ = ioutil.WriteFile(filepath.Join(root, ".ignore.sh"), []byte{}, 0755)
	_ = ioutil.WriteFile(filepath.Join(root, "test1.sh"), []byte{}, 0755)

	snippets, err := getSnippetsFromFolders(folders)
	if !assert.NoError(t, err) {
		return
	}

	expectedTemplates := []DevcontainerSnippet{
		{
			Name: "test1",
			Type: DevcontainerSnippetTypeSingleFile,
			Path: filepath.Join(root, "test1.sh"),
		},
	}

	assert.ElementsMatch(t, expectedTemplates, snippets)
}
func TestGetSnippets_ListsFolderTemplate(t *testing.T) {

	root1, err := ioutil.TempDir("", "devcontainer*")
	if !assert.NoError(t, err) {
		return
	}
	defer os.RemoveAll(root1)
	root2, err := ioutil.TempDir("", "devcontainer*")
	if !assert.NoError(t, err) {
		return
	}
	defer os.RemoveAll(root1)

	folders := []string{root1, root2}

	_ = os.MkdirAll(filepath.Join(root1, "test1"), 0755)
	_ = ioutil.WriteFile(filepath.Join(root1, "test1/snippet.json"), []byte{}, 0755)
	_ = ioutil.WriteFile(filepath.Join(root2, "test1.sh"), []byte{}, 0755)

	snippets, err := getSnippetsFromFolders(folders)
	if !assert.NoError(t, err) {
		return
	}

	expectedTemplates := []DevcontainerSnippet{
		{
			Name: "test1",
			Type: DevcontainerSnippetTypeFolder,
			Path: filepath.Join(root1, "test1"), // Uses root1 as it is in the list first
		},
	}

	assert.ElementsMatch(t, expectedTemplates, snippets)
}
func TestGetSnippets_TakesFilesInPriorityOrder(t *testing.T) {

	root1, err := ioutil.TempDir("", "devcontainer*")
	if !assert.NoError(t, err) {
		return
	}
	defer os.RemoveAll(root1)
	root2, err := ioutil.TempDir("", "devcontainer*")
	if !assert.NoError(t, err) {
		return
	}
	defer os.RemoveAll(root1)

	folders := []string{root1, root2}

	_ = ioutil.WriteFile(filepath.Join(root1, "test1.sh"), []byte{}, 0755)
	_ = ioutil.WriteFile(filepath.Join(root2, "test1.sh"), []byte{}, 0755)

	snippets, err := getSnippetsFromFolders(folders)
	if !assert.NoError(t, err) {
		return
	}

	expectedTemplates := []DevcontainerSnippet{
		{
			Name: "test1",
			Type: DevcontainerSnippetTypeSingleFile,
			Path: filepath.Join(root1, "test1.sh"), // Uses root1 as it is in the list first
		},
	}

	assert.ElementsMatch(t, expectedTemplates, snippets)
}

func TestSingleFileAddSnippet_NoInsertionPoint(t *testing.T) {

	root, _ := ioutil.TempDir("", "devcontainer*")
	defer os.RemoveAll(root)

	// set up snippet
	snippetFolder := filepath.Join(root, "snippets")
	_ = os.MkdirAll(snippetFolder, 0755)
	snippetFilename := filepath.Join(snippetFolder, "test1.sh")
	_ = ioutil.WriteFile(snippetFilename, []byte("# dummy file"), 0755)

	// set up devcontainer
	targetFolder := filepath.Join(root, "target")
	devcontainerFolder := filepath.Join(targetFolder, ".devcontainer")
	_ = os.MkdirAll(devcontainerFolder, 0755)

	_ = ioutil.WriteFile(filepath.Join(devcontainerFolder, "Dockerfile"), []byte(`FROM foo
RUN echo hi
`), 0755)
	_ = ioutil.WriteFile(filepath.Join(devcontainerFolder, "devcontainer.json"), []byte(`{
		"name" : "testname"
	}`), 0755)

	// Add snippet
	snippet := DevcontainerSnippet{
		Name: "test",
		Path: snippetFilename,
		Type: DevcontainerSnippetTypeSingleFile,
	}
	err := addSingleFileSnippetToDevContainer(targetFolder, &snippet)
	if !assert.NoError(t, err) {
		return
	}

	buf, err := ioutil.ReadFile(filepath.Join(devcontainerFolder, "scripts", "test1.sh"))
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, "# dummy file", string(buf))

	buf, err = ioutil.ReadFile(filepath.Join(devcontainerFolder, "Dockerfile"))
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, `FROM foo
RUN echo hi

# test
COPY scripts/test1.sh /tmp/
RUN /tmp/test1.sh
`, string(buf))
}
func TestSingleFileAddSnippet_WithInsertionPoint(t *testing.T) {

	root, _ := ioutil.TempDir("", "devcontainer*")
	defer os.RemoveAll(root)

	// set up snippet
	snippetFolder := filepath.Join(root, "snippets")
	_ = os.MkdirAll(snippetFolder, 0755)
	snippetFilename := filepath.Join(snippetFolder, "test1.sh")
	_ = ioutil.WriteFile(snippetFilename, []byte("# dummy file"), 0755)

	// set up devcontainer
	targetFolder := filepath.Join(root, "target")
	devcontainerFolder := filepath.Join(targetFolder, ".devcontainer")
	_ = os.MkdirAll(devcontainerFolder, 0755)

	_ = ioutil.WriteFile(filepath.Join(devcontainerFolder, "Dockerfile"), []byte(`FROM foo
RUN echo hi
# __DEVCONTAINER_SNIPPET_INSERT__ 

RUN echo hi2
`), 0755)
	_ = ioutil.WriteFile(filepath.Join(devcontainerFolder, "devcontainer.json"), []byte(`{
	"name" : "testname"
}`), 0755)

	// Add snippet
	snippet := DevcontainerSnippet{
		Name: "test",
		Path: snippetFilename,
		Type: DevcontainerSnippetTypeSingleFile,
	}
	err := addSingleFileSnippetToDevContainer(targetFolder, &snippet)
	if !assert.NoError(t, err) {
		return
	}

	buf, err := ioutil.ReadFile(filepath.Join(devcontainerFolder, "scripts", "test1.sh"))
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, "# dummy file", string(buf))

	buf, err = ioutil.ReadFile(filepath.Join(devcontainerFolder, "Dockerfile"))
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, `FROM foo
RUN echo hi
# test
COPY scripts/test1.sh /tmp/
RUN /tmp/test1.sh

# __DEVCONTAINER_SNIPPET_INSERT__ 

RUN echo hi2
`, string(buf))
}

func TestFolderAddSnippet_MergesDevcontainerJSON(t *testing.T) {

	root, _ := ioutil.TempDir("", "devcontainer*")
	defer os.RemoveAll(root)

	// set up snippet
	snippetFolder := filepath.Join(root, "snippets/test1")
	_ = os.MkdirAll(snippetFolder, 0755)
	snippetJSONFilename := filepath.Join(snippetFolder, "snippet.json")
	_ = ioutil.WriteFile(snippetJSONFilename, []byte(`{
		"actions": [
			{
				"type": "mergeJSON",
				"source": "devcontainer.json",
				"target": ".devcontainer/devcontainer.json"
			}
		]
	}`), 0755)

	snippetDevcontainerFilename := filepath.Join(snippetFolder, "devcontainer.json")
	_ = ioutil.WriteFile(snippetDevcontainerFilename, []byte(`// For format details, see https://aka.ms/vscode-remote/devcontainer.json or this file's README at:
// https://github.com/microsoft/vscode-dev-containers/tree/v0.117.1/containers/go
{
	"runArgs": [
		// Mount go mod cache
		"-v", "devcontainer-cli-gomodcache:/go/pkg",
	],

	// Set *default* container specific settings.json values on container create.
	"settings": {
		"go.gopath": "/go",
		"go.useLanguageServer": true,
		"[go]": {
			"editor.snippetSuggestions": "none",
			"editor.formatOnSave": true,
			"editor.codeActionsOnSave": {
				"source.organizeImports": true,
			}
		},
		"gopls": {
			"usePlaceholders": true, // add parameter placeholders when completing a function
			// Experimental settings
			"completeUnimported": true, // autocomplete unimported packages
			"watchFileChanges": true, // watch file changes outside of the editor
			"deepCompletion": true, // enable deep completion
		},
	},
	
	// Add the IDs of extensions you want installed when the container is created.
	"extensions": [
		"stuartleeks.vscode-go-by-example",
		"golang.go",
	],
}`), 0755)

	// set up devcontainer
	targetFolder := filepath.Join(root, "target")
	devcontainerFolder := filepath.Join(targetFolder, ".devcontainer")
	_ = os.MkdirAll(devcontainerFolder, 0755)

	_ = ioutil.WriteFile(filepath.Join(devcontainerFolder, "devcontainer.json"), []byte(`// For format details, see https://aka.ms/vscode-remote/devcontainer.json or this file's README at:
// https://github.com/microsoft/vscode-dev-containers/tree/v0.117.1/containers/go
{
	"name": "test",
	"dockerFile": "Dockerfile",
	"runArgs": [
		// Use host network
		"--network=host",
	],

	// Set *default* container specific settings.json values on container create.
	"settings": {
		"terminal.integrated.shell.linux": "/bin/bash",
	},
	
	// Add the IDs of extensions you want installed when the container is created.
	"extensions": [
		"example.test",
	],

	// Use 'forwardPorts' to make a list of ports inside the container available locally.
	// "forwardPorts": [],

	// Uncomment to connect as a non-root user. See https://aka.ms/vscode-remote/containers/non-root.
	"remoteUser": "vscode"
}`), 0755)

	// Add snippet
	snippet := DevcontainerSnippet{
		Name: "test",
		Path: snippetFolder,
		Type: DevcontainerSnippetTypeFolder,
	}
	err := addSnippetToDevcontainer(targetFolder, &snippet)
	if !assert.NoError(t, err) {
		return
	}

	buf, err := ioutil.ReadFile(filepath.Join(devcontainerFolder, "devcontainer.json"))
	if !assert.NoError(t, err) {
		return
	}
	stringContent := string(buf)
	assert.Equal(t, `// For format details, see https://aka.ms/vscode-remote/devcontainer.json or this file's README at:
// https://github.com/microsoft/vscode-dev-containers/tree/v0.117.1/containers/go
{
	"name": "test",
	"dockerFile": "Dockerfile",
	"runArgs": [
		// Use host network
		"--network=host",
		// Mount go mod cache
		"-v", "devcontainer-cli-gomodcache:/go/pkg",
	],

	// Set *default* container specific settings.json values on container create.
	"settings": {
		"terminal.integrated.shell.linux": "/bin/bash",
		"go.gopath": "/go",
		"go.useLanguageServer": true,
		"[go]": {
			"editor.snippetSuggestions": "none",
			"editor.formatOnSave": true,
			"editor.codeActionsOnSave": {
				"source.organizeImports": true,
			}
		},
		"gopls": {
			"usePlaceholders": true, // add parameter placeholders when completing a function
			// Experimental settings
			"completeUnimported": true, // autocomplete unimported packages
			"watchFileChanges": true, // watch file changes outside of the editor
			"deepCompletion": true, // enable deep completion
		},
	},
	
	// Add the IDs of extensions you want installed when the container is created.
	"extensions": [
		"example.test",
		"stuartleeks.vscode-go-by-example",
		"golang.go",
	],

	// Use 'forwardPorts' to make a list of ports inside the container available locally.
	// "forwardPorts": [],

	// Uncomment to connect as a non-root user. See https://aka.ms/vscode-remote/containers/non-root.
	"remoteUser": "vscode"
}`, stringContent)
}

func TestFolderAddSnippet_CopiesScriptAndUpdatesDockerfile(t *testing.T) {

	root, _ := ioutil.TempDir("", "devcontainer*")
	defer os.RemoveAll(root)

	// set up snippet
	snippetFolder := filepath.Join(root, "snippets/test1")
	_ = os.MkdirAll(snippetFolder, 0755)
	snippetJSONFilename := filepath.Join(snippetFolder, "snippet.json")
	_ = ioutil.WriteFile(snippetJSONFilename, []byte(`{
		"actions": [
			{
				"type": "copyAndRun",
				"source": "script.sh"
			}
		]
	}`), 0755)

	scriptFilename := filepath.Join(snippetFolder, "script.sh")
	_ = ioutil.WriteFile(scriptFilename, []byte("# dummy file"), 0755)

	// set up devcontainer
	targetFolder := filepath.Join(root, "target")
	devcontainerFolder := filepath.Join(targetFolder, ".devcontainer")
	_ = os.MkdirAll(devcontainerFolder, 0755)

	_ = ioutil.WriteFile(filepath.Join(devcontainerFolder, "Dockerfile"), []byte(`FROM foo
RUN echo hi

# __DEVCONTAINER_SNIPPET_INSERT__ 

RUN echo hi2
`), 0755)
	_ = ioutil.WriteFile(filepath.Join(devcontainerFolder, "devcontainer.json"), []byte(`{
	"name" : "testname"
}`), 0755)

	// Add snippet
	snippet := DevcontainerSnippet{
		Name: "test",
		Path: snippetFolder,
		Type: DevcontainerSnippetTypeFolder,
	}
	err := addSnippetToDevcontainer(targetFolder, &snippet)
	if !assert.NoError(t, err) {
		return
	}

	buf, err := ioutil.ReadFile(filepath.Join(devcontainerFolder, "scripts", "script.sh"))
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, "# dummy file", string(buf))

	buf, err = ioutil.ReadFile(filepath.Join(devcontainerFolder, "Dockerfile"))
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, `FROM foo
RUN echo hi

# test
COPY scripts/script.sh /tmp/
RUN /tmp/script.sh

# __DEVCONTAINER_SNIPPET_INSERT__ 

RUN echo hi2
`, string(buf))
}

func TestFolderAddSnippet_InsertsSnippetsInDockerfile(t *testing.T) {

	root, _ := ioutil.TempDir("", "devcontainer*")
	defer os.RemoveAll(root)

	// set up snippet
	snippetFolder := filepath.Join(root, "snippets/test1")
	_ = os.MkdirAll(snippetFolder, 0755)
	snippetJSONFilename := filepath.Join(snippetFolder, "snippet.json")
	_ = ioutil.WriteFile(snippetJSONFilename, []byte(`{
		"actions": [
			{
				"type": "dockerfileSnippet",
				"content": "ENV FOO=BAR"
			},
			{
				"type": "dockerfileSnippet",
				"content": "# testing\nENV WIBBLE=BIBBLE"
			}
		]
	}`), 0755)

	// set up devcontainer
	targetFolder := filepath.Join(root, "target")
	devcontainerFolder := filepath.Join(targetFolder, ".devcontainer")
	_ = os.MkdirAll(devcontainerFolder, 0755)

	_ = ioutil.WriteFile(filepath.Join(devcontainerFolder, "Dockerfile"), []byte(`FROM foo
RUN echo hi

# __DEVCONTAINER_SNIPPET_INSERT__ 

RUN echo hi2
`), 0755)
	_ = ioutil.WriteFile(filepath.Join(devcontainerFolder, "devcontainer.json"), []byte(`{
	"name" : "testname"
}`), 0755)

	// Add snippet
	snippet := DevcontainerSnippet{
		Name: "test",
		Path: snippetFolder,
		Type: DevcontainerSnippetTypeFolder,
	}
	err := addSnippetToDevcontainer(targetFolder, &snippet)
	if !assert.NoError(t, err) {
		return
	}

	buf, err := ioutil.ReadFile(filepath.Join(devcontainerFolder, "Dockerfile"))
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, `FROM foo
RUN echo hi

ENV FOO=BAR

# testing
ENV WIBBLE=BIBBLE

# __DEVCONTAINER_SNIPPET_INSERT__ 

RUN echo hi2
`, string(buf))
}

func TestFolderAddSnippet_PerformsSubstitutionWithoutUserName(t *testing.T) {

	root, _ := ioutil.TempDir("", "devcontainer*")
	defer os.RemoveAll(root)

	// set up snippet
	snippetFolder := filepath.Join(root, "snippets/test1")
	_ = os.MkdirAll(snippetFolder, 0755)
	snippetJSONFilename := filepath.Join(snippetFolder, "snippet.json")
	_ = ioutil.WriteFile(snippetJSONFilename, []byte(`{
		"actions": [
			{
				"type": "dockerfileSnippet",
				"content": "ENV DC_NAME=__DEVCONTAINER_NAME__\nENV DC_USER_NAME=__DEVCONTAINER_USER_NAME__\nENV DC_HOME=__DEVCONTAINER_HOME__"
			}
		]
	}`), 0755)

	// set up devcontainer
	targetFolder := filepath.Join(root, "target")
	devcontainerFolder := filepath.Join(targetFolder, ".devcontainer")
	_ = os.MkdirAll(devcontainerFolder, 0755)

	_ = ioutil.WriteFile(filepath.Join(devcontainerFolder, "Dockerfile"), []byte(`FROM foo
RUN echo hi

# __DEVCONTAINER_SNIPPET_INSERT__ 

RUN echo hi2
`), 0755)
	_ = ioutil.WriteFile(filepath.Join(devcontainerFolder, "devcontainer.json"), []byte(`{
	"name" : "testname"
}`), 0755)

	// Add snippet
	snippet := DevcontainerSnippet{
		Name: "test",
		Path: snippetFolder,
		Type: DevcontainerSnippetTypeFolder,
	}
	err := addSnippetToDevcontainer(targetFolder, &snippet)
	if !assert.NoError(t, err) {
		return
	}

	buf, err := ioutil.ReadFile(filepath.Join(devcontainerFolder, "Dockerfile"))
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, `FROM foo
RUN echo hi

ENV DC_NAME=testname
ENV DC_USER_NAME=root
ENV DC_HOME=/root

# __DEVCONTAINER_SNIPPET_INSERT__ 

RUN echo hi2
`, string(buf))
}
func TestFolderAddSnippet_PerformsSubstitutionWithUserName(t *testing.T) {

	root, _ := ioutil.TempDir("", "devcontainer*")
	defer os.RemoveAll(root)

	// set up snippet
	snippetFolder := filepath.Join(root, "snippets/test1")
	_ = os.MkdirAll(snippetFolder, 0755)
	snippetJSONFilename := filepath.Join(snippetFolder, "snippet.json")
	_ = ioutil.WriteFile(snippetJSONFilename, []byte(`{
		"actions": [
			{
				"type": "dockerfileSnippet",
				"content": "ENV DC_NAME=__DEVCONTAINER_NAME__\nENV DC_USER_NAME=__DEVCONTAINER_USER_NAME__\nENV DC_HOME=__DEVCONTAINER_HOME__"
			}
		]
	}`), 0755)

	// set up devcontainer
	targetFolder := filepath.Join(root, "target")
	devcontainerFolder := filepath.Join(targetFolder, ".devcontainer")
	_ = os.MkdirAll(devcontainerFolder, 0755)

	_ = ioutil.WriteFile(filepath.Join(devcontainerFolder, "Dockerfile"), []byte(`FROM foo
RUN echo hi

# __DEVCONTAINER_SNIPPET_INSERT__ 

RUN echo hi2
`), 0755)
	_ = ioutil.WriteFile(filepath.Join(devcontainerFolder, "devcontainer.json"), []byte(`{
	"name" : "testname",
	"remoteUser": "dcuser"
}`), 0755)

	// Add snippet
	snippet := DevcontainerSnippet{
		Name: "test",
		Path: snippetFolder,
		Type: DevcontainerSnippetTypeFolder,
	}
	err := addSnippetToDevcontainer(targetFolder, &snippet)
	if !assert.NoError(t, err) {
		return
	}

	buf, err := ioutil.ReadFile(filepath.Join(devcontainerFolder, "Dockerfile"))
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, `FROM foo
RUN echo hi

ENV DC_NAME=testname
ENV DC_USER_NAME=dcuser
ENV DC_HOME=/home/dcuser

# __DEVCONTAINER_SNIPPET_INSERT__ 

RUN echo hi2
`, string(buf))
}
