package cmd

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
)

var enableBinaryCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable a binary for this ide project",
	Long:  "Enable a binary for this ide project",
	Args:  cobra.MinimumNArgs(2),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	PreRunE: loadProject,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := project.EnableBinary(args[0], strings.Join(args[1:], " "))
		if err != nil {
			return err
		}

		log.Printf("Binary %s enabled\n", args[0])

		return nil
	},
}

func init() {
	binaryCmd.AddCommand(enableBinaryCmd)
}
