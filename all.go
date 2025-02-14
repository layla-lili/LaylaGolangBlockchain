// This file is auto-generated. Do not edit directly.
// Last updated: 2025-02-15 00:00:25

package main

// import (
//     "bytes"
//     "context"
//     "crypto/ecdsa"
//     "crypto/elliptic"
//     "crypto/rand"
//     "crypto/sha256"
//     "encoding/hex"
//     "encoding/json"
//     "errors"
//     "fmt"
//     "net/http"
//     "os"
//     "strings"
//     "sync"
//     "testing"
//     "time"

//     "github.com/cbergoon/merkletree"
//     "github.com/gorilla/mux"
//     "github.com/gorilla/websocket"
//     libp2p "github.com/libp2p/go-libp2p"
//     "github.com/libp2p/go-libp2p/core/host"
//     "github.com/libp2p/go-libp2p/core/network"
//     "github.com/libp2p/go-libp2p/core/peer"
//     mdns "github.com/libp2p/go-libp2p/p2p/discovery/mdns"
//     "github.com/multiformats/go-multiaddr"
// )

// /*
// # LaylaGolangBlockchain
// */

// // ======================
// // main.go
// // ======================

// package main

// import (
// 	"fmt"
// 	"os"

// 	// "sync"

// 	// "github.com/libp2p/go-libp2p/core/host"
// )

// // Global variable for wallet

// var Blockchain []Block
// // var pendingTransactions []Transaction
// // var p2pHost host.Host
// // var blockchainState *BlockchainState

// // var blockchainMutex sync.Mutex
// // var transactionMutex sync.Mutex

// func main() {
// 	// Initialize blockchain state
// 	state := NewBlockchainState()

// 	// Create and add genesis block
// 	genesisBlock := CreateGenesisBlock()
// 	if err := state.AddBlock(genesisBlock); err != nil {
// 		fmt.Printf("❌ Failed to add genesis block: %v\n", err)
// 		os.Exit(1)
// 	}

// 	// Initialize wallet
// 	wallet, err := NewWallet()
// 	if err != nil {
// 		fmt.Printf("❌ Failed to create wallet: %v\n", err)
// 		os.Exit(1)
// 	}
// 	state.SetWallet(wallet)

// 	// Initialize P2P host
// 	p2pHost, err := CreateLibp2pHost()
// 	if err != nil {
// 		fmt.Printf("❌ Failed to create libp2p host: %v\n", err)
// 		os.Exit(1)
// 	}
// 	state.SetP2PHost(p2pHost)

// 	// Setup P2P discovery and stream handler
// 	if err := SetupDiscovery(p2pHost); err != nil {
// 		fmt.Printf("❌ Failed to setup discovery: %v\n", err)
// 		os.Exit(1)
// 	}
// 	SetupStreamHandler(p2pHost)

// 	// Create and start server
// 	server := NewServer(state)
// 	apiPort := "8080"
// 	if len(os.Args) > 1 {
// 		apiPort = os.Args[1]
// 	}

// 	if err := server.Start(apiPort); err != nil {
// 		fmt.Printf("❌ Server failed: %v\n", err)
// 		os.Exit(1)
// 	}
// }

// // ======================
// // block.go
// // ======================

// package main

// import (
// 	"crypto/sha256"
// 	"encoding/hex"
// 	"fmt"
// 	"strings"
// 	"time"
// )

// // Update the Block struct
// type Block struct {
// 	Index        int           `json:"index"`
// 	Timestamp    string        `json:"timestamp"`
// 	Transactions []Transaction `json:"transactions"`
// 	PrevHash     string        `json:"prevHash"`
// 	Hash         string        `json:"hash"`
// 	Nonce        int           `json:"nonce"`
// 	MerkleRoot   []byte        `json:"merkleRoot"`
// 	Difficulty   int           `json:"difficulty"`
// }

