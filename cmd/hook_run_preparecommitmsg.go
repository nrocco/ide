package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var runPrepareCommitMsgHookCmd = &cobra.Command{
	Use:   "prepare-commit-msg",
	Short: "Run the prepare-commit-msg hook for the ide project",
	Long:  "Run the prepare-commit-msg hook for the ide project",
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	PreRunE: loadProject,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return nil
		}

		commitMsgBytes, readErr := ioutil.ReadFile(args[0])
		if readErr != nil {
			return readErr
		}

		commitMsgFile, createErr := os.Create(args[0])
		if createErr != nil {
			return createErr
		}

		defer commitMsgFile.Close()

		w := bufio.NewWriter(commitMsgFile)

		re := regexp.MustCompile("[a-zA-Z]+-[0-9]+")

		key := re.FindString(project.Branch())
		if key != "" {
			fmt.Fprintln(w, strings.ToUpper(key)+" ")
		}

		w.Write(commitMsgBytes)

		return w.Flush()
	},
}

func init() {
	runHookCmd.AddCommand(runPrepareCommitMsgHookCmd)
}
