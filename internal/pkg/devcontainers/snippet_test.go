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

	ioutil.WriteFile(filepath.Join(root, "test1.sh"), []byte{}, 0755)
	ioutil.WriteFile(filepath.Join(root, "test2.sh"), []byte{}, 0755)

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

	ioutil.WriteFile(filepath.Join(root, "_ignore.sh"), []byte{}, 0755)
	ioutil.WriteFile(filepath.Join(root, ".ignore.sh"), []byte{}, 0755)
	ioutil.WriteFile(filepath.Join(root, "test1.sh"), []byte{}, 0755)

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

	ioutil.WriteFile(filepath.Join(root1, "test1.sh"), []byte{}, 0755)
	ioutil.WriteFile(filepath.Join(root2, "test1.sh"), []byte{}, 0755)

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
