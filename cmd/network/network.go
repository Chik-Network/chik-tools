package network

import (
	"github.com/spf13/cobra"

	"github.com/chik-network/chik-tools/cmd"
)

// networkCmd represents the config command
var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Utilities for working with chik networks",
}

func init() {
	cmd.RootCmd.AddCommand(networkCmd)
}
