package devcontainers

import (
	"io/ioutil"
	"os"
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
