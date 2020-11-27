package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"github.com/nrocco/ide/server"
)

var runPostMergeHookCmd = &cobra.Command{
	Use:   "post-merge",
	Short: "Run the post-merge hook for the ide project",
	Long:  "Run the post-merge hook for the ide project",
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	PreRunE: loadProject,
	RunE: func(cmd *cobra.Command, args []string) error {
		port := viper.GetString("server.address")

		if port == "" {
			return project.RefreshCtags()
		}

		conn, err := grpc.Dial(port, grpc.WithInsecure())
		if err != nil {
			return err
		}
		defer conn.Close()

		c := server.NewServerClient(conn)

		r, err := c.RefreshCtags(context.Background(), &server.RefreshCtagsRequest{Directory: project.Location()})
		if err != nil {
			return err
		}

		log.Printf("Writing ctags file to %s\n", r.File)

		return nil
	},
}

func init() {
	runHookCmd.AddCommand(runPostMergeHookCmd)
}
