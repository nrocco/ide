package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

type ctagsEntry struct {
	Type      string `json:"_type"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	Pattern   string `json:"pattern"`
	Line      int    `json:"line"`
	Kind      string `json:"kind"`
	Scope     string `json:"scope"`
	Signature string `json:"signature"`
	Typeref   string `json:"typeref"`
	ScopeKind string `json:"scopeKind"`
	Language  string `json:"language"`
	Roles     string `json:"roles"`
	End       int    `json:"end"`
}

var parseCmd = &cobra.Command{
	Use:   "parse <path>",
	Short: "Parse ctags output for a path",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		ctags := exec.Command("ctags",
			"-f-",
			"--excmd=number",
			"--recurse=yes",
			"--sort=no",
			"--totals=no",
			"--with-list-header=no",
			"--machinable=yes",
			"--kinds-php=f",
			"--kinds-go=f",
			"--kinds-python=cfm",
			"--fields=aCeEfFikKlmnNpPrRsStxzZ",
			"--output-format=json",
			path,
		)

		stdout, err := ctags.StdoutPipe()
		if err != nil {
			return err
		}

		if err := ctags.Start(); err != nil {
			return err
		}

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			var entry ctagsEntry
			if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
				return err
			}
			entry.Language = strings.ToLower(entry.Language)
			entry.Typeref = strings.TrimPrefix(entry.Typeref, "typename:")

			if cmd.Flags().GetBool("debug") {
				fmt.Printf("%+v\n", entry)
				continue
			}

			fmt.Println("we must do smth now")
		}

		if err := scanner.Err(); err != nil {
			return err
		}

		return ctags.Wait()
	},
}

func init() {
	parseCmd.Flags().Bool("debug", false, "Output raw JSON instead of formatted output")
	parseCmd.Flags().Bool("no-public", false, "Exclude public functions") // TODO use this
	parseCmd.Flags().Bool("no-private", false, "Exclude private functions") // TODO use this

	toolCmd.AddCommand(parseCmd)
}
