package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var runPostMergeHookCmd = &cobra.Command{
	Use:   "post-merge",
	Short: "Run the post-merge hook for the ide project",
	Long:  "Run the post-merge hook for the ide project",
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
	runHookCmd.AddCommand(runPostMergeHookCmd)
}