// // Generate hash for a block
// func CalculateBlockHash(block Block) string {
// 	record := fmt.Sprintf("%d%s%s%d%s%d",
// 		block.Index,
// 		block.Timestamp,
// 		block.PrevHash,
// 		block.Nonce,
// 		hex.EncodeToString(block.MerkleRoot),
// 		block.Difficulty,
// 	)
// 	hash := sha256.Sum256([]byte(record))
// 	return hex.EncodeToString(hash[:])
// }

// func GenerateBlock(prevBlock Block, transactions []Transaction) Block {
// 	newBlock := Block{
// 		Index:        prevBlock.Index + 1,
// 		Timestamp:    time.Now().String(),
// 		Transactions: transactions,
// 		PrevHash:     prevBlock.Hash,
// 		Nonce:        0,
// 		Difficulty:   4, // Adjust difficulty as needed
// 	}

// 	// Calculate Merkle root using GetMerkleRoot instead of CalculateMerkleRoot
// 	merkleRoot, err := GetMerkleRoot(transactions)
// 	if err != nil {
// 		fmt.Printf("❌ Failed to create Merkle root: %v\n", err)
// 	} else {
// 		newBlock.MerkleRoot = merkleRoot
// 	}

// 	// Mine the block
// 	target := strings.Repeat("0", newBlock.Difficulty)
// 	for !strings.HasPrefix(newBlock.Hash, target) {
// 		newBlock.Nonce++
// 		newBlock.Hash = CalculateBlockHash(newBlock)
// 	}

// 	return newBlock
// }

// // Genesis Block (first block)
// func CreateGenesisBlock() Block {
// 	timestamp := time.Now().String()
// 	nonce := 0
// 	record := fmt.Sprintf("Genesis-%s-%d", timestamp, nonce)
// 	hash := sha256.Sum256([]byte(record))

// 	genesisBlock := Block{
// 		Index:     0,
// 		Timestamp: timestamp,
// 		PrevHash:  "",
// 		Hash:      hex.EncodeToString(hash[:]),
// 		Nonce:     nonce,
// 	}

// 	fmt.Printf("🔨 Created dynamic Genesis Block: %v\n", genesisBlock)
// 	return genesisBlock
// }

// func ValidateBlockchain(chain []Block) bool {
// 	for i := 1; i < len(chain); i++ {
// 		currentBlock := chain[i]
// 		previousBlock := chain[i-1]

// 		if currentBlock.PrevHash != previousBlock.Hash {
// 			return false
// 		}

// 		if currentBlock.Hash != CalculateBlockHash(currentBlock) {
// 			return false
// 		}
// 	}
// 	return true
// }

// func ValidateBlock(block Block, previousBlock Block) error {
// 	if block.Index != previousBlock.Index+1 {
// 		return fmt.Errorf("invalid block index")
// 	}
// 	if block.PrevHash != previousBlock.Hash {
// 		return fmt.Errorf("invalid previous hash")
// 	}
// 	if block.Hash != CalculateBlockHash(block) {
// 		return fmt.Errorf("invalid block hash")
// 	}

// 	return nil
// }

// func NewBlock(transactions []Transaction, prevBlock Block) (*Block, error) {
// 	block := &Block{
// 		Index:        prevBlock.Index + 1,
// 		Timestamp:    time.Now().String(),
// 		Transactions: transactions,
// 		PrevHash:     prevBlock.Hash,
// 		Difficulty:   4, // Default difficulty
// 		Nonce:        0,
// 	}

// 	// Calculate Merkle root
// 	merkleRoot, err := GetMerkleRoot(transactions)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to calculate merkle root: %w", err)
// 	}
// 	block.MerkleRoot = merkleRoot

// 	// Mine the block
// 	block.Hash = MineBlock(block)

// 	return block, nil
// }

// func MineBlock(block *Block) string {
// 	target := strings.Repeat("0", block.Difficulty)
// 	for {
// 		hash := CalculateBlockHash(*block)
// 		if strings.HasPrefix(hash, target) {
// 			return hash
// 		}
// 		block.Nonce++
// 	}
// }

// // ======================
// // block_test.go
// // ======================

// package main

// import (
// 	"bytes"
// 	"strings"
// 	"testing"
// )

