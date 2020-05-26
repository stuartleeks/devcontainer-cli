package main

import (
	"bufio"
	"fmt"
	"log"
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
		Run: func(cmd *cobra.Command, args []string) {
			latest, err := update.CheckForUpdate(version)
			if err != nil {
				log.Println("Error occurred while checking for updates:", err)
				return
			}

			if latest == nil {
				fmt.Println("No updates available")
				return
			}

			fmt.Printf("\n\n UPDATE AVAILABLE: %s \n \n Release notes: %s\n", latest.Version, latest.ReleaseNotes)

			if checkOnly {
				return
			}

			fmt.Print("Do you want to update? (y/n): ")
			if !yes {
				input, err := bufio.NewReader(os.Stdin).ReadString('\n')
				if err != nil || (input != "y\n" && input != "y\r\n") {
					// error or something other than `y`
					return
				}
			}
			fmt.Println("Applying...")

			exe, err := os.Executable()
			if err != nil {
				log.Panicln("Could not locate executable path")
			}
			if err := selfupdate.UpdateTo(latest.AssetURL, exe); err != nil {
				log.Panicln("Error occurred while updating binary:", err)
			}
			log.Println("Successfully updated to version", latest.Version)

		},
	}
	cmd.Flags().BoolVar(&checkOnly, "check-only", false, "Check for an update without applying")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Automatically apply any updates (i.e. answer yes) ")

	return cmd
}
