package certs

import (
	"github.com/spf13/cobra"

	"github.com/chik-network/chik-tools/cmd"
)

// certsCmd represents the config command
var certsCmd = &cobra.Command{
	Use:   "certs",
	Short: "Utilities for working with chik certificates",
}

func init() {
	cmd.RootCmd.AddCommand(certsCmd)
}