// func TestMerkleTreeVerification(t *testing.T) {
// 	transactions := []Transaction{
// 		{Sender: "Alice", Receiver: "Bob", Amount: 10},
// 		{Sender: "Bob", Receiver: "Charlie", Amount: 5},
// 	}

// 	// Calculate TxIDs
// 	for i := range transactions {
// 		transactions[i].TxID = CalculateTxID(transactions[i])
// 	}

// 	merkleTree, err := NewMerkleTree(transactions)
// 	if err != nil {
// 		t.Fatalf("Failed to create Merkle tree: %v", err)
// 	}

// 	block := GenerateBlock(Block{}, transactions)

// 	if !bytes.Equal(block.MerkleRoot, merkleTree.MerkleRoot()) {
// 		t.Error("Merkle root mismatch")
// 	}
// }

// func TestBlockGeneration(t *testing.T) {
// 	genesis := CreateGenesisBlock()
// 	tx := Transaction{Sender: "Alice", Receiver: "Bob", Amount: 5.0}
// 	tx.TxID = CalculateTxID(tx)
// 	block := GenerateBlock(genesis, []Transaction{tx})

// 	if block.Index != genesis.Index+1 {
// 		t.Errorf("Expected block index %d, got %d", genesis.Index+1, block.Index)
// 	}
// 	if block.PrevHash != genesis.Hash {
// 		t.Errorf("Previous hash mismatch")
// 	}
// 	if !ValidateBlockchain([]Block{genesis, block}) {
// 		t.Error("Blockchain validation failed")
// 	}
// }

// func TestBlockMining(t *testing.T) {
// 	tests := []struct {
// 		name       string
// 		difficulty int
// 		wantPrefix string
// 	}{
// 		{"Difficulty 1", 1, "0"},
// 		{"Difficulty 2", 2, "00"},
// 		{"Difficulty 3", 3, "000"},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			block := &Block{
// 				Index:      1,
// 				Timestamp:  "test",
// 				Difficulty: tt.difficulty,
// 			}

// 			hash := MineBlock(block)
// 			if !strings.HasPrefix(hash, tt.wantPrefix) {
// 				t.Errorf("MineBlock() = %v, want prefix %v", hash, tt.wantPrefix)
// 			}
// 		})
// 	}
// }

// func TestBlockValidation(t *testing.T) {
// 	// Create test transactions
// 	tx := Transaction{
// 		Sender:   "Alice",
// 		Receiver: "Bob",
// 		Amount:   10.0,
// 	}
// 	transactions := []Transaction{tx}

// 	// Create and mine a block
// 	prevBlock := Block{Index: 0, Hash: "genesis"}
// 	block, err := NewBlock(transactions, prevBlock)
// 	if err != nil {
// 		t.Fatalf("NewBlock() error = %v", err)
// 	}

// 	// Verify block properties
// 	if block.Index != prevBlock.Index+1 {
// 		t.Errorf("NewBlock() index = %v, want %v", block.Index, prevBlock.Index+1)
// 	}
// 	if block.PrevHash != prevBlock.Hash {
// 		t.Errorf("NewBlock() prevHash = %v, want %v", block.PrevHash, prevBlock.Hash)
// 	}
// 	if block.MerkleRoot == nil {
// 		t.Error("NewBlock() merkleRoot is nil")
// 	}
// }

// // ======================
// // transaction.go
// // ======================

// package main

// import (
// 	"crypto/ecdsa"
// 	"crypto/elliptic"
// 	"crypto/sha256"
// 	"encoding/hex"
// 	"fmt"
// 	"time"
// )

// // Transaction represents a simple transaction
// type Transaction struct {
// 	Sender    string  // Public Key (Address)
// 	Receiver  string  // Public Key (Address)
// 	Amount    float64 // Amount of currency
// 	TxID      string  // Transaction ID (Hash)
// 	Signature []byte
// 	Timestamp time.Time
// 	Fee       float64
// }

// // type Transaction struct {
// // 	Sender    string `json:"sender"`
// // 	Receiver string `json:"receiver"`
// // 	Amount    float64 `json:"amount"`
// // 	Signature string `json:"signature"`
// // }

