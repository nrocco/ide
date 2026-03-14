package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/nrocco/ide/pkg/ide/tools"
	"github.com/spf13/cobra"
)

var rgitCmd = &cobra.Command{
	Use:                "rgit",
	Short:              "Run a git command in multiple git projects",
	Long:               "Run a git command in multiple git projects",
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		allOptionsParsed := false
		options := []string{}
		repositories := []string{}

		rgitDryRun := false
		rgitRecursive := false
		rgitSleep := false
		rgitVerbose := false

		// Parse all arguments from the command line, splitting them in repositories and options for git
		for _, arg := range args {
			if arg == "--" {
				allOptionsParsed = true
			} else if arg == "--rgit-recursive" {
				rgitRecursive = true
			} else if arg == "--rgit-dry-run" {
				rgitDryRun = true
			} else if arg == "--rgit-verbose" {
				rgitVerbose = true
			} else if arg == "--rgit-sleep" {
				rgitSleep = true
			} else if allOptionsParsed {
				repositories = append(repositories, arg)
			} else {
				options = append(options, arg)
			}
		}

		if rgitVerbose {
			fmt.Printf("os.Args:      %v\n", os.Args)
			fmt.Printf("args:         %v\n", args)
			fmt.Printf("options:      %v\n", options)
			fmt.Printf("repositories: %v\n", repositories)
		}

		if len(options) == 0 {
			return errors.New("no git options given")
		}

		if len(repositories) == 0 {
			repositories, _ = filepath.Glob("*")
		}

		options = append([]string{"--no-pager"}, options...)

		// For every repository execute the git command
		return tools.WalkGitRepositories(rgitRecursive, func(repo string) error {
			fmt.Println("==> Working on " + color.GreenString(repo))
			if rgitVerbose {
				fmt.Println(options)
			}
			if rgitDryRun {
				return nil
			}
			// fmt.Println(options)
			cmd := exec.Command("git", options...)
			cmd.Dir = repo
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				fmt.Println(err)
			}
			if rgitSleep {
				time.Sleep(time.Duration(2) * time.Second)
			}
			return nil
		}, repositories...)
	},
}

func init() {
	toolCmd.AddCommand(rgitCmd)
}
