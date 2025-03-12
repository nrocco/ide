package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var addHookCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a git hook for this ide project",
	Long:  "Add a git hook for this ide project",
	Args:  cobra.MinimumNArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"commit-msg", "prepare-commit-msg"}, cobra.ShellCompDirectiveNoFileComp
	},
	PreRunE: loadProject,
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, hook := range args {
			err := project.HookAdd(hook)
			if err != nil {
				log.Println(err)
			} else {
				log.Printf("Hook %s add\n", hook)
			}
		}

		return nil
	},
}

func init() {
	hookCmd.AddCommand(addHookCmd)
}
