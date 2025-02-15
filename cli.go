package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type CLI struct {
	baseURL string
}

func NewCLI(port string) *CLI {
	return &CLI{
		baseURL: fmt.Sprintf("http://localhost:%s", port),
	}
}

func (cli *CLI) Start() {
	for {
		fmt.Println("\n🚀 Blockchain CLI")
		fmt.Println("1. View blockchain")
		fmt.Println("2. Create transaction")
		fmt.Println("3. Mine block")
		fmt.Println("4. View peers")
		fmt.Println("5. Exit")

		var choice int
		fmt.Print("Enter choice (1-5): ")
		fmt.Scan(&choice)

		switch choice {
		case 1:
			cli.viewBlockchain()
		case 2:
			cli.createTransaction()
		case 3:
			cli.mineBlock()
		case 4:
			cli.viewPeers()
		case 5:
			return
		}
	}
}

func (cli *CLI) viewBlockchain() {
	resp, err := http.Get(cli.baseURL + "/chain")
	if err != nil {
		log.Printf("Error fetching blockchain: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var chain []Block
	if err := json.NewDecoder(resp.Body).Decode(&chain); err != nil {
		log.Printf("Error decoding chain: %v\n", err)
		return
	}

	fmt.Println("\n📦 Blockchain:")
	for _, block := range chain {
		fmt.Printf("\nBlock #%d\n", block.Index)
		fmt.Printf("Hash: %s\n", block.Hash)
		fmt.Printf("Transactions: %d\n", len(block.Transactions))
	}
}

func (cli *CLI) createTransaction() {
	tx := Transaction{
		Receiver:  "recipient123",
		Amount:    10.0,
		Timestamp: time.Now(),
	}

	jsonData, _ := json.Marshal(tx)
	resp, err := http.Post(cli.baseURL+"/transaction", "application/json",
		bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal("Error creating transaction:", err)
	}

	var result Transaction
	json.NewDecoder(resp.Body).Decode(&result)
	fmt.Printf("\n✅ Transaction created: %s\n", result.TxID)
}

func (cli *CLI) mineBlock() {
	resp, err := http.Get(cli.baseURL + "/mine")
	if err != nil {
		log.Fatal("Error mining block:", err)
	}

	var block Block
	json.NewDecoder(resp.Body).Decode(&block)
	fmt.Printf("\n⛏️  Mined block #%d\n", block.Index)
	fmt.Printf("Hash: %s\n", block.Hash)
}

func (cli *CLI) viewPeers() {
	resp, err := http.Get(cli.baseURL + "/peers")
	if err != nil {
		log.Fatal("Error fetching peers:", err)
	}

	var peers []string
	json.NewDecoder(resp.Body).Decode(&peers)
	fmt.Printf("\n👥 Connected peers: %d\n", len(peers))
	for _, peer := range peers {
		fmt.Println(peer)
	}
}
