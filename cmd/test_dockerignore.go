package cmd

import (
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var testDockerIgnoreCmd = &cobra.Command{
	Use:   "test-dockerignore",
	Short: "Test a .dockerignore in the current directory",
	Long:  "Test a .dockerignore in the current directory",
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	PreRunE: loadProject,
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat(".dockerignore"); err != nil {
			return err
		}

		dockerFile := strings.Join([]string{
			"FROM busybox",
			"COPY . /build-context",
			"WORKDIR /build-context",
			"CMD find .",
		}, "\n")

		command := exec.Command("docker", "image", "build", "-t", "build-context", "-f", "-", ".")
		command.Stdin = strings.NewReader(dockerFile)
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		if err := command.Run(); err != nil {
			return err
		}

		command = exec.Command("docker", "container", "run", "--rm", "build-context")
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		if err := command.Run(); err != nil {
			return err
		}

		command = exec.Command("docker", "image", "rm", "build-context")
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		if err := command.Run(); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	toolsCmd.AddCommand(testDockerIgnoreCmd)
}
