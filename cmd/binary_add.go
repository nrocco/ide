package cmd

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
)

var addBinaryCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a binary for this ide project",
	Long:  "Add a binary for this ide project",
	Args:  cobra.MinimumNArgs(2),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	PreRunE: loadProject,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := project.AddBinary(args[0], strings.Join(args[1:], " "))
		if err != nil {
			return err
		}

		log.Printf("Binary %s add\n", args[0])

		return project.AddGitBinToPath()
	},
}

func init() {
	binaryCmd.AddCommand(addBinaryCmd)
}
