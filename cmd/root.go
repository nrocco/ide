package cmd

import (
	"os"
	"path"
	"strings"

	"github.com/nrocco/ide/pkg/ide"
	"github.com/spf13/cobra"
)

var project *ide.Project

func loadProject(cmd *cobra.Command, args []string) error {
	if project != nil {
		return nil
	}

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	project, err = ide.NewProject(dir)
	if err != nil {
		return err
	}

	return nil
}

var rootCmd = &cobra.Command{
	Use:          "ide",
	Short:        "ide provides a powerful ide that gets out of your way",
	SilenceUsage: true,
}

// Execute executes the rootCmd logic and is the main entry point
func Execute() {
	if strings.Contains(os.Args[0], ".git/hooks") {
		args := []string{os.Args[0], "hook", "run", path.Base(os.Args[0]), "--"}
		os.Args = append(args, os.Args[1:]...)
	} else if os.Args[0] == "compare" {
		args := []string{"ide", "tools", "compare", "--"}
		os.Args = append(args, os.Args[1:]...)
	} else if os.Args[0] == "consume" {
		args := []string{"ide", "tools", "consume", "--"}
		os.Args = append(args, os.Args[1:]...)
	} else if os.Args[0] == "rgit" {
		args := []string{"ide", "tools", "rgit", "--"}
		os.Args = append(args, os.Args[1:]...)
	} else if strings.Contains(os.Args[0], "/go-build") {
		// TODO this exists for development only
	} else if !strings.Contains(os.Args[0], "ide") {
		args := []string{os.Args[0], "shim", "run", "--", path.Base(os.Args[0])}
		os.Args = append(args, os.Args[1:]...)
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
