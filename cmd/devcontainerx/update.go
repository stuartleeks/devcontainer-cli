package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
	"github.com/stuartleeks/devcontainer-cli/internal/pkg/update"
)

func createUpdateCommand() *cobra.Command {

	var checkOnly bool
	var yes bool

	cmd := &cobra.Command{
		Use:   "update",
		Short: "update cli",
		Long:  "Apply the latest update",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// do nothing - suppress root PersistentPreRun which does periodic update check
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			latest, err := update.CheckForUpdate(version)
			if err != nil {
				return fmt.Errorf("Error occurred while checking for updates: %v", err)
			}

			if latest == nil {
				fmt.Println("No updates available")
				return nil
			}

			fmt.Printf("\n\n UPDATE AVAILABLE: %s \n \n Release notes: %s\n", latest.Version, latest.ReleaseNotes)

			if checkOnly {
				return nil
			}

			fmt.Print("Do you want to update? (y/n): ")
			if !yes {
				input, err := bufio.NewReader(os.Stdin).ReadString('\n')
				if err != nil || (input != "y\n" && input != "y\r\n") {
					// error or something other than `y`
					return err
				}
			}
			fmt.Println("Applying...")

			exe, err := os.Executable()
			if err != nil {
				return fmt.Errorf("Could not locate executable path: %v", err)
			}
			if err := selfupdate.UpdateTo(latest.AssetURL, exe); err != nil {
				return fmt.Errorf("Error occurred while updating binary: %v", err)
			}
			fmt.Printf("Successfully updated to version %s\n", latest.Version)
			return nil
		},
	}
	cmd.Flags().BoolVar(&checkOnly, "check-only", false, "Check for an update without applying")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Automatically apply any updates (i.e. answer yes) ")

	return cmd
}
