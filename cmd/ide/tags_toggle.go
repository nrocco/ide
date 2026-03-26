package main

import (
	"github.com/spf13/cobra"
)

var tagsToggleCmd = &cobra.Command{
	Use:   "toggle",
	Short: "Toggle the notags option",
	Long:  "Toggle the notags option. When enabled, tag generation is disabled for this project.",
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	PreRunE: loadProject,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !cmd.Flags().Changed("on") && !cmd.Flags().Changed("off") {
			return project.SetNoTags(!project.NoTags())
		}

		off, _ := cmd.Flags().GetBool("off")
		return project.SetNoTags(off)
	},
}

func init() {
	tagsToggleCmd.Flags().Bool("on", false, "Enable notags")
	tagsToggleCmd.Flags().Bool("off", false, "Disable notags")
	tagsToggleCmd.MarkFlagsMutuallyExclusive("on", "off")

	tagsCmd.AddCommand(tagsToggleCmd)
}
