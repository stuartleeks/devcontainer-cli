package devcontainers

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetDevcontainerName(t *testing.T) {

	f, err := ioutil.TempFile("", "test.json")
	assert.NoError(t, err)
	defer os.Remove(f.Name())

	_, _ = f.WriteString(`{
	"name": "initial",
	// here's a comment!
	"otherProperties": [
		"something",
		"here"
	]
}`)

	err = SetDevcontainerName(f.Name(), "newName")
	assert.NoError(t, err)

	buf, err := ioutil.ReadFile(f.Name())
	assert.NoError(t, err)

	assert.Equal(t, `{
	"name": "newName",
	// here's a comment!
	"otherProperties": [
		"something",
		"here"
	]
}`, string(buf))
}

func TestGetDevContainerUserName_Uncommented(t *testing.T) {

	f, err := ioutil.TempFile("", "test.json")
	assert.NoError(t, err)
	defer os.Remove(f.Name())

	_, _ = f.WriteString(`{
	"name": "initial",
	// here's a comment!
	"otherProperties": [
		"something",
		"here"
	],
	"remoteUser": "vscode"
}`)

	user, err := GetDevContainerUserName(f.Name())
	assert.NoError(t, err)
	assert.Equal(t, "vscode", user)
}

func TestGetDevContainerUserName_NotSet(t *testing.T) {

	f, err := ioutil.TempFile("", "test.json")
	assert.NoError(t, err)
	defer os.Remove(f.Name())

	_, _ = f.WriteString(`{
	"name": "initial",
	// here's a comment!
	"otherProperties": [
		"something",
		"here"
	]
}`)

	user, err := GetDevContainerUserName(f.Name())
	assert.NoError(t, err)
	assert.Equal(t, "", user)
}

func TestGetDevContainerUserName_Commented(t *testing.T) {

	f, err := ioutil.TempFile("", "test.json")
	assert.NoError(t, err)
	defer os.Remove(f.Name())

	_, _ = f.WriteString(`{
	"name": "initial",
	// here's a comment!
	"otherProperties": [
		"something",
		"here"
	],
	// "remoteUser": "vscode"
}`)

	user, err := GetDevContainerUserName(f.Name())
	assert.NoError(t, err)
	assert.Equal(t, "", user)
}

func TestGetTemplateFolders_ListsFoldersWithDevcontainers(t *testing.T) {

	root, err := ioutil.TempDir("", "devcontainer*")
	assert.NoError(t, err)
	defer os.RemoveAll(root)

	folders := []string{root}

	os.MkdirAll(filepath.Join(root, "test1", ".devcontainer"), 0755)
	ioutil.WriteFile(filepath.Join(root, "test1", ".devcontainer", "devcontainer.json"), []byte{}, 0755)

	os.MkdirAll(filepath.Join(root, "test2", ".devcontainer"), 0755)
	ioutil.WriteFile(filepath.Join(root, "test2", ".devcontainer", "devcontainer.json"), []byte{}, 0755)

	templates, err := getTemplatesFromFolders(folders)
	assert.NoError(t, err)

	expectedTemplates := []DevcontainerTemplate{
		{
			Name: "test1",
			Path: filepath.Join(root, "test1", ".devcontainer"),
		},
		{
			Name: "test2",
			Path: filepath.Join(root, "test2", ".devcontainer"),
		},
	}

	assert.ElementsMatch(t, expectedTemplates, templates)
}
func TestGetTemplateFolders_TakesFolderInPrioirtyOrder(t *testing.T) {

	root1, err := ioutil.TempDir("", "devcontainer*")
	assert.NoError(t, err)
	defer os.RemoveAll(root1)
	root2, err := ioutil.TempDir("", "devcontainer*")
	assert.NoError(t, err)
	defer os.RemoveAll(root2)

	folders := []string{root1, root2}

	os.MkdirAll(filepath.Join(root1, "test1", ".devcontainer"), 0755)
	ioutil.WriteFile(filepath.Join(root1, "test1", ".devcontainer", "devcontainer.json"), []byte{}, 0755)

	os.MkdirAll(filepath.Join(root2, "test1", ".devcontainer"), 0755)
	ioutil.WriteFile(filepath.Join(root2, "test1", ".devcontainer", "devcontainer.json"), []byte{}, 0755)

	templates, err := getTemplatesFromFolders(folders)
	assert.NoError(t, err)

	expectedTemplates := []DevcontainerTemplate{
		{
			Name: "test1",
			Path: filepath.Join(root1, "test1", ".devcontainer"),
		},
	}

	assert.ElementsMatch(t, expectedTemplates, templates)
}
func TestGetTemplateFolders_IgnoresFolderWithoutDevcontainer(t *testing.T) {

	root1, err := ioutil.TempDir("", "devcontainer*")
	assert.NoError(t, err)
	defer os.RemoveAll(root1)
	root2, err := ioutil.TempDir("", "devcontainer*")
	assert.NoError(t, err)
	defer os.RemoveAll(root2)

	folders := []string{root1, root2}

	os.MkdirAll(filepath.Join(root1, "test1", ".devcontainer"), 0755)
	ioutil.WriteFile(filepath.Join(root1, "test1", ".devcontainer", "not-a-devcontainer.json"), []byte{}, 0755)

	os.MkdirAll(filepath.Join(root2, "test1", ".devcontainer"), 0755)
	ioutil.WriteFile(filepath.Join(root2, "test1", ".devcontainer", "devcontainer.json"), []byte{}, 0755)

	templates, err := getTemplatesFromFolders(folders)
	assert.NoError(t, err)

	expectedTemplates := []DevcontainerTemplate{
		{
			Name: "test1",
			Path: filepath.Join(root2, "test1", ".devcontainer"), // Takes root2 because root1 doesn't have devcontainer.json
		},
	}

	assert.ElementsMatch(t, expectedTemplates, templates)
}
