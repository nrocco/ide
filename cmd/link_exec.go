package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

var execLinkCmd = &cobra.Command{
	Use:   "exec [name]",
	Short: "Execute a linked program in this ide project's environment",
	Long:  "Execute a linked program in this ide project's environment",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Usage()
			return errors.New("You must supply a name for a link to execute")
		}

		executable, err := Project.GetExecutable(args[0])

		if err != nil {
			return err
		}

		if !executable.Containerized() {
			return errors.New("This is not a containerized executable")
		}

		args[0] = executable.Program()

		options := append([]string{"exec", "-i", "-t", executable.Container()}, args...)

		command := exec.Command("docker", options...)
		command.Stdout = os.Stdout
		command.Stdin = os.Stdin
		command.Stderr = os.Stderr

		// TODO return exit code from docker exec

		return command.Run()
	},
}

func init() {
	linkCmd.AddCommand(execLinkCmd)
}
