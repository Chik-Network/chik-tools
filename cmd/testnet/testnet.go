package testnet

import (
	"github.com/spf13/cobra"

	"github.com/chik-network/chik-tools/cmd"
)

// testnetCmd represents the config command
var testnetCmd = &cobra.Command{
	Use:        "testnet",
	Short:      "Utilities for working with chik testnets",
	Deprecated: "\nThe testnet subcommand is deprecated. Please use the 'network' subcommand instead\n",
}

func init() {
	cmd.RootCmd.AddCommand(testnetCmd)
}
