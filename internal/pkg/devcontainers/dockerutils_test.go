package devcontainers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetClosestPathMatchForPath_ReturnsLongestMatch(t *testing.T) {

	inputs := []DevcontainerInfo{
		{LocalFolderPath: "/path/to/project"},
		{LocalFolderPath: "/path/to/somewhere/else"},
		{LocalFolderPath: "/path"},
	}

	actual, err := GetClosestPathMatchForPath(inputs, "/path/to/project")
	if assert.NoError(t, err) {
		assert.Equal(t, "/path/to/project", actual.LocalFolderPath)
	}
}

func TestGetClosestPathMatchForPath_ReturnsLongestMatchWithTrailingSlash(t *testing.T) {

	inputs := []DevcontainerInfo{
		{LocalFolderPath: "/path/to/project"},
		{LocalFolderPath: "/path/to/somewhere/else"},
		{LocalFolderPath: "/path"},
	}

	actual, err := GetClosestPathMatchForPath(inputs, "/path/to/project/")
	if assert.NoError(t, err) {
		assert.Equal(t, "/path/to/project", actual.LocalFolderPath)
	}
}

func TestGetClosestPathMatchForPath_ReturnsLongestMatchForChildFolder(t *testing.T) {

	inputs := []DevcontainerInfo{
		{LocalFolderPath: "/path/to/project"},
		{LocalFolderPath: "/path/to/somewhere/else"},
		{LocalFolderPath: "/path"},
	}

	actual, err := GetClosestPathMatchForPath(inputs, "/path/to/project/with/child")
	if assert.NoError(t, err) {
		assert.Equal(t, "/path/to/project", actual.LocalFolderPath)
	}
}
