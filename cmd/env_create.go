package cmd

import (
	"errors"
	"github.com/nrocco/ide/pkg/ide"
	"github.com/spf13/cobra"
)

var createEnvCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an environment for this ide project",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("You must supply an image name")
		}

		environment := ide.Environment{
			Project: Project,
			Image:   args[0],
		}

		return environment.Create()

		// TODO: get default docker create options from ide config file (such as --add-host and --env ssh-auth-sock etc)
		// TODO: use docker go client to create the container (but do not start it)
		// TODO: add current working dir as labels so this container can be better identified

		return nil
	},
}

func init() {
	envCmd.AddCommand(createEnvCmd)
}
