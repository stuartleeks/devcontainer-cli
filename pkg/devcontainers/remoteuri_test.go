package devcontainers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHexToString(t *testing.T) {
	input := "\\\\wsl$\\Ubuntusl\\home\\stuart\\source\\kips-operator"
	expected := "5c5c77736c245c5562756e7475736c5c686f6d655c7374756172745c736f757263655c6b6970732d6f70657261746f72"
	actual := convertToHexString(input)
	assert.Equal(t, expected, actual)
}

func TestGetWorkspaceFolder_withWorkspaceFolder(t *testing.T) {

	content := `{
		"someProp": 2,
		// add a content here for good measure
		"workspaceFolder": "/workspace/wibble",
	}`
	result, err := getWorkspaceMountPathFromDevcontainerDefinition([]byte(content))

	assert.NoError(t, err)
	assert.Equal(t, "/workspace/wibble", result)
}
func TestGetWorkspaceFolder_withCommentedWorkspaceFolder(t *testing.T) {

	content := `{
		"someProp": 2,
		// add a content here for good measure
		//"workspaceFolder": "/workspace/wibble",
	}`
	result, err := getWorkspaceMountPathFromDevcontainerDefinition([]byte(content))

	assert.NoError(t, err)
	assert.Equal(t, "", result)
}
func TestGetWorkspaceFolder_withNoWorkspaceFolder(t *testing.T) {

	content := `{
		"someProp": 2,
		// add a content here for good measure
	}`
	result, err := getWorkspaceMountPathFromDevcontainerDefinition([]byte(content))

	assert.NoError(t, err)
	assert.Equal(t, "", result)
}
