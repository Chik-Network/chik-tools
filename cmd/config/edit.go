package config

import (
	"path"

	"github.com/chik-network/go-chik-libs/pkg/config"
	"github.com/chik-network/go-modules/pkg/slogs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// editCmd generates a new chik config
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit an existing chik configuration file",
	Example: `chik-tools config edit --config ~/.chik/mainnet/config/config.yaml --set full_node.port=59678 --set full_node.target_peer_count=10

# The following version will discover the config file by inspecting CHIK_ROOT or using the default CHIK_ROOT
chik-tools config edit --set full_node.port=59678 --set full_node.target_peer_count=10

# Show what changes would be made without actually making them
chik-tools config edit --set full_node.port=59678 --dry-run`,
	Run: func(cmd *cobra.Command, args []string) {
		chikRoot, err := config.GetChikRootPath()
		if err != nil {
			slogs.Logr.Fatal("Unable to determine CHIK_ROOT", "error", err)
		}

		if len(args) > 0 {
			slogs.Logr.Fatal("Unexpected number of arguments provided")
		}

		cfgPath := viper.GetString("config")
		if cfgPath == "" {
			// Use default chik root
			cfgPath = path.Join(chikRoot, "config", "config.yaml")
		}

		cfg, err := config.LoadConfigAtRoot(cfgPath, chikRoot)
		if err != nil {
			slogs.Logr.Fatal("error loading chik config", "error", err)
		}

		err = cfg.FillValuesFromEnvironment()
		if err != nil {
			slogs.Logr.Fatal("error filling values from environment", "error", err)
		}

		dryRun := viper.GetBool("dry-run")
		if dryRun {
			slogs.Logr.Info("DRY RUN: The following changes would be made to the config file")
		}

		valuesToSet := viper.GetStringMapString("edit-set")
		for path, value := range valuesToSet {
			pathMap := config.ParsePathsFromStrings([]string{path}, false)
			var key string
			var pathSlice []string
			for key, pathSlice = range pathMap {
				break
			}

			// Get the current value using GetFieldByPath
			currentValue, err := cfg.GetFieldByPath(pathSlice)
			if err != nil {
				slogs.Logr.Info("Config value not found", "path", path)
			}

			if dryRun {
				slogs.Logr.Info("Would change config value",
					"path", path,
					"current_value", currentValue,
					"new_value", value)
				continue
			}

			err = cfg.SetFieldByPath(pathSlice, value)
			if err != nil {
				slogs.Logr.Fatal("error setting path in config", "key", key, "value", value, "error", err)
			}
		}

		if dryRun {
			slogs.Logr.Info("DRY RUN: No changes were made to the config file")
			return
		}

		err = cfg.Save()
		if err != nil {
			slogs.Logr.Fatal("error saving config", "error", err)
		}
	},
}

func init() {
	editCmd.PersistentFlags().StringToStringP("set", "s", nil, "Paths and values to set in the config")

	cobra.CheckErr(viper.BindPFlag("edit-set", editCmd.PersistentFlags().Lookup("set")))

	configCmd.AddCommand(editCmd)
}
