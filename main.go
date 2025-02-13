package main

import (
	"fmt"
	"os"

	"github.com/libp2p/go-libp2p/core/host"
)

var Blockchain []Block
var pendingTransactions []Transaction
var p2pHost host.Host // Add this global variable

func main() {
	// Create Genesis Block
	genesisBlock := Block{Index: 0, Timestamp: "2024-02-13", PrevHash: "", Hash: "genesis"}
	Blockchain = append(Blockchain, genesisBlock)

	// Initialize libp2p host
	var err error
	p2pHost, err = CreateLibp2pHost()
	if err != nil {
		fmt.Printf("Failed to create libp2p host: %v\n", err)
		return
	}

	// Setup P2P discovery and stream handler
	if err := SetupDiscovery(p2pHost); err != nil {
		fmt.Printf("Failed to setup discovery: %v\n", err)
		return
	}
	SetupStreamHandler(p2pHost)

	// Start Servers
	apiPort := "8080"
	if len(os.Args) > 1 {
		apiPort = os.Args[1]
	}

	startAPIServer(apiPort)
}
