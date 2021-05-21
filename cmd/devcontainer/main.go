package main

import (
	"github.com/spf13/cobra"
	"github.com/stuartleeks/devcontainer-cli/internal/pkg/update"
	"github.com/stuartleeks/devcontainer-cli/pkg/config"
)

// Overridden via ldflags
var (
	version   = "99.0.1-devbuild"
	commit    = "unknown"
	date      = "unknown"
	goversion = "unknown"
)

func main() {

	rootCmd := &cobra.Command{
		Use: "devcontainer",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			update.PeriodicCheckForUpdate(version)
		},
	}

	rootCmd.AddCommand(createCompleteCommand(rootCmd))
	rootCmd.AddCommand(createConfigCommand())
	rootCmd.AddCommand(createExecCommand())
	rootCmd.AddCommand(createListCommand())
	rootCmd.AddCommand(createTemplateCommand())
	if config.GetExperimentalFeaturesEnabled() {
		rootCmd.AddCommand(createSnippetCommand())
	}
	rootCmd.AddCommand(createUpdateCommand())
	rootCmd.AddCommand(createOpenInCodeCommand())
	rootCmd.AddCommand(createOpenInCodeInsidersCommand())
	rootCmd.AddCommand(createVersionCommand())

	_ = rootCmd.Execute()
}
