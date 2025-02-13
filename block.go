package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// Block with transactions
type Block struct {
	Index        int
	Timestamp    string
	Transactions []Transaction
	Data      string
	PrevHash     string
	Hash         string
	Nonce        int
}

// Generate hash for a block
func CalculateBlockHash(block Block) string {
	record := fmt.Sprintf("%d%s%s%d", block.Index, block.Timestamp, block.PrevHash, block.Nonce)
	hash := sha256.Sum256([]byte(record))
	return hex.EncodeToString(hash[:])
}

// Create a new block with transactions
func GenerateBlock(prevBlock Block, transactions []Transaction) Block {
	newBlock := Block{
		Index:        prevBlock.Index + 1,
		Timestamp:    time.Now().String(),
		Transactions: transactions,
		PrevHash:     prevBlock.Hash,
		Nonce:        0, // PoW (Implemented later)
	}
	newBlock.Hash = CalculateBlockHash(newBlock)
	return newBlock
}

// Genesis Block (first block)
func CreateGenesisBlock() Block {
	return Block{
		Index:     0,
		Timestamp: time.Now().String(),
		Data:      "Genesis Block",
		PrevHash:  "",
		Hash:      "",
		Nonce:     0,
	}
}