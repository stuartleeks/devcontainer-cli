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
	if !assert.NoError(t, err) {
		return
	}
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
	if !assert.NoError(t, err) {
		return
	}

	buf, err := ioutil.ReadFile(f.Name())
	if !assert.NoError(t, err) {
		return
	}

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
	if !assert.NoError(t, err) {
		return
	}
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
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, "vscode", user)
}

func TestGetDevContainerUserName_NotSet(t *testing.T) {

	f, err := ioutil.TempFile("", "test.json")
	if !assert.NoError(t, err) {
		return
	}
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
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, "", user)
}

func TestGetDevContainerUserName_Commented(t *testing.T) {

	f, err := ioutil.TempFile("", "test.json")
	if !assert.NoError(t, err) {
		return
	}
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
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, "", user)
}

func TestSetDevcontainerName_SubstitutionValue(t *testing.T) {

	f, err := ioutil.TempFile("", "test.json")
	if !assert.NoError(t, err) {
		return
	}
	defer os.Remove(f.Name())

	_, _ = f.WriteString(`{
	"name": "initial",
	// here's a comment!
	"otherProperties": [
		"something-__DEVCONTAINER_NAME__",
		"here-__DEVCONTAINER_NAME__"
	]
}`)

	err = SetDevcontainerName(f.Name(), "newName")
	if !assert.NoError(t, err) {
		return
	}

	buf, err := ioutil.ReadFile(f.Name())
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, `{
	"name": "newName",
	// here's a comment!
	"otherProperties": [
		"something-newName",
		"here-newName"
	]
}`, string(buf))
}

func TestGetTemplateFolders_ListsFoldersWithDevcontainers(t *testing.T) {

	root, err := ioutil.TempDir("", "devcontainer*")
	if !assert.NoError(t, err) {
		return
	}
	defer os.RemoveAll(root)

	folders := []string{root}

	_ = os.MkdirAll(filepath.Join(root, "test1", ".devcontainer"), 0755)
	_ = ioutil.WriteFile(filepath.Join(root, "test1", ".devcontainer", "devcontainer.json"), []byte{}, 0755)

	_ = os.MkdirAll(filepath.Join(root, "test2", ".devcontainer"), 0755)
	_ = ioutil.WriteFile(filepath.Join(root, "test2", ".devcontainer", "devcontainer.json"), []byte{}, 0755)

	templates, err := getTemplatesFromFolders(folders)
	if !assert.NoError(t, err) {
		return
	}

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
	if !assert.NoError(t, err) {
		return
	}
	defer os.RemoveAll(root1)

	root2, err := ioutil.TempDir("", "devcontainer*")
	if !assert.NoError(t, err) {
		return
	}
	defer os.RemoveAll(root2)

	folders := []string{root1, root2}

	_ = os.MkdirAll(filepath.Join(root1, "test1", ".devcontainer"), 0755)
	_ = ioutil.WriteFile(filepath.Join(root1, "test1", ".devcontainer", "devcontainer.json"), []byte{}, 0755)

	_ = os.MkdirAll(filepath.Join(root2, "test1", ".devcontainer"), 0755)
	_ = ioutil.WriteFile(filepath.Join(root2, "test1", ".devcontainer", "devcontainer.json"), []byte{}, 0755)

	templates, err := getTemplatesFromFolders(folders)
	if !assert.NoError(t, err) {
		return
	}

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
	if !assert.NoError(t, err) {
		return
	}
	defer os.RemoveAll(root1)

	root2, err := ioutil.TempDir("", "devcontainer*")
	if !assert.NoError(t, err) {
		return
	}
	defer os.RemoveAll(root2)

	folders := []string{root1, root2}

	_ = os.MkdirAll(filepath.Join(root1, "test1", ".devcontainer"), 0755)
	_ = ioutil.WriteFile(filepath.Join(root1, "test1", ".devcontainer", "not-a-devcontainer.json"), []byte{}, 0755)

	_ = os.MkdirAll(filepath.Join(root2, "test1", ".devcontainer"), 0755)
	_ = ioutil.WriteFile(filepath.Join(root2, "test1", ".devcontainer", "devcontainer.json"), []byte{}, 0755)

	templates, err := getTemplatesFromFolders(folders)
	if !assert.NoError(t, err) {
		return
	}

	expectedTemplates := []DevcontainerTemplate{
		{
			Name: "test1",
			Path: filepath.Join(root2, "test1", ".devcontainer"), // Takes root2 because root1 doesn't have devcontainer.json
		},
	}

	assert.ElementsMatch(t, expectedTemplates, templates)
}

