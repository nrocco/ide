package cmd

import (
	"errors"
	"log"
	"net"

	"github.com/nrocco/ide/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run an ide server for processing long running operations",
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		port := viper.GetString("server.address")
		if port == "" {
			return errors.New("No server.address specified in config file")
		}

		lis, err := net.Listen("tcp", port)
		if err != nil {
			return err
		}

		log.Printf("Running ide server on port %s", port)

		s := server.NewServer()
		return s.Serve(lis)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
