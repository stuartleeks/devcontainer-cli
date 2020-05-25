package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
)

func createUpdateCommand() *cobra.Command {

	var checkOnly bool
	var yes bool

	cmd := &cobra.Command{
		Use:   "update",
		Short: "update cli",
		Long:  "Apply the latest update",
		Run: func(cmd *cobra.Command, args []string) {
			latest, found, err := selfupdate.DetectLatest("stuartleeks/devcontainer-cli")
			if err != nil {
				log.Println("Error occurred while detecting version:", err)
				return
			}

			v, err := semver.Parse(version)
			if err != nil {
				log.Panicln(err.Error())
			}
			if !found || latest.Version.LTE(v) {
				log.Println("Current version is the latest")
				return
			}

			fmt.Print("\n\n UPDATE AVAILABLE \n \n Release notes: "+latest.ReleaseNotes+" \n Do you want to update to: ", latest.Version, "? (y/n): ")
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
	cmd.Flags().BoolVarP(&checkOnly, "yes", "y", false, "Automatically apply any updates (i.e. answer yes) ")

	return cmd
}
