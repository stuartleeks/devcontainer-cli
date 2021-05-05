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
	assert.NoError(t, err)
	defer os.RemoveAll(root)

	folders := []string{root}

	_ = ioutil.WriteFile(filepath.Join(root, "test1.sh"), []byte{}, 0755)
	_ = ioutil.WriteFile(filepath.Join(root, "test2.sh"), []byte{}, 0755)

	snippets, err := getSnippetsFromFolders(folders)
	assert.NoError(t, err)

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
	assert.NoError(t, err)
	defer os.RemoveAll(root)

	folders := []string{root}

	_ = ioutil.WriteFile(filepath.Join(root, "_ignore.sh"), []byte{}, 0755)
	_ = ioutil.WriteFile(filepath.Join(root, ".ignore.sh"), []byte{}, 0755)
	_ = ioutil.WriteFile(filepath.Join(root, "test1.sh"), []byte{}, 0755)

	snippets, err := getSnippetsFromFolders(folders)
	assert.NoError(t, err)

	expectedTemplates := []DevcontainerSnippet{
		{
			Name: "test1",
			Type: DevcontainerSnippetTypeSingleFile,
			Path: filepath.Join(root, "test1.sh"),
		},
	}

	assert.ElementsMatch(t, expectedTemplates, snippets)
}
func TestGetSnippets_TakesFilesInPriorityOrder(t *testing.T) {

	root1, err := ioutil.TempDir("", "devcontainer*")
	assert.NoError(t, err)
	defer os.RemoveAll(root1)
	root2, err := ioutil.TempDir("", "devcontainer*")
	assert.NoError(t, err)
	defer os.RemoveAll(root1)

	folders := []string{root1, root2}

	_ = ioutil.WriteFile(filepath.Join(root1, "test1.sh"), []byte{}, 0755)
	_ = ioutil.WriteFile(filepath.Join(root2, "test1.sh"), []byte{}, 0755)

	snippets, err := getSnippetsFromFolders(folders)
	assert.NoError(t, err)

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

	root, err := ioutil.TempDir("", "devcontainer*")
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

	// Add snippet
	snippet := DevcontainerSnippet{
		Name: "test",
		Path: snippetFilename,
		Type: DevcontainerSnippetTypeSingleFile,
	}
	err = addSingleFileSnippetToDevContainer(targetFolder, &snippet)
	assert.NoError(t, err)

	buf, err := ioutil.ReadFile(filepath.Join(devcontainerFolder, "scripts", "test1.sh"))
	assert.NoError(t, err)
	assert.Equal(t, "# dummy file", string(buf))

	buf, err = ioutil.ReadFile(filepath.Join(devcontainerFolder, "Dockerfile"))
	assert.NoError(t, err)
	assert.Equal(t, `FROM foo
RUN echo hi

# test
COPY scripts/test1.sh /tmp/
RUN /tmp/test1.sh
`, string(buf))
}
func TestSingleFileAddSnippet_WithInsertionPoint(t *testing.T) {

	root, err := ioutil.TempDir("", "devcontainer*")
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

	// Add snippet
	snippet := DevcontainerSnippet{
		Name: "test",
		Path: snippetFilename,
		Type: DevcontainerSnippetTypeSingleFile,
	}
	err = addSingleFileSnippetToDevContainer(targetFolder, &snippet)
	assert.NoError(t, err)

	buf, err := ioutil.ReadFile(filepath.Join(devcontainerFolder, "scripts", "test1.sh"))
	assert.NoError(t, err)
	assert.Equal(t, "# dummy file", string(buf))

	buf, err = ioutil.ReadFile(filepath.Join(devcontainerFolder, "Dockerfile"))
	assert.NoError(t, err)
	assert.Equal(t, `FROM foo
RUN echo hi
# __DEVCONTAINER_SNIPPET_INSERT__ 

# test
COPY scripts/test1.sh /tmp/
RUN /tmp/test1.sh

RUN echo hi2
`, string(buf))
}
