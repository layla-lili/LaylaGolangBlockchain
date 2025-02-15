package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// Parse command line flags
	httpPort := flag.String("http", "8080", "HTTP server port")
	p2pPort := flag.String("p2p", "6001", "P2P network port")
	flag.Parse()

	// Override with positional args if provided
	args := flag.Args()
	if len(args) >= 1 {
		*httpPort = args[0]
	}
	if len(args) >= 2 {
		*httpPort = args[1]
	}

	fmt.Printf("🚀 Starting blockchain node...\n")
	fmt.Printf("HTTP Port: %s, P2P Port: %s\n", *httpPort, *p2pPort)

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

	// Initialize P2P host with specific port
	p2pHost, err := CreateLibp2pHost(*p2pPort)
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
	SetupStreamHandler(p2pHost, state)

	fmt.Println("🔍 DEBUG: Starting server initialization...")
	server := NewServer(state)

	// Start server with error handling
	if err := server.Start(*httpPort); err != nil {
		fmt.Printf("❌ Server failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✅ Node is ready! Access the following endpoints:\n")
	fmt.Printf("   REST API: http://localhost:%s\n", *httpPort)
	fmt.Printf("   P2P Network: /ip4/0.0.0.0/tcp/%s\n", *p2pPort)

	// Print available commands
	fmt.Println("\n📝 Available Commands:")
	fmt.Println("   GET  /chain       - View the blockchain")
	fmt.Println("   POST /transaction - Create a new transaction")
	fmt.Println("   GET  /mine        - Mine a new block")
	fmt.Println("   GET  /peers       - View connected peers")

	// Start CLI
	fmt.Println("\n💻 Starting CLI interface...")
	cli := NewCLI(*httpPort)
	cli.Start()
}
