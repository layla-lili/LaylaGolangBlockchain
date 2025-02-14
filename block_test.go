package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestMerkleTreeVerification(t *testing.T) {
	transactions := []Transaction{
		{Sender: "Alice", Receiver: "Bob", Amount: 10},
		{Sender: "Bob", Receiver: "Charlie", Amount: 5},
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
	tx := Transaction{Sender: "Alice", Receiver: "Bob", Amount: 5.0}
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
		Sender:   "Alice",
		Receiver: "Bob",
		Amount:   10.0,
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
