package cmd

import (
	"errors"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var runExecCmd = &cobra.Command{
	Use:   "exec [name]",
	Short: "Execute a program in this ide project's environment",
	Long:  "Execute a program in this ide project's environment",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Usage()
			return errors.New("You must supply the name of the executable to run")
		}

		executable, err := project.GetExecutable(args[0])

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
	execCmd.AddCommand(runExecCmd)
}