func TestAddTemplate_PerformsSubstitutionWithUserName(t *testing.T) {

	root, err := ioutil.TempDir("", "devcontainer*")
	if !assert.NoError(t, err) {
		return
	}
	defer os.RemoveAll(root)

	_ = os.MkdirAll(filepath.Join(root, "test1", ".devcontainer"), 0755)
	_ = ioutil.WriteFile(filepath.Join(root, "test1", ".devcontainer", "devcontainer.json"), []byte(`{
	"name": "expect this to be replaced",
	"settings": {
		"DC_NAME": "__DEVCONTAINER_NAME__",
		"DC_USER_NAME": "__DEVCONTAINER_USER_NAME__",
		"DC_HOME": "__DEVCONTAINER_HOME__"
	},
	"remoteUser": "dcuser"
}`), 0755)
	_ = ioutil.WriteFile(filepath.Join(root, "test1", ".devcontainer", "Dockerfile"), []byte(`FROM foo
RUN echo hi

ENV DC_NAME=__DEVCONTAINER_NAME__
ENV DC_USER_NAME=__DEVCONTAINER_USER_NAME__
ENV DC_HOME=__DEVCONTAINER_HOME__

# __DEVCONTAINER_SNIPPET_INSERT__ 

RUN echo hi2
`), 0755)

	// set up devcontainer
	targetFolder := filepath.Join(root, "target")
	_ = os.MkdirAll(targetFolder, 0755)

	// Add template
	err = CopyTemplateToFolder(filepath.Join(root, "test1", ".devcontainer"), targetFolder, "NewName")
	if !assert.NoError(t, err) {
		return
	}

	devcontainerFolder := filepath.Join(targetFolder, ".devcontainer")
	buf, err := ioutil.ReadFile(filepath.Join(devcontainerFolder, "Dockerfile"))
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, `FROM foo
RUN echo hi

ENV DC_NAME=NewName
ENV DC_USER_NAME=dcuser
ENV DC_HOME=/home/dcuser

# __DEVCONTAINER_SNIPPET_INSERT__ 

RUN echo hi2
`, string(buf))

	buf, err = ioutil.ReadFile(filepath.Join(devcontainerFolder, "devcontainer.json"))
	if !assert.NoError(t, err) {
		return
	}
	stringContent := string(buf)
	assert.Equal(t, `{
	"name": "NewName",
	"settings": {
		"DC_NAME": "NewName",
		"DC_USER_NAME": "dcuser",
		"DC_HOME": "/home/dcuser"
	},
	"remoteUser": "dcuser"
}`, stringContent)

}
func TestAddTemplate_PerformsSubstitutionWithoutUserName(t *testing.T) {

	root, err := ioutil.TempDir("", "devcontainer*")
	if !assert.NoError(t, err) {
		return
	}
	defer os.RemoveAll(root)

	_ = os.MkdirAll(filepath.Join(root, "test1", ".devcontainer"), 0755)
	_ = ioutil.WriteFile(filepath.Join(root, "test1", ".devcontainer", "devcontainer.json"), []byte(`{
	"name": "expect this to be replaced",
	"settings": {
		"DC_NAME": "__DEVCONTAINER_NAME__",
		"DC_USER_NAME": "__DEVCONTAINER_USER_NAME__",
		"DC_HOME": "__DEVCONTAINER_HOME__"
	},
}`), 0755)
	_ = ioutil.WriteFile(filepath.Join(root, "test1", ".devcontainer", "Dockerfile"), []byte(`FROM foo
RUN echo hi

ENV DC_NAME=__DEVCONTAINER_NAME__
ENV DC_USER_NAME=__DEVCONTAINER_USER_NAME__
ENV DC_HOME=__DEVCONTAINER_HOME__

# __DEVCONTAINER_SNIPPET_INSERT__ 

RUN echo hi2
`), 0755)

	// set up devcontainer
	targetFolder := filepath.Join(root, "target")
	_ = os.MkdirAll(targetFolder, 0755)

	// Add template
	err = CopyTemplateToFolder(filepath.Join(root, "test1", ".devcontainer"), targetFolder, "NewName")
	if !assert.NoError(t, err) {
		return
	}

	devcontainerFolder := filepath.Join(targetFolder, ".devcontainer")
	buf, err := ioutil.ReadFile(filepath.Join(devcontainerFolder, "Dockerfile"))
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, `FROM foo
RUN echo hi

ENV DC_NAME=NewName
ENV DC_USER_NAME=root
ENV DC_HOME=/root

# __DEVCONTAINER_SNIPPET_INSERT__ 

RUN echo hi2
`, string(buf))

	buf, err = ioutil.ReadFile(filepath.Join(devcontainerFolder, "devcontainer.json"))
	if !assert.NoError(t, err) {
		return
	}
	stringContent := string(buf)
	assert.Equal(t, `{
	"name": "NewName",
	"settings": {
		"DC_NAME": "NewName",
		"DC_USER_NAME": "root",
		"DC_HOME": "/root"
	},
}`, stringContent)

}
