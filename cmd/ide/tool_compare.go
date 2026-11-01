package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var compareCmdOpts struct {
	Recursive bool
	With      string
	First     int
	Last      int
	DryRun    bool
}

var compareCmd = &cobra.Command{
	Use:   "compare",
	Short: "Compare a file in the current project with another project",
	Long:  "Compare a file in the current project with another project",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if compareCmdOpts.First > 8 || compareCmdOpts.Last > 8 {
			return fmt.Errorf("--first and --last can not be higher than 8")
		}

		file := args[len(args)-1]
		projects := args[0 : len(args)-1]

		if compareCmdOpts.Recursive {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			var found []string
			for _, dir := range projects {
				entries, err := os.ReadDir(dir)
				if err != nil {
					return err
				}

				for _, entry := range entries {
					if !entry.IsDir() {
						continue
					}

					project := filepath.Join(dir, entry.Name())

					// Skip the current project, it is always part of the comparison
					if abs, err := filepath.Abs(project); err == nil && abs == cwd {
						continue
					}

					if _, err := os.Stat(filepath.Join(project, compareCmdOpts.With)); err != nil {
						continue
					}

					found = append(found, project)
				}
			}

			if compareCmdOpts.First > 0 && len(found) > compareCmdOpts.First {
				found = found[:compareCmdOpts.First]
			}

			if compareCmdOpts.Last > 0 && len(found) > compareCmdOpts.Last {
				found = found[len(found)-compareCmdOpts.Last:]
			}

			projects = found
		}

		files := []string{file}

		for _, project := range projects {
			files = append(files, filepath.Join(project, file))
		}

		if compareCmdOpts.DryRun {
			for _, f := range files {
				fmt.Println(f)
			}
			return nil
		}

		options := append([]string{"-d"}, files...)

		command := exec.Command("nvim", options...) // #nosec G204 -- fixed binary, arguments are passed without a shell
		command.Stdin = os.Stdin
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		return command.Run()
	},
}

func init() {
	compareCmd.Flags().BoolVar(&compareCmdOpts.Recursive, "recursive", false, "compare with all projects found in the given directories")
	compareCmd.Flags().StringVar(&compareCmdOpts.With, "with", "", "only include projects that contain this file, relative to the project directory")
	compareCmd.Flags().IntVar(&compareCmdOpts.First, "first", 0, "only include the first n projects found, 0 means all")
	compareCmd.Flags().IntVar(&compareCmdOpts.Last, "last", 0, "only include the last n projects found, 0 means all")
	compareCmd.MarkFlagsMutuallyExclusive("first", "last")
	compareCmd.Flags().BoolVar(&compareCmdOpts.DryRun, "dry-run", false, "only print the files that would be compared")
	toolCmd.AddCommand(compareCmd)
}
