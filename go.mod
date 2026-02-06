module github.com/stuartleeks/devcontainer-cli

go 1.23.0

toolchain go1.23.7

require (
	github.com/blang/semver v3.5.1+incompatible
	github.com/bradford-hamilton/dora v0.1.1
	github.com/fatih/color v1.18.0
	github.com/jmespath/go-jmespath v0.4.0
	github.com/mattn/go-colorable v0.1.14
	github.com/neilotoole/jsoncolor v0.7.1
	github.com/rhysd/go-github-selfupdate v1.2.2
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.8.4
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fsnotify/fsnotify v1.4.7 // indirect
	github.com/google/go-github/v30 v30.1.0 // indirect
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/go-update v0.0.0-20160112193335-8152e7eb6ccf // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/magiconair/properties v1.8.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/mapstructure v1.1.2 // indirect
	github.com/pelletier/go-toml v1.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/spf13/afero v1.1.2 // indirect
	github.com/spf13/cast v1.3.0 // indirect
	github.com/spf13/jwalterweatherman v1.0.0 // indirect
	github.com/spf13/pflag v1.0.3 // indirect
	github.com/tcnksm/go-gitconfig v0.1.2 // indirect
	github.com/ulikunitz/xz v0.5.5 // indirect
	golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2 // indirect
	golang.org/x/oauth2 v0.27.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/term v0.14.0 // indirect
	golang.org/x/text v0.3.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/bradford-hamilton/dora v0.1.1 => github.com/stuartleeks/dora v0.1.5

replace github.com/jmespath/go-jmespath => github.com/stuartleeks/go-jmespath v0.0.1
