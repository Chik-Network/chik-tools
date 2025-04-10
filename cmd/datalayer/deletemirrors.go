package datalayer

import (
	"fmt"

	"github.com/chik-network/go-chik-libs/pkg/rpc"
	"github.com/chik-network/go-chik-libs/pkg/types"
	"github.com/chik-network/go-modules/pkg/slogs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// deleteMirrorsCmd Deletes all owned mirrors for all datalayer subscriptions
var deleteMirrorsCmd = &cobra.Command{
	Use:     "delete-mirrors",
	Short:   "Deletes all owned mirrors for all datalayer subscriptions",
	Example: "chik-tools data delete-mirrors --all\nchik-tools data delete-mirrors --id abcd1234\nchik-tools data delete-mirrors --all --dry-run",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		all := viper.GetBool("delete-mirror-all")
		subID := viper.GetString("delete-mirror-id")
		if !all && subID == "" {
			return fmt.Errorf("must provide a subscription ID with --id flag or use --all option")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		client, err := rpc.NewClient(rpc.ConnectionModeHTTP, rpc.WithAutoConfig())
		if err != nil {
			slogs.Logr.Fatal("error creating chik RPC client", "error", err)
		}

		// Figure out what fee we are using
		feeXCK := viper.GetFloat64("delete-mirror-fee")
		feeMojos := uint64(feeXCK * 1000000000000)
		slogs.Logr.Info("fee for all transactions", "xck", feeXCK, "mojos", feeMojos)

		all := viper.GetBool("delete-mirror-all")
		subID := viper.GetString("delete-mirror-id")
		dryRun := viper.GetBool("delete-mirror-dry-run")

		if dryRun {
			slogs.Logr.Info("DRY RUN: No mirrors will be deleted")
		}

		if all {
			slogs.Logr.Info("deleting all owned mirrors for all subscriptions")
			subscriptions, _, err := client.DataLayerService.GetSubscriptions(&rpc.DatalayerGetSubscriptionsOptions{})
			if err != nil {
				slogs.Logr.Fatal("error getting list of datalayer subscriptions", "error", err)
			}

			for _, subscription := range subscriptions.StoreIDs {
				deleteMirrorsForSubscription(client, subscription, feeMojos, dryRun)
			}
		} else {
			deleteMirrorsForSubscription(client, subID, feeMojos, dryRun)
		}
	},
}

func deleteMirrorsForSubscription(client *rpc.Client, subscription string, feeMojos uint64, dryRun bool) {
	slogs.Logr.Info("checking subscription", "store", subscription)

	mirrors, _, err := client.DataLayerService.GetMirrors(&rpc.DatalayerGetMirrorsOptions{
		ID: subscription,
	})
	if err != nil {
		slogs.Logr.Fatal("error fetching mirrors for subscription", "store", subscription, "error", err)
	}
	var ownedMirrors []types.Bytes32

	for _, mirror := range mirrors.Mirrors {
		if mirror.Ours {
			ownedMirrors = append(ownedMirrors, mirror.CoinID)
		}
	}

	if len(ownedMirrors) == 0 {
		slogs.Logr.Info("no owned mirrors for this datastore", "store", subscription)
		return
	}

	for _, coinID := range ownedMirrors {
		if dryRun {
			slogs.Logr.Info("DRY RUN: Would delete mirror", "store", subscription, "mirror", coinID.String())
			continue
		}

		slogs.Logr.Info("deleting mirror", "store", subscription, "mirror", coinID.String())
		resp, _, err := client.DataLayerService.DeleteMirror(&rpc.DatalayerDeleteMirrorOptions{
			CoinID: coinID.String(),
			Fee:    feeMojos,
		})
		if err != nil {
			slogs.Logr.Fatal("error deleting mirror for store", "store", subscription, "mirror", coinID, "error", err)
		}
		if !resp.Success {
			slogs.Logr.Fatal("unknown error when deleting mirror for store", "store", subscription, "mirror", coinID)
		}
	}
}

func init() {
	deleteMirrorsCmd.PersistentFlags().Float64P("fee", "m", 0, "Fee to use when deleting the mirrors. The fee is used per mirror. Units are XCK")
	deleteMirrorsCmd.PersistentFlags().Bool("all", false, "Delete all owned mirrors for all subscriptions")
	deleteMirrorsCmd.PersistentFlags().String("id", "", "The subscription ID to delete mirrors for")
	deleteMirrorsCmd.PersistentFlags().Bool("dry-run", false, "Show what mirrors would be deleted without actually deleting them")

	cobra.CheckErr(viper.BindPFlag("delete-mirror-fee", deleteMirrorsCmd.PersistentFlags().Lookup("fee")))
	cobra.CheckErr(viper.BindPFlag("delete-mirror-all", deleteMirrorsCmd.PersistentFlags().Lookup("all")))
	cobra.CheckErr(viper.BindPFlag("delete-mirror-id", deleteMirrorsCmd.PersistentFlags().Lookup("id")))
	cobra.CheckErr(viper.BindPFlag("delete-mirror-dry-run", deleteMirrorsCmd.PersistentFlags().Lookup("dry-run")))

	datalayerCmd.AddCommand(deleteMirrorsCmd)
}