// // UTXO represents an unspent transaction output
// type UTXO struct {
// 	TxID   string
// 	Index  int
// 	Amount float64
// }

// // Calculate transaction hash (TxID)
// func CalculateTxID(tx Transaction) string {
// 	record := tx.Sender + tx.Receiver + fmt.Sprintf("%f", tx.Amount)
// 	hash := sha256.Sum256([]byte(record))
// 	return hex.EncodeToString(hash[:])
// }

// func CalculateBlockReward(block Block) float64 {
// 	baseReward := 50.0 // Base mining reward
// 	totalFees := 0.0

// 	for _, tx := range block.Transactions {
// 		totalFees += tx.Fee
// 	}
// 	return baseReward + totalFees
// }


// func ValidateTransaction(tx Transaction, publicKey []byte) bool {
// 	if tx.Amount <= 0 {
// 		return false
// 	}

// 	txHash := sha256.Sum256([]byte(tx.TxID))
// 	x, y := elliptic.Unmarshal(elliptic.P256(), publicKey)
// 	if x == nil {
// 		return false
// 	}

// 	publicKeyECDSA := &ecdsa.PublicKey{
// 		Curve: elliptic.P256(),
// 		X:     x,
// 		Y:     y,
// 	}

// 	return ecdsa.VerifyASN1(publicKeyECDSA, txHash[:], tx.Signature)
// }

// // ======================
// // chain.go
// // ======================

// package main



// const (
// 	BlockProtocol = "/block/1.0.0"
// 	TxProtocol    = "/tx/1.0.0"
// 	ChainProtocol = "/chain/1.0.0"
// )

// type Message struct {
// 	Type    string      `json:"type"`
// 	Payload interface{} `json:"payload"`
// }

// type BlockchainDB interface {
// 	SaveBlock(block Block) error
// 	GetBlock(hash string) (Block, error)
// 	SaveChain(chain []Block) error
// 	LoadChain() ([]Block, error)
// }

// func SelectBestChain(chains [][]Block) []Block {
// 	var bestChain []Block
// 	maxLength := 0

// 	for _, chain := range chains {
// 		if len(chain) > maxLength && ValidateBlockchain(chain) {
// 			maxLength = len(chain)
// 			bestChain = chain
// 		}
// 	}
// 	return bestChain
// }

// // ======================
// // mempool.go
// // ======================

// package main

// import (
// 	"fmt"
// 	"sync"
// 	"time"
// )

// type Mempool struct {
// 	transactions map[string]Transaction
// 	mutex        sync.RWMutex
// }

// func NewMempool() *Mempool {
// 	return &Mempool{
// 		transactions: make(map[string]Transaction),
// 	}
// }

// func (m *Mempool) AddTransaction(tx Transaction) error {
// 	m.mutex.Lock()
// 	defer m.mutex.Unlock()

// 	if !ValidateTransaction(tx, []byte(tx.Sender)) {
// 		return fmt.Errorf("invalid transaction")
// 	}

// 	m.transactions[tx.TxID] = tx
// 	return nil
// }

// func (m *Mempool) GetTransactions() []Transaction {
// 	m.mutex.RLock()
// 	defer m.mutex.RUnlock()

// 	txs := make([]Transaction, 0, len(m.transactions))
// 	for _, tx := range m.transactions {
// 		txs = append(txs, tx)
// 	}
// 	return txs
// }

// func (m *Mempool) RemoveTransactions(txs []Transaction) {
// 	m.mutex.Lock()
// 	defer m.mutex.Unlock()

// 	for _, tx := range txs {
// 		delete(m.transactions, tx.TxID)
// 	}
// }

// // Cleanup old transactions
// func (m *Mempool) CleanupOldTransactions(maxAge time.Duration) {
// 	m.mutex.Lock()
// 	defer m.mutex.Unlock()

// 	now := time.Now()
// 	for txID, tx := range m.transactions {
// 		if now.Sub(tx.Timestamp) > maxAge {
// 			delete(m.transactions, txID)
// 		}
// 	}
// }

