package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/nrocco/ide/pkg/ide"
	"github.com/spf13/cobra"
)

var parseCmd = &cobra.Command{
	Use:   "parse <path>...",
	Short: "Parse ctags output for a path",
	Args:  cobra.MinimumNArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveDefault
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		debug, _ := cmd.Flags().GetBool("debug")
		noPublic, _ := cmd.Flags().GetBool("no-public")
		noPrivate, _ := cmd.Flags().GetBool("no-private")
		currentFile := ""

		return project.CtagsParseCode(func(entry ide.CtagsEntry) {
			if noPublic && entry.IsPublic() {
				return
			}
			if noPrivate && entry.IsPrivate() {
				return
			}
			if debug {
				fmt.Printf("%+v\n", entry)
				return
			}
			if currentFile != entry.Path {
				if currentFile != "" {
					fmt.Printf("\n\n")
				}
				currentFile = entry.Path
				color.Magenta(entry.Path)
				fmt.Printf("\n")
			}
			color.New(color.FgHiBlack).Printf("%5d  ", entry.Line)
			if entry.Language == "python" && entry.Kind == "member" {
				fmt.Printf("%s.", entry.Scope)
			}
			if entry.Language == "go" && entry.ScopeKind == "struct" {
				fmt.Printf("(%s) ", entry.Scope)
			}
			if entry.IsPublic() {
				color.New(color.FgHiBlue).Printf("%s", entry.Name)
			} else {
				color.New(color.FgHiRed).Printf("%s", entry.Name)
			}
			color.New(color.FgHiGreen).Printf("%s", entry.Signature)
			if entry.Typeref != "" {
				color.New(color.FgHiBlack).Printf(": ")
				color.New(color.FgHiBlue).Printf("%s", entry.Typeref)
			}
			fmt.Printf("\n")
		}, args...)
	},
}

func init() {
	parseCmd.Flags().Bool("debug", false, "Output raw JSON instead of formatted output")
	parseCmd.Flags().Bool("no-public", false, "Exclude public functions")   // TODO use this
	parseCmd.Flags().Bool("no-private", false, "Exclude private functions") // TODO use this

	toolCmd.AddCommand(parseCmd)
}
