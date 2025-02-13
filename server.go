package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Get full blockchain
func getBlockchain(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Blockchain)
}

// Create a new transaction
func createTransaction(w http.ResponseWriter, r *http.Request) {
	var tx Transaction
	_ = json.NewDecoder(r.Body).Decode(&tx)
	tx.TxID = CalculateTxID(tx)

	// Add transaction to pending transactions
	pendingTransactions = append(pendingTransactions, tx)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tx)
}

// Mine a new block
func mineBlock(w http.ResponseWriter, r *http.Request) {
	if len(pendingTransactions) == 0 {
		http.Error(w, "No transactions to mine", http.StatusBadRequest)
		return
	}

	newBlock := GenerateBlock(Blockchain[len(Blockchain)-1], pendingTransactions)
	Blockchain = append(Blockchain, newBlock)
	pendingTransactions = []Transaction{}

	// Broadcast updated blockchain using libp2p
	BroadcastBlockchain(p2pHost, Blockchain)

	json.NewEncoder(w).Encode(newBlock)
}

// Get connected peers should use libp2p now
func getPeers(w http.ResponseWriter, r *http.Request) {
	peers := p2pHost.Network().Peers()
	peerList := make([]string, 0, len(peers))
	for _, peer := range peers {
		peerList = append(peerList, peer.String())
	}
	json.NewEncoder(w).Encode(peerList)
}

// Start HTTP API Server
func startAPIServer(port string) {
	router := mux.NewRouter()
	router.HandleFunc("/chain", getBlockchain).Methods("GET")
	router.HandleFunc("/transactions", createTransaction).Methods("POST")
	router.HandleFunc("/mine", mineBlock).Methods("GET")
	router.HandleFunc("/peers", getPeers).Methods("GET")

	fmt.Println("API Server running on port", port)
	http.ListenAndServe(":"+port, router)
}
