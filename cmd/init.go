package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

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
		language := Project.AutoDetectLanguage()

		log.Printf("Setting the project language to %s\n", language)
		Project.SetLanguage(language)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
