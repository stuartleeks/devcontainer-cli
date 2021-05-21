package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitialise(t *testing.T) {

	err := Initialise(`{
		"templatepaths": ["test"]
}`)
	if assert.NoError(t, err) {
		return
	}

	templatePaths := GetTemplateFolders()
	if !assert.Equal(t, 1, len(templatePaths)) {
		return
	}

	assert.Equal(t, "test", templatePaths[0])
}
