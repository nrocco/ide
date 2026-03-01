package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var runShimCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a shim in the context of an ide project",
	Long:  "Run a shim in the context of an ide project",
	Args:  cobra.MinimumNArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		shims := []string{}
		if len(args) == 0 {
			if err := loadProject(cmd, args); err == nil {
				for shim := range project.ShimList() {
					shims = append(shims, shim)
				}
			}
		}
		return shims, cobra.ShellCompDirectiveNoFileComp
	},
	PreRunE: loadProject,
	RunE: func(cmd *cobra.Command, args []string) error {
		command := project.ShimGet(args[0])
		if command == "" {
			return fmt.Errorf("shim %s does not exist", args[0])
		}
		return project.ShimRun(command, args)
	},
}

func init() {
	shimCmd.AddCommand(runShimCmd)
}
