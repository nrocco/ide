package main

import (
	"fmt"

	"github.com/nrocco/ide/pkg/ide/consume"
	"github.com/spf13/cobra"
)

var consumeCmd = &cobra.Command{
	Use:   "consume [host] as [session] [httpie args...]",
	Short: "HTTP client wrapper using httpie sessions",
	Long:  "HTTP client wrapper that uses httpie session files with optional op:// value resolving",
	Args:  cobra.MinimumNArgs(4),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		switch len(args) {
		case 0:
			hosts, err := consume.ListHosts()
			if err != nil {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return hosts, cobra.ShellCompDirectiveNoFileComp
		case 1:
			return []string{"as"}, cobra.ShellCompDirectiveNoFileComp
		case 2:
			sessions, err := consume.ListSessionsForHost(args[0])
			if err != nil {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return sessions, cobra.ShellCompDirectiveNoFileComp
		default:
			return nil, cobra.ShellCompDirectiveDefault
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if args[1] != "as" {
			return fmt.Errorf("expected 'as' keyword, got '%s'", args[1])
		}
		return consume.Run(args[0], args[2], args[3:])
	},
}

func init() {
	toolCmd.AddCommand(consumeCmd)
}
