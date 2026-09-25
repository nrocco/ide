package main

import (
	"io"
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

		// Tee s_client's stdout and stderr to the terminal and into x509's stdin
		pr, pw, err := os.Pipe()
		if err != nil {
			return err
		}

		sClient := exec.Command("openssl", "s_client", "-no-interactive", "-servername", host, "-showcerts", host+":"+port)
		sClient.Stdout = &teeWriter{w: os.Stdout, pipe: pw}
		sClient.Stderr = &teeWriter{w: os.Stderr, pipe: pw}

		x509Args := append([]string{"x509", "-in", "/dev/stdin"}, opensslArgs...)
		x509 := exec.Command("openssl", x509Args...)
		x509.Stdin = pr
		x509.Stdout = os.Stdout
		x509.Stderr = os.Stderr

		if err := x509.Start(); err != nil {
			pr.Close()
			pw.Close()
			return err
		}
		pr.Close()

		if err := sClient.Start(); err != nil {
			pw.Close()
			x509.Wait()
			return err
		}

		sClient.Wait()
		pw.Close()
		return x509.Wait()
	},
}

// teeWriter writes to w and, best effort, to pipe. Errors writing to pipe are
// ignored so output keeps flowing to w after the reading process has exited.
type teeWriter struct {
	w    io.Writer
	pipe io.Writer
}

func (t *teeWriter) Write(p []byte) (int, error) {
	t.pipe.Write(p)
	return t.w.Write(p)
}

func init() {
	toolCmd.AddCommand(checkcertCmd)
}
