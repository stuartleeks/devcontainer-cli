module github.com/stuartleeks/devcontainer-cli

go 1.14

require (
	github.com/blang/semver v3.5.1+incompatible
	github.com/bradford-hamilton/dora v0.1.1
	github.com/rhysd/go-github-selfupdate v1.2.2
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.8.2
	golang.org/x/sys v0.10.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/bradford-hamilton/dora v0.1.1 => github.com/stuartleeks/dora v0.1.5
