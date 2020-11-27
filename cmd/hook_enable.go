package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var enableHookCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable a git hook for this ide project",
	Long:  "Enable a git hook for this ide project",
	Args:  cobra.MinimumNArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"commit-msg", "post-checkout", "post-merge", "prepare-commit-msg"}, cobra.ShellCompDirectiveNoFileComp // TODO make this dynamic, hide already enabled hooks
	},
	PreRunE: loadProject,
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, hook := range args {
			err := project.EnableHook(hook)
			if err != nil {
				log.Println(err)
			} else {
				log.Printf("Hook %s enabled\n", hook)
			}
		}

		return nil
	},
}

func init() {
	hookCmd.AddCommand(enableHookCmd)
}
