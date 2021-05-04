package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

var initialised bool = false

// EnsureInitialised reads the config. Will quit if config is invalid
func EnsureInitialised() {
	if !initialised {
		viper.SetConfigName("devcontainer-cli")
		viper.SetConfigType("json")

		viper.AddConfigPath(getConfigPath())

		viper.SetDefault("templatePaths", []string{})

		// TODO - allow env var for config
		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				// Config file not found; ignore error if desired
			} else {
				fmt.Printf("Error loading config file: %s\n", err)
				os.Exit(1)
			}
		}
		initialised = true
	}
}
func getConfigPath() string {
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
	EnsureInitialised()
	return viper.GetStringSlice("templatePaths")
}
func GetSnippetFolders() []string {
	EnsureInitialised()
	return viper.GetStringSlice("snippetPaths")
}
func GetLastUpdateCheck() time.Time {
	EnsureInitialised()
	return viper.GetTime("lastUpdateCheck")
}
func SetLastUpdateCheck(t time.Time) {
	EnsureInitialised()
	viper.Set("lastUpdateCheck", t)
}
func GetAll() map[string]interface{} {
	EnsureInitialised()
	return viper.AllSettings()
}

func SaveConfig() error {
	EnsureInitialised()
	configPath := getConfigPath()
	configPath = os.ExpandEnv(configPath)
	if err := os.MkdirAll(configPath, 0755); err != nil {
		return err
	}
	configFilePath := filepath.Join(configPath, "devcontainer-cli.json")
	return viper.WriteConfigAs(configFilePath)
}
