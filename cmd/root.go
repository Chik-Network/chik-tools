package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/chik-network/go-modules/pkg/slogs"
)

var (
	gitVersion string
	buildTime  string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:     "chik-tools",
	Short:   "Collection of CLI tools for working with Chik Blockchain",
	Version: fmt.Sprintf("%s (%s)", gitVersion, buildTime),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(InitLogs)

	RootCmd.PersistentFlags().String("log-level", "info", "The log-level for the application, can be one of info, warn, error, debug.")
	RootCmd.PersistentFlags().Bool("dry-run", false, "Show what changes would be made without actually making them. For commands that modify data or configuration, this will show the old and new values.")

	cobra.CheckErr(viper.BindPFlag("log-level", RootCmd.PersistentFlags().Lookup("log-level")))
	cobra.CheckErr(viper.BindPFlag("dry-run", RootCmd.PersistentFlags().Lookup("dry-run")))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	// Search config in home directory with name ".chik-tools" (without extension).
	viper.AddConfigPath(home)
	viper.SetConfigType("yaml")
	viper.SetConfigName(".chik-tools")

	viper.SetEnvPrefix("CHIK_TOOLS")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

// InitLogs sets up the logger
func InitLogs() {
	slogs.Init(viper.GetString("log-level"))
}
