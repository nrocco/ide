package main

import (
	"errors"
	"fmt"

	"github.com/nrocco/ide/pkg/ide"
	"github.com/spf13/cobra"
)

var bumpCmdOpts struct {
	Strategy string
	Force    bool
}

var bumpCmd = &cobra.Command{
	Use:       "bump {major|minor|patch}",
	Short:     "Bump the semantic version of your ide project",
	Long:      "Bump the semantic version of your ide project. Without --force only the new version is shown.",
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	ValidArgs: []string{"major", "minor", "patch"},
	PreRunE:   loadProject,
	RunE: func(cmd *cobra.Command, args []string) error {
		strategy, err := project.VersionStrategy(bumpCmdOpts.Strategy)
		if err != nil {
			return err
		}

		if bumpCmdOpts.Strategy == ide.VersionStrategyAuto {
			fmt.Fprintf(cmd.OutOrStdout(), "using strategy %s\n", strategy.Name())
		}

		current, err := strategy.Current()
		latest := current.String()
		if errors.Is(err, ide.ErrNoVersion) {
			latest = "<none>"
			fmt.Fprintln(cmd.OutOrStdout(), "no existing version found, starting from 0.0.0")
		} else if err != nil {
			return err
		}

		next, err := current.Bump(args[0])
		if err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "%s -> %s\n", latest, next)

		if !bumpCmdOpts.Force {
			return nil
		}

		return strategy.Write(next)
	},
}

func init() {
	bumpCmd.Flags().StringVar(&bumpCmdOpts.Strategy, "strategy", ide.VersionStrategyAuto, "strategy used to read and write the version, auto detects the strategy")
	bumpCmd.Flags().BoolVar(&bumpCmdOpts.Force, "force", false, "actually write the new version instead of only showing it")

	bumpCmd.RegisterFlagCompletionFunc("strategy", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		names := []string{ide.VersionStrategyAuto}
		if err := loadProject(cmd, args); err == nil {
			for _, strategy := range project.VersionStrategies() {
				names = append(names, strategy.Name())
			}
		}
		return names, cobra.ShellCompDirectiveNoFileComp
	})

	rootCmd.AddCommand(bumpCmd)
}