// // ======================
// // merkle.go
// // ======================

// package main

// import (
// 	"crypto/sha256"
// 	"errors"
// 	"fmt"

// 	"github.com/cbergoon/merkletree"
// )

// // Make Transaction implement merkletree.Content interface
// func (tx Transaction) CalculateHash() ([]byte, error) {
// 	// Use existing TxID calculation
// 	hash := sha256.Sum256([]byte(tx.TxID))
// 	return hash[:], nil
// }

// func (tx Transaction) Equals(other merkletree.Content) (bool, error) {
// 	otherTx, ok := other.(Transaction)
// 	if !ok {
// 		return false, errors.New("invalid content type for comparison")
// 	}
// 	return tx.TxID == otherTx.TxID, nil
// }

// // NewMerkleTree creates a new Merkle tree from transactions
// func NewMerkleTree(transactions []Transaction) (*merkletree.MerkleTree, error) {
// 	if len(transactions) == 0 {
// 		return nil, errors.New("cannot create Merkle tree with no transactions")
// 	}

// 	var list []merkletree.Content
// 	for _, tx := range transactions {
// 		list = append(list, tx)
// 	}

// 	tree, err := merkletree.NewTree(list)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return tree, nil
// }

// // VerifyTransaction verifies if a transaction is part of the block
// func VerifyTransaction(block Block, tx Transaction) (bool, error) {
// 	tree, err := NewMerkleTree(block.Transactions)
// 	if err != nil {
// 		return false, err
// 	}

// 	return tree.VerifyContent(tx)
// }

// // GetMerkleRoot returns the Merkle root of transactions
// func GetMerkleRoot(transactions []Transaction) ([]byte, error) {
// 	tree, err := NewMerkleTree(transactions)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return tree.MerkleRoot(), nil
// }

// // VerifyTransactionInBlock verifies if a transaction is included in the block's Merkle tree
// func VerifyTransactionInBlock(block Block, tx Transaction) (bool, error) {
// 	// Create a new Merkle tree from block transactions
// 	tree, err := NewMerkleTree(block.Transactions)
// 	if err != nil {
// 		return false, fmt.Errorf("failed to create Merkle tree: %w", err)
// 	}

// 	// Verify the transaction is in the tree
// 	return tree.VerifyContent(tx)
// }

// // ======================
// // wallet.go
// // ======================

// package main

// import (
// 	"crypto/ecdsa"
// 	"crypto/elliptic"
// 	"crypto/rand"
// 	"crypto/sha256"
// 	"encoding/hex"
// 	"fmt"
// 	"math/big"
// )

// // Wallet structure
// type Wallet struct {
// 	PrivateKey *ecdsa.PrivateKey
// 	PublicKey  []byte
// 	Address    string
// 	UTXOs      []UTXO // <- Adding this field to keep track of UTXOs
// }

// // Create a new wallet
// func NewWallet() (*Wallet, error) {
// 	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to generate wallet: %v", err)
// 	}

// 	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
// 	address := generateAddress(publicKey)

// 	return &Wallet{
// 		PrivateKey: privateKey,
// 		PublicKey:  publicKey,
// 		Address:    address,
// 		UTXOs:      []UTXO{}, // Initialize the UTXOs slice
// 	}, nil
// }

// func generateAddress(publicKey []byte) string {
// 	hash := sha256.Sum256(publicKey)
// 	return hex.EncodeToString(hash[:])
// }

// // Sign a transaction
// func (w *Wallet) SignTransaction(tx *Transaction) error {
// 	txHash := sha256.Sum256([]byte(tx.TxID))
// 	signature, err := ecdsa.SignASN1(rand.Reader, w.PrivateKey, txHash[:])
// 	if err != nil {
// 		return fmt.Errorf("failed to sign transaction: %v", err)
// 	}
// 	tx.Signature = signature
// 	return nil
// }

