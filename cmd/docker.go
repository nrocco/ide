package cmd

import (
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

var dockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "Blaat",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		cli, err := client.NewEnvClient()
		if err != nil {
			return err
		}

		// _, err = cli.ImagePull(ctx, "docker.io/library/alpine", types.ImagePullOptions{})
		// if err != nil {
		// 	panic(err)
		// }

		resp, err := cli.ContainerCreate(ctx, &container.Config{
			Image: "alpine",
			Cmd:   []string{"echo", "hello world"},
		}, nil, nil, "")
		if err != nil {
			return err
		}

		if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
			return err
		}

		if _, err = cli.ContainerWait(ctx, resp.ID); err != nil {
			return err
		}

		out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
		if err != nil {
			return err
		}

		io.Copy(os.Stdout, out)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(dockerCmd)
}
