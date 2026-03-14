package main

import (
	"os"

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
