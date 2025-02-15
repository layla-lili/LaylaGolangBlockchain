package main

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/mr-tron/base58"
)

func TestMerkleTreeVerification(t *testing.T) {
	transactions := []Transaction{
		{SenderAddress: "Alice", Receiver: "Bob", Amount: 10},
		{SenderAddress: "Bob", Receiver: "Charlie", Amount: 5},
	}

	// Calculate TxIDs
	for i := range transactions {
		transactions[i].TxID = CalculateTxID(transactions[i])
	}

	merkleTree, err := NewMerkleTree(transactions)
	if err != nil {
		t.Fatalf("Failed to create Merkle tree: %v", err)
	}

	block := GenerateBlock(Block{}, transactions)

	if !bytes.Equal(block.MerkleRoot, merkleTree.MerkleRoot()) {
		t.Error("Merkle root mismatch")
	}
}

func TestBlockGeneration(t *testing.T) {
	genesis := CreateGenesisBlock()
	tx := Transaction{SenderAddress: "Alice", Receiver: "Bob", Amount: 5.0}
	tx.TxID = CalculateTxID(tx)
	block := GenerateBlock(genesis, []Transaction{tx})

	if block.Index != genesis.Index+1 {
		t.Errorf("Expected block index %d, got %d", genesis.Index+1, block.Index)
	}
	if block.PrevHash != genesis.Hash {
		t.Errorf("Previous hash mismatch")
	}
	if !ValidateBlockchain([]Block{genesis, block}) {
		t.Error("Blockchain validation failed")
	}
}

func TestBlockMining(t *testing.T) {
	tests := []struct {
		name       string
		difficulty int
		wantPrefix string
	}{
		{"Difficulty 1", 1, "0"},
		{"Difficulty 2", 2, "00"},
		{"Difficulty 3", 3, "000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block := &Block{
				Index:      1,
				Timestamp:  "test",
				Difficulty: tt.difficulty,
			}

			hash := MineBlock(block)
			if !strings.HasPrefix(hash, tt.wantPrefix) {
				t.Errorf("MineBlock() = %v, want prefix %v", hash, tt.wantPrefix)
			}
		})
	}
}

func TestBlockValidation(t *testing.T) {
	// Create test transactions
	tx := Transaction{
		SenderAddress: "Alice",
		Receiver:      "Bob",
		Amount:        10.0,
	}
	transactions := []Transaction{tx}

	// Create and mine a block
	prevBlock := Block{Index: 0, Hash: "genesis"}
	block, err := NewBlock(transactions, prevBlock)
	if err != nil {
		t.Fatalf("NewBlock() error = %v", err)
	}

	// Verify block properties
	if block.Index != prevBlock.Index+1 {
		t.Errorf("NewBlock() index = %v, want %v", block.Index, prevBlock.Index+1)
	}
	if block.PrevHash != prevBlock.Hash {
		t.Errorf("NewBlock() prevHash = %v, want %v", block.PrevHash, prevBlock.Hash)
	}
	if block.MerkleRoot == nil {
		t.Error("NewBlock() merkleRoot is nil")
	}
}

func TestConsensusValidation(t *testing.T) {
	// Increase timeout slightly but keep it reasonable
	timeout := time.After(3 * time.Second)
	done := make(chan bool)

	go func() {
		defer func() { done <- true }()

		state := NewBlockchainState()
		consensus := NewConsensus(state)

		// Create genesis block
		genesis := CreateGenesisBlock()
		genesis.Difficulty = 1
		if err := state.AddBlock(genesis); err != nil {
			t.Errorf("Failed to add genesis block: %v", err)
			return
		}

		// Create wallet and transaction
		wallet, err := NewWallet()
		if err != nil {
			t.Errorf("Failed to create wallet: %v", err)
			return
		}

		// Create transaction with fixed timestamp
		tx := Transaction{
			SenderAddress: wallet.GetAddress(),
			Receiver:      "Bob",
			Amount:        10.0,
			Timestamp:     time.Unix(1234567890, 0), // Use fixed timestamp
		}

		// Sign transaction
		if err := wallet.SignTransaction(&tx); err != nil {
			t.Errorf("Failed to sign transaction: %v", err)
			return
		}

		// Create and mine block
		block := Block{
			Index:        genesis.Index + 1,
			Timestamp:    time.Now().String(),
			Transactions: []Transaction{tx},
			PrevHash:     genesis.Hash,
			Difficulty:   1,
		}

		block.Hash = MineBlock(&block)
		t.Log("Block mined successfully")

		// Validate chain
		testChain := []Block{genesis, block}
		if !consensus.ValidateChain(testChain) {
			t.Error("Chain validation failed")
			// Add debug info
			t.Logf("Transaction details:\n"+
				"SenderAddress: %s\n"+
				"PublicKey: %x\n"+
				"Signature: %x\n"+
				"TxID: %s",
				tx.SenderAddress,
				tx.SenderPublicKey,
				tx.Signature,
				tx.TxID)
			return
		}
	}()

	// Wait for test completion or timeout
	select {
	case <-timeout:
		t.Fatal("Test timed out after 3 seconds")
	case <-done:
		t.Log("Test completed successfully")
	}
}

func TestGenerateAddress(t *testing.T) {
	wallet, err := NewWallet()
	if err != nil {
		t.Fatalf("Failed to create wallet: %v", err)
	}

	address := GenerateAddress(wallet.PublicKey)

	// Address should not be empty
	if address == "" {
		t.Error("Generated address is empty")
	}

	// Address should be decodable
	decoded, err := base58.Decode(address)
	if err != nil {
		t.Fatalf("Failed to decode address: %v", err)
	}

	// Check address length (1 version + 20 hash + 4 checksum bytes)
	if len(decoded) != 25 {
		t.Errorf("Invalid address length: got %d, want 25", len(decoded))
	}
}

func TestWalletAndTransactionConsistency(t *testing.T) {
	wallet, err := NewWallet()
	if err != nil {
		t.Fatalf("Failed to create wallet: %v", err)
	}

	// Create transaction with fixed timestamp for consistent hashing
	tx := Transaction{
		Receiver:  "recipient123",
		Amount:    10.0,
		Timestamp: time.Unix(1234567890, 0), // Use fixed timestamp
	}

	if err := wallet.SignTransaction(&tx); err != nil {
		t.Fatalf("Failed to sign transaction: %v", err)
	}

	// Debug logging
	t.Logf("Transaction details:\n"+
		"SenderAddress: %s\n"+
		"PublicKey: %x\n"+
		"Signature: %x\n"+
		"TxID: %s\n"+
		"Timestamp: %v",
		tx.SenderAddress,
		tx.SenderPublicKey,
		tx.Signature,
		tx.TxID,
		tx.Timestamp)

	if !ValidateTransaction(tx, tx.SenderPublicKey) {
		// Additional debug info if validation fails
		t.Error("Transaction validation failed")
		t.Logf("Derived address: %s", GenerateAddress(tx.SenderPublicKey))
	}
}
