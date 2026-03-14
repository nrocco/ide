package main

import (
	"errors"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var generateTagsCmd = &cobra.Command{
	Use:     "generate",
	Short:   "Generate tags",
	Long:    "Generate tags",
	PreRunE: loadProject,
	RunE: func(cmd *cobra.Command, args []string) error {
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		age := project.CtagsFileAge()
		if !age.IsZero() {
			fdArgs := []string{
				"--changed-after=" + age.UTC().Format("2006-01-02T15:04:05Z07:00"),
				"--type=file",
				"--no-ignore-vcs",
			}
			if !dryRun {
				fdArgs = append(fdArgs, "--quiet")
			}

			fd := exec.Command("fd", fdArgs...)
			if dryRun {
				fd.Stdout = os.Stdout
			}

			if err := fd.Run(); err != nil {
				var exitErr *exec.ExitError
				if errors.As(err, &exitErr) {
					// non-zero exit from fd means no files changed — skip regeneration
					return nil
				}
				return err
			}

			if dryRun {
				return nil
			}
		}

		return project.CtagsGenerate()
	},
}

func init() {
	generateTagsCmd.Flags().Bool("dry-run", false, "Show files that would trigger regeneration without regenerating")

	tagsCmd.AddCommand(generateTagsCmd)
}
