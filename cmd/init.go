package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// initCmd represents the status command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a git repository as an ide project",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if Project.IsConfigured() {
			log.Fatalln("The repository is already setup")
			return nil
		}

		log.Println("Setting up the repository as a ide project...")

		Project.Language = Project.AutoDetectLanguage()

		log.Printf("Setting the project language to %s\n", Project.Language)

		Project.Save()

		return nil
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}
