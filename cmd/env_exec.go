package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

var execEnvCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute a binary in this ide project's environment",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return nil // TODO return error here instead
		}

		container := Project.GetExecutable(args[0])
		options := append([]string{"exec", "-i", "-t", container}, args...)

		command := exec.Command("docker", options...)
		command.Stdout = os.Stdout
		command.Stdin = os.Stdin
		command.Stderr = os.Stderr

		return command.Run()
	},
}

func init() {
	envCmd.AddCommand(execEnvCmd)
}
