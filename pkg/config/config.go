package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var initialised bool = false

// ensureInitialised reads the config. Will quit if config is invalid
func ensureInitialised() {
	if !initialised {
		setupViper()

		viper.AddConfigPath(getConfigPath())
		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				fmt.Printf("Error loading config file: %s\n", err)
				os.Exit(1)
			}
		}
		initialised = true
	}
}

// Initialise provides a way to load config from a custom source
func Initialise(configJSON string) error {
	if initialised {
		return fmt.Errorf("already initialised")
	}
	reader := strings.NewReader(configJSON)
	setupViper()
	if err := viper.ReadConfig(reader); err != nil {
		return fmt.Errorf("error reading config: %s", err)
	}
	initialised = true
	return nil
}

func setupViper() {
	viper.SetConfigName("devcontainer-cli")
	viper.SetConfigType("json")

	viper.SetDefault("templatePaths", []string{})
	viper.SetDefault("settingPaths", []string{})
	viper.SetDefault("experimental", false)
}
func getConfigPath() string {
	// TODO - allow env var for config
	var path string
	if os.Getenv("HOME") != "" {
		path = filepath.Join("$HOME", ".devcontainer-cli/")
	} else {
		// if HOME not set, assume Windows and use USERPROFILE env var
		path = filepath.Join("$USERPROFILE", ".devcontainer-cli/")
	}
	return os.ExpandEnv(path)
}

func GetTemplateFolders() []string {
	ensureInitialised()
	return viper.GetStringSlice("templatePaths")
}
func GetSnippetFolders() []string {
	ensureInitialised()
	return viper.GetStringSlice("snippetPaths")
}
func GetExperimentalFeaturesEnabled() bool {
	ensureInitialised()
	return viper.GetBool("experimental")
}
func GetLastUpdateCheck() time.Time {
	ensureInitialised()
	return viper.GetTime("lastUpdateCheck")
}
func SetLastUpdateCheck(t time.Time) {
	ensureInitialised()
	viper.Set("lastUpdateCheck", t)
}
func GetAll() map[string]interface{} {
	ensureInitialised()
	return viper.AllSettings()
}

func SaveConfig() error {
	ensureInitialised()
	configPath := getConfigPath()
	configPath = os.ExpandEnv(configPath)
	if err := os.MkdirAll(configPath, 0755); err != nil {
		return err
	}
	configFilePath := filepath.Join(configPath, "devcontainer-cli.json")
	return viper.WriteConfigAs(configFilePath)
}
