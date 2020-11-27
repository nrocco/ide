package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a git repository as an ide project",
	Long:  "Initialize a git repository as an ide project",
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	PreRunE: loadProject,
	RunE: func(cmd *cobra.Command, args []string) error {
		if project.IsConfigured() {
			log.Fatalln("The repository is already setup")
			return nil
		}

		log.Println("Setting up the repository as a ide project...")
		language := project.AutoDetectLanguage()

		log.Printf("Setting the project language to %s\n", language)
		project.SetLanguage(language)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
