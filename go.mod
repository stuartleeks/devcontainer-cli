module github.com/stuartleeks/devcontainer-cli

go 1.14

require (
	github.com/blang/semver v3.5.1+incompatible
	github.com/bradford-hamilton/dora v0.1.1
	github.com/kyoh86/richgo v0.3.10 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/rhysd/go-github-selfupdate v1.2.2
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.7.0
	golang.org/x/sys v0.0.0-20220114195835-da31bd327af9 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/bradford-hamilton/dora v0.1.1 => github.com/stuartleeks/dora v0.1.5
