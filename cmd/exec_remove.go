package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var removeExecCmd = &cobra.Command{
	Use:   "rm [name]",
	Short: "Remove an executable from this ide project",
	Long:  "Remove an executable from this ide project",
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, name := range args {
			err := project.RemoveExecutable(name)
			if err != nil {
				log.Println(err)
			} else {
				log.Printf("Executable %s removed\n", name)
			}
		}

		return nil
	},
}

func init() {
	execCmd.AddCommand(removeExecCmd)
}
