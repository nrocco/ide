package cmd

import (
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/nrocco/ide/pkg/ide"
)

var cfgFile string

var project *ide.Project

func loadProject(cmd *cobra.Command, args []string) error {
	if project != nil {
		return nil
	}

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	project, err = ide.LoadProject(dir)
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
		name := path.Base(os.Args[0])
		args := []string{os.Args[0], "hook", "run", name, "--"}
		os.Args = append(args, os.Args[1:]...)
	} else if !strings.Contains(os.Args[0], "ide") {
		name := path.Base(os.Args[0])
		args := []string{os.Args[0], "binary", "run", "--", name}
		os.Args = append(args, os.Args[1:]...)
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .ide.yaml in $PWD, $HOME, /etc)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName(".ide")
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME")
		viper.AddConfigPath("/etc/")
	}

	viper.SetEnvPrefix("ide")
	viper.AutomaticEnv()

	viper.ReadInConfig()
}
