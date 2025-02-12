package config

import (
	"encoding/hex"
	"net"
	"path"
	"strconv"

	"github.com/chik-network/go-chik-libs/pkg/config"
	"github.com/chik-network/go-chik-libs/pkg/peerprotocol"
	"github.com/chik-network/go-modules/pkg/slogs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/chik-network/chik-tools/internal/utils"
)

var (
	skipConfirm bool
)

// addTrustedPeerCmd Adds a trusted peer to the config
var addTrustedPeerCmd = &cobra.Command{
	Use:   "add-trusted-peer",
	Short: "Adds a trusted peer to the config file",
	Example: `chik-tools config add-trusted-peer 1.2.3.4

# The following version will also override the port to use when connecting to this peer
chik-tools config add-trusted-peer 1.2.3.4 19678`,
	Run: func(cmd *cobra.Command, args []string) {
		chikRoot, err := config.GetChikRootPath()
		if err != nil {
			slogs.Logr.Fatal("Unable to determine CHIK_ROOT", "error", err)
		}

		// 1: Peer IP
		// 2: Optional, port
		if len(args) < 1 || len(args) > 2 {
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

		peer := args[0]
		port := cfg.FullNode.Port
		if len(args) > 1 {
			port64, err := strconv.ParseUint(args[1], 10, 16)
			if err != nil {
				slogs.Logr.Fatal("Invalid port provided")
			}
			port = uint16(port64)
		}

		ip := net.ParseIP(peer)
		if ip == nil {
			slogs.Logr.Fatal("Invalid IP address", "id", peer)
		}
		slogs.Logr.Info("Attempting to get peer id", "peer", peer, "port", port)

		keypair, err := cfg.FullNode.SSL.LoadPublicKeyPair(chikRoot)
		if err != nil {
			slogs.Logr.Fatal("Error loading certs from CHIK_ROOT", "CHIK_ROOT", chikRoot, "error", err)
		}
		if keypair == nil {
			slogs.Logr.Fatal("Error loading certs from CHIK_ROOT", "CHIK_ROOT", chikRoot, "error", "keypair was nil")
		}
		conn, err := peerprotocol.NewConnection(
			&ip,
			peerprotocol.WithPeerPort(port),
			peerprotocol.WithNetworkID(*cfg.SelectedNetwork),
			peerprotocol.WithPeerKeyPair(*keypair),
		)
		if err != nil {
			slogs.Logr.Fatal("Error creating connection", "error", err)
		}
		peerID, err := conn.PeerID()
		if err != nil {
			slogs.Logr.Fatal("Error getting peer id", "error", err)
		}
		peerIDStr := hex.EncodeToString(peerID[:])
		slogs.Logr.Info("peer id received", "peer", peerIDStr)
		if !utils.ConfirmAction("Would you like trust this peer? (y/N)", skipConfirm) {
			slogs.Logr.Error("Cancelled")
		}
		cfg.Wallet.TrustedPeers[peerIDStr] = "Does_not_matter"

		peerToAdd := config.Peer{
			Host: ip.String(),
			Port: port,
		}

		foundPeer := false
		for idx, peer := range cfg.Wallet.FullNodePeers {
			if peer.Host == ip.String() {
				foundPeer = true
				cfg.Wallet.FullNodePeers[idx] = peerToAdd
			}
		}
		if !foundPeer {
			cfg.Wallet.FullNodePeers = append(cfg.Wallet.FullNodePeers, peerToAdd)
		}

		err = cfg.Save()
		if err != nil {
			slogs.Logr.Fatal("error saving config", "error", err)
		}

		slogs.Logr.Info("Added trusted peer. Restart your chik services for the configuration to take effect")
	},
}

func init() {
	addTrustedPeerCmd.Flags().BoolVarP(&skipConfirm, "yes", "y", false, "Skip confirmation")
	configCmd.AddCommand(addTrustedPeerCmd)
}
