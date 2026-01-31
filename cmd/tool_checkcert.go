package cmd

import (
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var checkcertCmd = &cobra.Command{
	Use:                "checkcert host[:port] [openssl args...]",
	Short:              "Check TLS certificate for a host",
	Long:               "Connect to a host and display its TLS certificate using openssl",
	Args:               cobra.MinimumNArgs(1),
	DisableFlagParsing: true,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var host, port string

		if strings.Contains(args[0], ":") {
			parts := strings.SplitN(args[0], ":", 2)
			host = parts[0]
			port = parts[1]
		} else {
			host = args[0]
			port = "443"
		}

		opensslArgs := []string{"-text"}
		if len(args) > 1 {
			opensslArgs = args[1:]
		}

		sClient := exec.Command("openssl", "s_client", "-no-interactive", "-servername", host, "-showcerts", host+":"+port)
		sClient.Stderr = nil

		x509Args := append([]string{"x509", "-in", "/dev/stdin", "-noout"}, opensslArgs...)
		x509 := exec.Command("openssl", x509Args...)
		x509.Stdout = os.Stdout
		x509.Stderr = os.Stderr

		pipe, err := sClient.StdoutPipe()
		if err != nil {
			return err
		}
		x509.Stdin = pipe

		if err := sClient.Start(); err != nil {
			return err
		}
		if err := x509.Start(); err != nil {
			return err
		}

		sClient.Wait()
		return x509.Wait()
	},
}

func init() {
	toolCmd.AddCommand(checkcertCmd)
}
