package main

import (
	"encoding/json"
	"fmt"
	"os"

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
		Run: func(cmd *cobra.Command, args []string) {
			c := config.GetAll()
			jsonConfig, err := json.MarshalIndent(c, "", "  ")
			if err != nil {
				fmt.Printf("Error converting to JSON: %s\n", err)
				os.Exit(1)
			}
			fmt.Println(string(jsonConfig))
		},
	}
	return cmd
}
func createConfigWriteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "write",
		Short: "write config",
		Long:  "Write out the config file to ~/.devcontainer-cli/devcontainer-cli.json",
		Run: func(cmd *cobra.Command, args []string) {
			if err := config.SaveConfig(); err != nil {
				fmt.Printf("Error saving config: %s\n", err)
			} else {
				fmt.Println("Config saved")
			}
		},
	}
	return cmd

}
