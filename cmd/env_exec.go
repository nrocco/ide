package cmd

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"io"
	"os"
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

		ctx := context.Background()
		cli, err := client.NewEnvClient()
		if err != nil {
			return err
		}

		execOpts := types.ExecConfig{
			Cmd:          args,
			Tty:          true,
			AttachStdin:  true,
			AttachStdout: true,
			AttachStderr: true,
		}

		execInstance, err := cli.ContainerExecCreate(ctx, container, execOpts)
		if err != nil {
			return err
		}

		att, err := cli.ContainerExecAttach(ctx, execInstance.ID, execOpts)
		if err != nil {
			return err
		}

		execStartOpts := types.ExecStartCheck{}
		err = cli.ContainerExecStart(ctx, execInstance.ID, execStartOpts)
		if err != nil {
			return err
		}

		io.Copy(os.Stdout, att.Reader)

		return nil
	},
}

func init() {
	envCmd.AddCommand(execEnvCmd)
}
