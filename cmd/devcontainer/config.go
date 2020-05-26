package main

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stuartleeks/devcontainer-cli/internal/pkg/config"
)

func createConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "config",
	}
	cmd.AddCommand(createConfigShowCommand())
	cmd.AddCommand(createConfigWriteCommand())
	return cmd
}
func createConfigShowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "show the current config",
		Long:  "load the current config and print it out",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := config.GetAll()
			jsonConfig, err := json.MarshalIndent(c, "", "  ")
			if err != nil {
				return fmt.Errorf("Error converting to JSON: %s\n", err)
			}
			fmt.Println(string(jsonConfig))
			return nil
		},
	}
	return cmd
}
func createConfigWriteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "write",
		Short: "write config",
		Long:  "Write out the config file to ~/.devcontainer-cli/devcontainer-cli.json",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := config.SaveConfig(); err != nil {
				return fmt.Errorf("Error saving config: %s\n", err)
			}
			fmt.Println("Config saved")
			return nil
		},
	}
	return cmd

}
