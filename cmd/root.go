package cmd

import (
	"fmt"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/nrocco/ide/pkg/ide"
)

var cfgFile string

// Project represents an instance of ide.Project
var Project *ide.Project

var rootCmd = &cobra.Command{
	Use:          "ide",
	Short:        "ide provides a powerful ide that gets out of your way",
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		Project, err = ide.LoadProject(dir)
		if err != nil {
			return err
		}

		return nil
	},
}

// Execute executes the rootCmd logic and is the main entry point for ide
func Execute() {
	if strings.Contains(os.Args[0], ".git/hooks") {
		elems := strings.Split(os.Args[0], "/")

		args := []string{os.Args[0], "hook", "run", elems[2]}
		os.Args = append(args, os.Args[1:]...)
	} else if !strings.Contains(os.Args[0], "ide") {
		args := []string{os.Args[0], "exec", "run", "--", os.Args[0]}
		os.Args = append(args, os.Args[1:]...)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ide.yaml)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err == nil {
			viper.AddConfigPath(home)
			viper.SetConfigName(".ide")
		}
	}

	viper.SetEnvPrefix("ide")
	viper.AutomaticEnv()

	viper.ReadInConfig()
}
