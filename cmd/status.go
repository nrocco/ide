package cmd

import (
	"fmt"
	"log"
	"strings"

	humanize "github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get the current status of your ide project",
	Long:  "Get the current status of your ide project",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !project.IsConfigured() {
			log.Fatalf("The project is not configured yet\n")
			return nil
		}

		fmt.Printf("Ide\n")
		fmt.Printf("  Name: %s\n", project.Name())
		fmt.Printf("  Branch: %s\n", project.Branch())
		fmt.Printf("  Language: %s\n", project.Language())
		fmt.Printf("  Location: %s\n", project.Location())
		fmt.Printf("  Ctags:\n")
		fmt.Printf("    File: %s\n", project.CtagsFile())
		fmt.Printf("    Age: %s\n", humanize.Time(project.CtagsFileAge()))
		fmt.Printf("    Size: %s\n", humanize.Bytes(project.CtagsFileSize()))
		fmt.Printf("  Hooks: %s\n", strings.Join(project.ListHooks(), " "))
		fmt.Printf("  Executables: %s\n", strings.Join(project.ListExecutables(), " "))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
