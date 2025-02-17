package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	// blockchain *Blockchain
	state *BlockchainState
}

func NewServer(state *BlockchainState) *Server {
	return &Server{state: state}
}

// GET /chain - Get full blockchain
func (s *Server) getBlockchain(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(s.state.GetChain()); err != nil {
		http.Error(w, "Failed to encode blockchain", http.StatusInternalServerError)
		return
	}
}

// POST /transaction - Create a new transaction
func (s *Server) createTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var tx Transaction
	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		http.Error(w, "Invalid transaction data", http.StatusBadRequest)
		return
	}

	// Generate TxID if not present
	if tx.TxID == "" {
		tx.TxID = CalculateTxID(tx)
	}

	// Validate and sign transaction using state wallet
	if err := s.state.GetWallet().SignTransaction(&tx); err != nil {
		http.Error(w, fmt.Sprintf("Failed to sign transaction: %v", err), http.StatusInternalServerError)
		return
	}

	if err := s.state.AddTransaction(tx); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tx)
}

// GET /mine - Mine a new block
func (s *Server) mineBlock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	transactions := s.state.GetPendingTransactions()
	if len(transactions) == 0 {
		http.Error(w, "No transactions to mine", http.StatusBadRequest)
		return
	}

	lastBlock := s.state.GetLastBlock()
	newBlock := GenerateBlock(lastBlock, transactions)

	if err := s.state.AddBlock(newBlock); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Broadcast new block to peers
	BroadcastBlockchain(s.state.GetP2PHost(), s.state.GetChain())

	json.NewEncoder(w).Encode(newBlock)
}

// GET /peers - Get connected peers
func (s *Server) getPeers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	peers := s.state.GetP2PHost().Network().Peers()
	peerList := make([]string, 0, len(peers))
	for _, peer := range peers {
		peerList = append(peerList, peer.String())
	}

	if err := json.NewEncoder(w).Encode(peerList); err != nil {
		http.Error(w, "Failed to encode peer list", http.StatusInternalServerError)
		return
	}
}

// Start starts the HTTP API server
func (s *Server) Start(port string) error {
	fmt.Println("🔍 DEBUG: Server.Start called")

	// Setup routes with logging middleware
	router := s.setupRoutesWithLogging()
	fmt.Println("🔍 DEBUG: Routes setup complete")

	addr := ":" + port
	fmt.Printf("🔍 DEBUG: About to listen on %s\n", addr)

	// Start server in a separate goroutine
	errChan := make(chan error)
	go func() {
		errChan <- http.ListenAndServe(addr, router)
	}()

	// Check for immediate errors
	select {
	case err := <-errChan:
		return fmt.Errorf("server failed: %w", err)
	case <-time.After(100 * time.Millisecond):
		fmt.Printf("✅ Server successfully started on port %s\n", port)
		return nil
	}
}

// Uncomment and update the logging middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		fmt.Printf("📨 Started %s %s\n", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
		fmt.Printf("✅ Completed %s %s in %v\n", r.Method, r.RequestURI, time.Since(startTime))
	})
}

// Add a new method to setup routes with logging
func (s *Server) setupRoutesWithLogging() http.Handler {
	router := s.setupRoutes()
	return loggingMiddleware(router)
}

// setupRoutes configures and returns the router with all HTTP endpoints
func (s *Server) setupRoutes() *http.ServeMux {
	router := http.NewServeMux()

	// Define routes
	router.HandleFunc("/chain", s.getBlockchain)
	router.HandleFunc("/transaction", s.createTransaction)
	router.HandleFunc("/mine", s.mineBlock)
	router.HandleFunc("/peers", s.getPeers)

	return router
}
