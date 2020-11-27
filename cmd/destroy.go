package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Remove all ide configuration for a repository",
	Long:  "Remove all ide configuration for a repository",
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	PreRunE: loadProject,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !project.IsConfigured() {
			log.Fatalf("The repository is already destroyed")
			return nil
		}

		err := project.Destroy()
		if err != nil {
			return err
		}

		log.Println("Repository is no longer an ide project")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)
}
