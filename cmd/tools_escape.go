package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var escapeCmd = &cobra.Command{
	Use:   "escape [file]",
	Short: "Escape newlines in input",
	Long:  "Read from stdin or a file and convert all newlines to literal \\n",
	Args:  cobra.MaximumNArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveDefault
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var input io.Reader

		if len(args) == 0 {
			input = os.Stdin
		} else {
			file, err := os.Open(args[0])
			if err != nil {
				return err
			}
			defer file.Close()
			input = file
		}

		scanner := bufio.NewScanner(input)

		for scanner.Scan() {
			fmt.Print(scanner.Text())
			fmt.Print("\\n")
		}

		fmt.Print("\n")

		return scanner.Err()
	},
}

func init() {
	toolsCmd.AddCommand(escapeCmd)
}
