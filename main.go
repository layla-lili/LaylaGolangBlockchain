package main

import (
	"fmt"
	"os"

	// "sync"

	"github.com/libp2p/go-libp2p/core/host"
)

// Global variable for wallet

var Blockchain []Block
var pendingTransactions []Transaction
var p2pHost host.Host
var blockchainState *BlockchainState

// var blockchainMutex sync.Mutex
// var transactionMutex sync.Mutex

func main() {
	// Initialize blockchain state
	state := NewBlockchainState()

	// Create and add genesis block
	genesisBlock := CreateGenesisBlock()
	if err := state.AddBlock(genesisBlock); err != nil {
		fmt.Printf("❌ Failed to add genesis block: %v\n", err)
		os.Exit(1)
	}

	// Initialize wallet
	wallet, err := NewWallet()
	if err != nil {
		fmt.Printf("❌ Failed to create wallet: %v\n", err)
		os.Exit(1)
	}
	state.SetWallet(wallet)

	// Initialize P2P host
	p2pHost, err := CreateLibp2pHost()
	if err != nil {
		fmt.Printf("❌ Failed to create libp2p host: %v\n", err)
		os.Exit(1)
	}
	state.SetP2PHost(p2pHost)

	// Setup P2P discovery and stream handler
	if err := SetupDiscovery(p2pHost); err != nil {
		fmt.Printf("❌ Failed to setup discovery: %v\n", err)
		os.Exit(1)
	}
	SetupStreamHandler(p2pHost)

	// Create and start server
	server := NewServer(state)
	apiPort := "8080"
	if len(os.Args) > 1 {
		apiPort = os.Args[1]
	}

	if err := server.Start(apiPort); err != nil {
		fmt.Printf("❌ Server failed: %v\n", err)
		os.Exit(1)
	}
}
