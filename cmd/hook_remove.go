package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var removeHookCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a git hook for this ide project",
	Long:  "Remove a git hook for this ide project",
	Args:  cobra.MinimumNArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if err := loadProject(cmd, args); err == nil {
			return []string{}, cobra.ShellCompDirectiveNoFileComp
		}
		return project.ListHooks(), cobra.ShellCompDirectiveNoFileComp
	},
	PreRunE: loadProject,
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, hook := range args {
			err := project.RemoveHook(hook)
			if err != nil {
				log.Println(err)
			} else {
				log.Printf("Hook %s removed\n", hook)
			}
		}

		return nil
	},
}

func init() {
	hookCmd.AddCommand(removeHookCmd)
}
