package cmd

import (
	"bufio"
	"fmt"
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
		if strings.Contains(args[0], "MERGE_MSG") {
			return nil
		}

		commitMsgBytes, readErr := os.ReadFile(args[0])
		if readErr != nil {
			return readErr
		}

		firstLine, _, _ := strings.Cut(string(commitMsgBytes[:]), "\n")

		jiraKeyRegexp := regexp.MustCompile("[a-zA-Z]{2,}-[0-9]{1,}")

		if jiraKeyRegexp.MatchString(firstLine) {
			return nil
		}

		key := strings.ToUpper(jiraKeyRegexp.FindString(project.Branch()))
		if key == "" {
			return fmt.Errorf("aborting commit, missing JIRA key in: %s", firstLine)
		}

		commitMsgFile, createErr := os.Create(args[0])
		if createErr != nil {
			return createErr
		}
		defer commitMsgFile.Close()

		w := bufio.NewWriter(commitMsgFile)
		w.WriteString(key+" ")
		w.Write(commitMsgBytes)

		return w.Flush()
	},
}

func init() {
	runHookCmd.AddCommand(runCommitMsgHookCmd)
}
