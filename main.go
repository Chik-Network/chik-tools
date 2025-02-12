package main

import (
	"github.com/chik-network/chik-tools/cmd"
	_ "github.com/chik-network/chik-tools/cmd/certs"
	_ "github.com/chik-network/chik-tools/cmd/config"
	_ "github.com/chik-network/chik-tools/cmd/datalayer"
	_ "github.com/chik-network/chik-tools/cmd/debug"
	_ "github.com/chik-network/chik-tools/cmd/network"
	_ "github.com/chik-network/chik-tools/cmd/testnet"
)

func main() {
	cmd.Execute()
}
