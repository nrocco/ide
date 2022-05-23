package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var ctagsCmd = &cobra.Command{
	Use:   "ctags",
	Short: "Update ctags for the ide project",
	Long:  "Update ctags for the ide project",
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	PreRunE: loadProject,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Printf("Writing ctags file to %s\n", project.CtagsFile())

		return project.RefreshCtags()
	},
}

func init() {
	rootCmd.AddCommand(ctagsCmd)
}
