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

	f.WriteString(`{
	"name": "initial",
	// here's a comment!
	"otherProperties": [
		"something",
		"here"
	]
}`)

	SetDevcontainerName(f.Name(), "newName")

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
