package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get the current status of your ide project",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !Project.IsConfigured() {
			log.Fatalf("The project is not configured yet\n")
			return nil
		}

		fmt.Printf("Ide\n")
		fmt.Printf("  Name:     %s\n", Project.Name())
		fmt.Printf("  Branch:   %s\n", Project.Branch())
		fmt.Printf("  Language: %s\n", Project.Language())
		fmt.Printf("  Location: %s\n", Project.Location())
		fmt.Printf("  Hooks:    %s\n", strings.Join(Project.ListHooks(), " "))

		return nil
	},
}

func init() {
	RootCmd.AddCommand(statusCmd)
}