// // Verify transaction signature
// func VerifySignature(tx *Transaction, signature string, pubKey ecdsa.PublicKey) bool {
// 	txHash := sha256.Sum256([]byte(tx.TxID))
// 	r := new(big.Int)
// 	s := new(big.Int)
// 	sigLen := len(signature) / 2
// 	r.SetString(signature[:sigLen], 16)
// 	s.SetString(signature[sigLen:], 16)

// 	return ecdsa.Verify(&pubKey, txHash[:], r, s)
// }

// // func validateTransaction(tx Transaction, wallet *Wallet) bool {
// // 	// Check if the sender has enough funds by validating against the UTXO set
// // 	for _, utxo := range wallet.UTXOs {
// // 		if utxo.TxID == tx.TxID && utxo.Amount >= tx.Amount {
// // 			return true
// // 		}
// // 	}
// // 	return false
// // }

// // func removeSpentTransactions(tx Transaction) {
// // 	// Loop over pending transactions and remove spent ones
// // 	for i, t := range pendingTransactions {
// // 		if t.TxID == tx.TxID {
// // 			// Remove from pending transactions
// // 			pendingTransactions = append(pendingTransactions[:i], pendingTransactions[i+1:]...)
// // 			break
// // 		}
// // 	}
// // }

// // ======================
// // p2p.go
// // ======================

// package main

// import (
// 	"fmt"
// 	"net/http"
// )


// // Start WebSocket server
// func StartP2PServer(port string) {
// 	go http.ListenAndServe(":"+port, nil)
// 	fmt.Println("P2P Server running on port", port)
// }

// // ======================
// // p2plibp2p.go
// // ======================

// // p2plibp2p.go
// package main

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"fmt"

// 	// "time"

// 	libp2p "github.com/libp2p/go-libp2p"

// 	"github.com/libp2p/go-libp2p/core/host"
// 	"github.com/libp2p/go-libp2p/core/network"
// 	"github.com/libp2p/go-libp2p/core/peer"
// 	mdns "github.com/libp2p/go-libp2p/p2p/discovery/mdns"
// )

// const protocolID = "/blockchain/1.0.0"
// const DiscoveryServiceTag = "blockchain-discovery"

// // Notifee implements the mdns.Notifee interface for peer discovery.
// type Notifee struct {
// 	h host.Host
// }

// func (n *Notifee) HandlePeerFound(pi peer.AddrInfo) {
// 	fmt.Printf("🎯 Peer Discovered: %s - Attempting connection...\n", pi.ID)

// 	err := n.h.Connect(context.Background(), pi)
// 	if err != nil {
// 		fmt.Printf("❌ Failed to connect to peer %s: %v\n", pi.ID, err)
// 	} else {
// 		fmt.Printf("🔗 Successfully connected to peer: %s\n", pi.ID)
// 	}
// }

// // SetupDiscovery starts mDNS discovery service.
// // func SetupDiscovery(h host.Host) error {
// // 	// Create a new Notifee
// // 	n := &Notifee{h: h}

// // 	// Initialize MDNS service with the notifee
// // 	// Note: NewMdnsService only returns the service, not an error
// // 	_ = mdns.NewMdnsService(h, DiscoveryServiceTag, n)

// //		return nil
// //	}
// func SetupDiscovery(h host.Host) error {
// 	n := &Notifee{h: h}

// 	service := mdns.NewMdnsService(h, DiscoveryServiceTag, n)
// 	if service == nil {
// 		return fmt.Errorf("failed to start mDNS service")
// 	}

// 	fmt.Println("✅ mDNS discovery service started. Waiting for peers...")

// 	return nil
// }

// // SetupStreamHandler registers a handler for incoming streams on our protocol.
// func SetupStreamHandler(h host.Host) {
// 	h.SetStreamHandler(protocolID, func(s network.Stream) {
// 		fmt.Println("Received new stream from:", s.Conn().RemotePeer().String())
// 		var receivedChain []Block
// 		decoder := json.NewDecoder(s)
// 		if err := decoder.Decode(&receivedChain); err != nil {
// 			fmt.Println("Error decoding blockchain from stream:", err)
// 		} else {
// 			fmt.Println("Received blockchain from peer:")
// 			for _, b := range receivedChain {
// 				fmt.Printf("Block %d: %s\n", b.Index, b.Hash)
// 			}
// 			// Here, you could add logic to compare and merge chains.
// 		}

