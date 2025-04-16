package datalayer

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/chik-network/go-chik-libs/pkg/rpc"
	"github.com/chik-network/go-modules/pkg/slogs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// bulkSubCmd Subscribes to multiple datastores at once using the output of chik data get_subscriptions
var bulkSubCmd = &cobra.Command{
	Use:   "bulk-subscribe",
	Short: "Subscribes to multiple datastores at once using the output of chik data get_subscriptions",
	Example: `chik-tools data bulk-subscribe -f subscriptions.json

# Show what changes would be made without actually subscribing
chik-tools data bulk-subscribe -f subscriptions.json --dry-run`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := rpc.NewClient(rpc.ConnectionModeHTTP, rpc.WithAutoConfig())
		if err != nil {
			slogs.Logr.Fatal("error creating chik RPC client", "error", err)
		}

		var content []byte
		if len(args) != 0 {
			content = []byte(strings.Join(args, " "))
		} else {
			jsonFile := viper.GetString("bulksub-file")
			content, err = os.ReadFile(jsonFile)
			if err != nil {
				slogs.Logr.Fatal("Unable to read input file", "error", err)
			}
		}

		subs := &rpc.DatalayerGetSubscriptionsResponse{}
		err = json.Unmarshal(content, subs)
		if err != nil {
			slogs.Logr.Fatal("Could not parse the subscriptions json file", "error", err)
		}

		dryRun := viper.GetBool("dry-run")
		if dryRun {
			slogs.Logr.Info("DRY RUN: Would subscribe to the following stores")
		}

		for _, id := range subs.StoreIDs {
			if dryRun {
				slogs.Logr.Info("DRY RUN: Would subscribe to store", "id", id)
				continue
			}

			slogs.Logr.Info("Subscription to store", "id", id)
			_, _, err = client.DataLayerService.Subscribe(&rpc.DatalayerSubscribeOptions{
				ID: id,
			})
			if err != nil {
				slogs.Logr.Error("Error subscribing to datastore", "id", id, "error", err)
			}
		}

		if dryRun {
			slogs.Logr.Info("DRY RUN: No changes were made to subscriptions")
		}
	},
}

func init() {
	bulkSubCmd.PersistentFlags().StringP("file", "f", "", "The file containing the json of subscriptions to add")

	cobra.CheckErr(viper.BindPFlag("bulksub-file", bulkSubCmd.PersistentFlags().Lookup("file")))

	datalayerCmd.AddCommand(bulkSubCmd)
}
