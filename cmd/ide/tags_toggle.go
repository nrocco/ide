package main

import (
	"github.com/spf13/cobra"
)

var tagsToggleCmdOpts struct {
	On  bool
	Off bool
}

var tagsToggleCmd = &cobra.Command{
	Use:   "toggle",
	Short: "Toggle the notags option",
	Long:  "Toggle the notags option. When enabled, tag generation is disabled for this project.",
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	PreRunE: loadProject,
	RunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Flags().Changed("on") {
			return project.SetNoTags(tagsToggleCmdOpts.On)
		} else if cmd.Flags().Changed("off") {
			return project.SetNoTags(tagsToggleCmdOpts.Off)
		}
		return project.SetNoTags(!project.NoTags())
	},
}

func init() {
	tagsToggleCmd.Flags().BoolVar(&tagsToggleCmdOpts.On, "on", false, "Enable notags")
	tagsToggleCmd.Flags().BoolVar(&tagsToggleCmdOpts.Off, "off", false, "Disable notags")
	tagsToggleCmd.MarkFlagsMutuallyExclusive("on", "off")

	tagsCmd.AddCommand(tagsToggleCmd)
}
