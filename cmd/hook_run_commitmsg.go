package cmd

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var runCommitMsgHookCmd = &cobra.Command{
	Use:   "commit-msg",
	Short: "Run the commit-msg hook for the ide project",
	Long:  "Run the commit-msg hook for the ide project",
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	Args:    cobra.ExactArgs(1),
	PreRunE: loadProject,
	RunE: func(cmd *cobra.Command, args []string) error {
		commitMsgFile, osErr := os.Open(args[0])
		if osErr != nil {
			return osErr
		}

		defer commitMsgFile.Close()

		reader := bufio.NewReader(commitMsgFile)

		line, readError := reader.ReadString('\n')
		if readError != nil {
			return readError
		}

		line = strings.Trim(line, "\n")

		matched, regexErr := regexp.MatchString("^[a-zA-Z]{2,}-[0-9]{1,}|^Merge.*[a-zA-Z]{2,}-[0-9]{1,}", line)
		if regexErr != nil {
			return regexErr
		}

		if !matched {
			log.Fatalf("Aborting commit due to invalid commit message: %s\n", line)
		}

		return nil
	},
}

func init() {
	runHookCmd.AddCommand(runCommitMsgHookCmd)
}
