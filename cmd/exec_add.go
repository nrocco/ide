package cmd

import (
	"errors"
	"log"

	"github.com/spf13/cobra"
)

var execAddCmd = &cobra.Command{
	Use:   "add [program] [name]",
	Short: "Add an executable to this ide project",
	Long:  "Add an executable to this ide project",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Usage()
			return errors.New("You must supply a program to add")
		} else if len(args) > 2 {
			cmd.Usage()
			return errors.New("Only 2 arguments expected")
		}

		if len(args) == 1 {
			args = append(args, "")
		}

		executable, err := project.NewExecutable(args[0], args[1])
		if err != nil {
			return err
		}

		log.Printf("Created executable %s in %s", executable.Name(), executable.Target())

		return nil
	},
}

func init() {
	execCmd.AddCommand(execAddCmd)
}
