package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/chik-network/chik-tools/cmd"
)

var (
	skipConfirm bool
	retries     uint
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Utilities for working with chik config",
}

func init() {
	configCmd.PersistentFlags().String("config", "", "existing config file to use (default is to look in $CHIK_ROOT)")
	cobra.CheckErr(viper.BindPFlag("config", configCmd.PersistentFlags().Lookup("config")))

	cmd.RootCmd.AddCommand(configCmd)
}
