package network_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/chik-network/go-chik-libs/pkg/config"
	"github.com/chik-network/go-chik-libs/pkg/types"
	"github.com/stretchr/testify/assert"

	"github.com/chik-network/chik-tools/cmd"
	"github.com/chik-network/chik-tools/cmd/network"
)

const testnetwork = "unittestnet"

func setupDefaultConfig(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "chik-root-*")
	assert.NoError(t, err)

	err = os.Setenv("CHIK_ROOT", tempDir)
	assert.NoError(t, err)

	rootPath, err := config.GetChikRootPath()
	assert.NoError(t, err)

	err = os.MkdirAll(filepath.Join(rootPath, "config"), 0755)
	assert.NoError(t, err)

	defaultConfig, err := config.LoadDefaultConfig()
	assert.NoError(t, err)

	defaultConfig.NetworkOverrides.Constants[testnetwork] = config.NetworkConstants{
		AggSigMeAdditionalData:         "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		DifficultyConstantFactor:       types.Uint128From64(10052721566054),
		DifficultyStarting:             30,
		EpochBlocks:                    768,
		GenesisChallenge:               "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		GenesisPreFarmPoolPuzzleHash:   "08296fc227decd043aee855741444538e4cc9a31772c4d1a9e6242d1e777e42a",
		GenesisPreFarmFarmerPuzzleHash: "08296fc227decd043aee855741444538e4cc9a31772c4d1a9e6242d1e777e42a",
		MempoolBlockBuffer:             10,
		MinPlotSize:                    18,
		NetworkType:                    1,
		SubSlotItersStarting:           67108864,
	}
	defaultConfig.NetworkOverrides.Config[testnetwork] = config.NetworkConfig{
		AddressPrefix:       "txck",
		DefaultFullNodePort: 58445,
	}

	configPath := filepath.Join(rootPath, "config", "config.yaml")
	err = defaultConfig.SavePath(configPath)
	assert.NoError(t, err)
}

func TestNetworkSwitch(t *testing.T) {
	cmd.InitLogs()
	setupDefaultConfig(t)
	cfg, err := config.GetChikConfig()
	assert.NoError(t, err)
	assert.Equal(t, "mainnet", *cfg.SelectedNetwork)

	network.SwitchNetwork("unittestnet", false)

	// reload config from disk
	cfg, err = config.GetChikConfig()
	assert.NoError(t, err)

	port := uint16(58445)
	localpeer := config.Peer{
		Host: "localhost",
		Port: port,
	}
	assert.Equal(t, "unittestnet", *cfg.SelectedNetwork)
	assert.Equal(t, []config.Peer{localpeer}, cfg.Farmer.FullNodePeers)
	assert.Equal(t, "db/blockchain_v2_unittestnet.sqlite", cfg.FullNode.DatabasePath)
	assert.Equal(t, []string{"dns-introducer-unittestnet.chiknetwork.com"}, cfg.FullNode.DNSServers)
	assert.Equal(t, "db/peers-unittestnet.dat", cfg.FullNode.PeersFilePath)
	assert.Equal(t, port, cfg.FullNode.Port)
	assert.Equal(t, config.Peer{Host: "introducer-unittestnet.chiknetwork.com", Port: port}, cfg.FullNode.IntroducerPeer)
	assert.Equal(t, port, cfg.Introducer.Port)
	assert.Equal(t, port, cfg.Seeder.OtherPeersPort)
	assert.Equal(t, []string{"node-unittestnet.chiknetwork.com"}, cfg.Seeder.BootstrapPeers)
	assert.Equal(t, []config.Peer{localpeer}, cfg.Timelord.FullNodePeers)
	assert.Equal(t, []string{"dns-introducer-unittestnet.chiknetwork.com"}, cfg.Wallet.DNSServers)
	assert.Equal(t, []config.Peer{localpeer}, cfg.Wallet.FullNodePeers)
	assert.Equal(t, config.Peer{Host: "introducer-unittestnet.chiknetwork.com", Port: port}, cfg.Wallet.IntroducerPeer)
	assert.Equal(t, "wallet/db/wallet_peers-unittestnet.dat", cfg.Wallet.WalletPeersFilePath)
}

func TestNetworkSwitch_SettingRetention(t *testing.T) {
	cmd.InitLogs()
	setupDefaultConfig(t)
	cfg, err := config.GetChikConfig()
	assert.NoError(t, err)
	assert.Equal(t, "mainnet", *cfg.SelectedNetwork)

	// Set some custom dns introducers, and ensure they are back when swapping away and back to mainnet
	cfg.FullNode.DNSServers = []string{"dns-mainnet-1.example.com", "dns-mainnet-2.example.com"}
	cfg.Seeder.BootstrapPeers = []string{"bootstrap-mainnet-1.example.com"}
	cfg.Seeder.StaticPeers = []string{"static-peer-1.example.com"}
	cfg.FullNode.FullNodePeers = []config.Peer{{Host: "fn-peer-1.example.com", Port: 1234}}
	err = cfg.Save()
	assert.NoError(t, err)

	// reload config from disk to ensure the dns servers were persisted
	cfg, err = config.GetChikConfig()
	assert.NoError(t, err)
	assert.Equal(t, []string{"dns-mainnet-1.example.com", "dns-mainnet-2.example.com"}, cfg.FullNode.DNSServers)

	network.SwitchNetwork("unittestnet", false)
	// reload config from disk to ensure defaults are in the config now
	cfg, err = config.GetChikConfig()
	assert.NoError(t, err)
	assert.Equal(t, []string{"dns-introducer-unittestnet.chiknetwork.com"}, cfg.FullNode.DNSServers)
	assert.Equal(t, []string{"node-unittestnet.chiknetwork.com"}, cfg.Seeder.BootstrapPeers)
	assert.Equal(t, []string{}, cfg.Seeder.StaticPeers)
	assert.Equal(t, []config.Peer{}, cfg.FullNode.FullNodePeers)

	network.SwitchNetwork("mainnet", false)

	// reload config from disk
	cfg, err = config.GetChikConfig()
	assert.NoError(t, err)
	assert.Equal(t, []string{"dns-mainnet-1.example.com", "dns-mainnet-2.example.com"}, cfg.FullNode.DNSServers)
	assert.Equal(t, []string{"bootstrap-mainnet-1.example.com"}, cfg.Seeder.BootstrapPeers)
	assert.Equal(t, []string{"static-peer-1.example.com"}, cfg.Seeder.StaticPeers)
	assert.Equal(t, []config.Peer{{Host: "fn-peer-1.example.com", Port: 1234}}, cfg.FullNode.FullNodePeers)
}
