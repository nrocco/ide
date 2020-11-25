package cmd

import (
	"errors"
	"log"
	"strings"

	"github.com/spf13/cobra"
)

var enableBinaryCmd = &cobra.Command{
	Use:                    "enable",
	Short:                  "Enable a binary for this ide project",
	Long:                   "Enable a binary for this ide project",
	BashCompletionFunction: "_describe 'ide_binary_run_commands' ide_binary_run_commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("You must specify at least two arguments")
		}

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
