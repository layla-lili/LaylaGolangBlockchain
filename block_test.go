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
		defer func() { done <- true }() // Ensure done is always signaled

		// Setup test state
		state := NewBlockchainState()
		consensus := NewConsensus(state)

		// Create test chain with minimal difficulty
		genesis := CreateGenesisBlock()
		genesis.Difficulty = 1 // Minimal difficulty
		if err := state.AddBlock(genesis); err != nil {
			t.Errorf("Failed to add genesis block: %v", err)
			return
		}

		// Create and sign test transaction
		wallet, err := NewWallet()
		if err != nil {
			t.Errorf("Failed to create wallet: %v", err)
			return
		}

		tx := Transaction{
			Receiver:  "Bob",
			Amount:    10.0,
			Timestamp: time.Now(),
		}
		if err := wallet.SignTransaction(&tx); err != nil {
			t.Errorf("Failed to sign transaction: %v", err)
			return
		}

		// Create and mine block with minimal difficulty
		block := Block{
			Index:        genesis.Index + 1,
			Timestamp:    time.Now().String(),
			Transactions: []Transaction{tx},
			PrevHash:     genesis.Hash,
			Difficulty:   1, // Minimal difficulty
		}

		// Quick mining with timeout
		miningDone := make(chan bool)
		go func() {
			block.Hash = MineBlock(&block)
			miningDone <- true
		}()

		select {
		case <-time.After(1 * time.Second):
			t.Log("Mining timed out, using simple hash")
			block.Hash = CalculateBlockHash(block)
		case <-miningDone:
			t.Log("Mining completed successfully")
		}

		// Test chain validation
		testChain := []Block{genesis, block}
		if !consensus.ValidateChain(testChain) {
			t.Error("Chain validation failed")
			return
		}

		// Test proof of work
		if !consensus.ValidateProofOfWork(block) {
			t.Error("Proof of work validation failed")
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

func TestWalletAndTransactionConsistency(t *testing.T) {
	wallet, err := NewWallet()
	if err != nil {
		t.Fatalf("Failed to create wallet: %v", err)
	}

	// Address should be non-empty and in correct format
	if len(wallet.Address) == 0 {
		t.Error("Generated address is empty")
	}

	// Address should be Base58 encoded
	decoded, err := base58.Decode(wallet.Address)
	if err != nil {
		t.Fatalf("Failed to decode address: %v", err)
	}
	if len(decoded) != 25 { // 1 version + 20 hash + 4 checksum bytes
		t.Errorf("Invalid address length: got %d, want 25", len(decoded))
	}

	// Create wallet
	wallet, err = NewWallet()
	if err != nil {
		t.Fatalf("Failed to create wallet: %v", err)
	}

	// Create transaction
	tx := Transaction{
		Receiver:  "recipient123",
		Amount:    10.0,
		Timestamp: time.Now(),
	}

	// Sign transaction
	if err := wallet.SignTransaction(&tx); err != nil {
		t.Fatalf("Failed to sign transaction: %v", err)
	}

	// Verify all components are consistent
	tests := []struct {
		name string
		want bool
		got  bool
	}{
		{
			name: "Address matches public key",
			want: true,
			got:  GenerateAddress(tx.SenderPublicKey) == tx.SenderAddress,
		},
		{
			name: "Signature is valid",
			want: true,
			got:  ValidateTransaction(tx, tx.Signature),
		},
		{
			name: "Address matches wallet",
			want: true,
			got:  wallet.GetAddress() == tx.SenderAddress,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("%s: got %v, want %v", tt.name, tt.got, tt.want)
				if tt.name == "Signature is valid" {
					t.Logf("Transaction details:\n"+
						"SenderAddress: %s\n"+
						"PublicKey: %x\n"+
						"Signature: %x\n"+
						"TxID: %s",
						tx.SenderAddress,
						tx.SenderPublicKey,
						tx.Signature,
						tx.TxID)
				}
			}
		})
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