// 		// Verify received block's Merkle root
// 		for _, block := range receivedChain {
// 			merkleRoot, err := GetMerkleRoot(block.Transactions)
// 			if err != nil {
// 				fmt.Printf("❌ Failed to create Merkle tree for block %d: %v\n", block.Index, err)
// 				continue
// 			}

// 			if !bytes.Equal(merkleRoot, block.MerkleRoot) {
// 				fmt.Printf("❌ Invalid Merkle root in block %d\n", block.Index)
// 				continue
// 			}

// 			// Verify all transactions in the block
// 			for _, tx := range block.Transactions {
// 				isValid, err := VerifyTransactionInBlock(block, tx)
// 				if err != nil {
// 					fmt.Printf("⚠️ Warning: Failed to verify transaction %s: %v\n", tx.TxID, err)
// 					continue
// 				}

// 				if !isValid {
// 					fmt.Printf("❌ Invalid transaction detected in block %d: %s\n", block.Index, tx.TxID)
// 					// Handle invalid transaction (maybe reject the block)
// 				}
// 			}
// 		}

// 		s.Close()
// 	})
// }

// // BroadcastBlockchain sends the current blockchain to all connected peers.
// //
// //	func BroadcastBlockchain(h host.Host, blockchain []Block) {
// //		peers := h.Network().Peers()
// //		fmt.Println("Broadcasting blockchain to", len(peers), "peers.")
// //		for _, p := range peers {
// //			s, err := h.NewStream(context.Background(), p, protocolID)
// //			if err != nil {
// //				fmt.Println("Error opening stream to peer", p.String(), ":", err)
// //				continue
// //			}
// //			encoder := json.NewEncoder(s)
// //			if err := encoder.Encode(blockchain); err != nil {
// //				fmt.Println("Error sending blockchain to peer", p.String(), ":", err)
// //			}
// //			s.Close()
// //		}
// //	}
// func BroadcastBlockchain(h host.Host, blockchain []Block) {
// 	peers := h.Network().Peers()
// 	fmt.Println("Broadcasting blockchain to", len(peers), "peers.")
// 	for _, p := range peers {
// 		s, err := h.NewStream(context.Background(), p, protocolID)
// 		if err != nil {
// 			fmt.Println("Error opening stream to peer", p.String(), ":", err)
// 			continue
// 		}
// 		encoder := json.NewEncoder(s)
// 		if err := encoder.Encode(blockchain); err != nil {
// 			fmt.Println("Error sending blockchain to peer", p.String(), ":", err)
// 		} else {
// 			fmt.Printf("Blockchain sent to peer %s\n", p.String())
// 		}
// 		s.Close()
// 	}
// }

// // CreateLibp2pHost creates a new libp2p host.
// func CreateLibp2pHost() (host.Host, error) {
// 	// Create a new libp2p host with default options
// 	fmt.Println("🚀 Creating libp2p host...")
// 	h, err := libp2p.New()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create libp2p host: %w", err)
// 	}
// 	return h, nil
// }

// // ======================
// // server.go
// // ======================

// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"

// 	"github.com/gorilla/mux"
// )

// type Server struct {
// 	state *BlockchainState
// }

// func NewServer(state *BlockchainState) *Server {
// 	return &Server{state: state}
// }

// // GET /chain - Get full blockchain
// func (s *Server) getBlockchain(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	if err := json.NewEncoder(w).Encode(s.state.GetChain()); err != nil {
// 		http.Error(w, "Failed to encode blockchain", http.StatusInternalServerError)
// 		return
// 	}
// }

// // POST /transaction - Create a new transaction
// func (s *Server) createTransaction(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	var tx Transaction
// 	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
// 		http.Error(w, "Invalid transaction data", http.StatusBadRequest)
// 		return
// 	}

// 	// Generate TxID if not present
// 	if tx.TxID == "" {
// 		tx.TxID = CalculateTxID(tx)
// 	}

