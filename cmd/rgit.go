package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var rgitCmd = &cobra.Command{
	Use:   "rgit",
	Short: "Run a git command in multiple git projects",
	Long:  "Run a git command in multiple git projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		allOptionsParsed := false
		options := []string{}
		repositories := []string{}

		// Parse all arguments from the command line, splitting them in repositories and options for git
		for _, arg := range args {
			if arg == "--" {
				allOptionsParsed = true
			} else if allOptionsParsed {
				repositories = append(repositories, arg)
			} else {
				options = append(options, arg)
			}
		}

		if len(options) == 0 {
			return errors.New("no git options given")
		}

		if len(repositories) == 0 {
			repositories, _ = filepath.Glob("*")
		}

		options = append([]string{"--no-pager"}, options...)

		// For every repository execute the git command
		for _, repo := range repositories {
			if repoStat, err := os.Stat(repo); err == nil {
				if !repoStat.IsDir() {
					continue
				}
			}

			if _, err := os.Stat(filepath.Join(repo, ".git")); err != nil {
				fmt.Fprintf(os.Stderr, "[ERROR] Skipping %s as it is not a git repository\n", repo)
				continue
			}
			fmt.Printf("==> Working on \033[0;32m%s\033[0m\n", repo)
			cmd := exec.Command("git", options...)
			cmd.Dir = repo
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()

			// TODO support RGIT_SLEEP
		}
		return nil
	},
}

func init() {
	toolsCmd.AddCommand(rgitCmd)
}