// 	// Validate and sign transaction using state wallet
// 	if err := s.state.GetWallet().SignTransaction(&tx); err != nil {
// 		http.Error(w, fmt.Sprintf("Failed to sign transaction: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// 	if err := s.state.AddTransaction(tx); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(tx)
// }

// // GET /mine - Mine a new block
// func (s *Server) mineBlock(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	transactions := s.state.GetPendingTransactions()
// 	if len(transactions) == 0 {
// 		http.Error(w, "No transactions to mine", http.StatusBadRequest)
// 		return
// 	}

// 	lastBlock := s.state.GetLastBlock()
// 	newBlock := GenerateBlock(lastBlock, transactions)

// 	if err := s.state.AddBlock(newBlock); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Broadcast new block to peers
// 	BroadcastBlockchain(s.state.GetP2PHost(), s.state.GetChain())

// 	json.NewEncoder(w).Encode(newBlock)
// }

// // GET /peers - Get connected peers
// func (s *Server) getPeers(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	peers := s.state.GetP2PHost().Network().Peers()
// 	peerList := make([]string, 0, len(peers))
// 	for _, peer := range peers {
// 		peerList = append(peerList, peer.String())
// 	}

// 	if err := json.NewEncoder(w).Encode(peerList); err != nil {
// 		http.Error(w, "Failed to encode peer list", http.StatusInternalServerError)
// 		return
// 	}
// }

// // Start starts the HTTP API server
// func (s *Server) Start(port string) error {
// 	router := mux.NewRouter()

// 	// Register routes
// 	router.HandleFunc("/chain", s.getBlockchain).Methods("GET")
// 	router.HandleFunc("/transaction", s.createTransaction).Methods("POST")
// 	router.HandleFunc("/mine", s.mineBlock).Methods("GET")
// 	router.HandleFunc("/peers", s.getPeers).Methods("GET")

// 	// Add middleware
// 	router.Use(loggingMiddleware)

// 	fmt.Printf("🚀 API Server running on port %s\n", port)
// 	return http.ListenAndServe(":"+port, router)
// }

// // Middleware for logging requests
// func loggingMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Printf("📨 %s %s\n", r.Method, r.RequestURI)
// 		next.ServeHTTP(w, r)
// 	})
// }

// // ======================
// // consensus.go
// // ======================

// package main

// import (
// 	"bytes"
// 	"sync"
// )

// type Consensus struct {
// 	chainLock sync.RWMutex
// }

// func (c *Consensus) HandleChainSync(receivedChain []Block) bool {
// 	c.chainLock.Lock()
// 	defer c.chainLock.Unlock()

// 	if !ValidateChain(receivedChain) {
// 		return false
// 	}

// 	if len(receivedChain) <= len(Blockchain) {
// 		return false
// 	}

// 	// Verify work on all blocks
// 	for _, block := range receivedChain {
// 		if !ValidateProofOfWork(block) {
// 			return false
// 		}
// 	}

// 	// Replace our chain
// 	Blockchain = receivedChain
// 	return true
// }

// func ValidateChain(chain []Block) bool {
// 	for i := 1; i < len(chain); i++ {
// 		currentBlock := chain[i]
// 		previousBlock := chain[i-1]

// 		if !bytes.Equal([]byte(currentBlock.PrevHash), []byte(previousBlock.Hash)) {
// 			return false
// 		}

// 		if !ValidateBlockTransactions(currentBlock) {
// 			return false
// 		}
// 	}
// 	return true
// }

// func ValidateBlockTransactions(block Block) bool {
// 	for _, tx := range block.Transactions {
// 		if !ValidateTransaction(tx, []byte(tx.Sender)) {
// 			return false
// 		}
// 	}
// 	return true
// }

// func ValidateProofOfWork(block Block) bool {
// 	prefix := make([]byte, block.Difficulty)
// 	for i := 0; i < block.Difficulty; i++ {
// 		prefix[i] = '0'
// 	}

// 	hash := CalculateBlockHash(block)
// 	return bytes.HasPrefix([]byte(hash), prefix)
// }

