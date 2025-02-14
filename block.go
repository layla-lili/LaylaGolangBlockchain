package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

// Update the Block struct
type Block struct {
	Index        int           `json:"index"`
	Timestamp    string        `json:"timestamp"`
	Transactions []Transaction `json:"transactions"`
	PrevHash     string        `json:"prevHash"`
	Hash         string        `json:"hash"`
	Nonce        int           `json:"nonce"`
	MerkleRoot   []byte        `json:"merkleRoot"`
	Difficulty   int           `json:"difficulty"`
}

// Generate hash for a block
func CalculateBlockHash(block Block) string {
	record := fmt.Sprintf("%d%s%s%d%s%d",
		block.Index,
		block.Timestamp,
		block.PrevHash,
		block.Nonce,
		hex.EncodeToString(block.MerkleRoot),
		block.Difficulty,
	)
	hash := sha256.Sum256([]byte(record))
	return hex.EncodeToString(hash[:])
}

func GenerateBlock(prevBlock Block, transactions []Transaction) Block {
	newBlock := Block{
		Index:        prevBlock.Index + 1,
		Timestamp:    time.Now().String(),
		Transactions: transactions,
		PrevHash:     prevBlock.Hash,
		Nonce:        0,
		Difficulty:   4, // Adjust difficulty as needed
	}

	// Calculate Merkle root using GetMerkleRoot instead of CalculateMerkleRoot
	merkleRoot, err := GetMerkleRoot(transactions)
	if err != nil {
		fmt.Printf("❌ Failed to create Merkle root: %v\n", err)
	} else {
		newBlock.MerkleRoot = merkleRoot
	}

	// Mine the block
	target := strings.Repeat("0", newBlock.Difficulty)
	for !strings.HasPrefix(newBlock.Hash, target) {
		newBlock.Nonce++
		newBlock.Hash = CalculateBlockHash(newBlock)
	}

	return newBlock
}

// Genesis Block (first block)
func CreateGenesisBlock() Block {
	timestamp := time.Now().String()
	nonce := 0
	record := fmt.Sprintf("Genesis-%s-%d", timestamp, nonce)
	hash := sha256.Sum256([]byte(record))

	genesisBlock := Block{
		Index:     0,
		Timestamp: timestamp,
		PrevHash:  "",
		Hash:      hex.EncodeToString(hash[:]),
		Nonce:     nonce,
	}

	fmt.Printf("🔨 Created dynamic Genesis Block: %v\n", genesisBlock)
	return genesisBlock
}

func ValidateBlockchain(chain []Block) bool {
	for i := 1; i < len(chain); i++ {
		currentBlock := chain[i]
		previousBlock := chain[i-1]

		if currentBlock.PrevHash != previousBlock.Hash {
			return false
		}

		if currentBlock.Hash != CalculateBlockHash(currentBlock) {
			return false
		}
	}
	return true
}

func ValidateBlock(block Block, previousBlock Block) error {
	if block.Index != previousBlock.Index+1 {
		return fmt.Errorf("invalid block index")
	}
	if block.PrevHash != previousBlock.Hash {
		return fmt.Errorf("invalid previous hash")
	}
	if block.Hash != CalculateBlockHash(block) {
		return fmt.Errorf("invalid block hash")
	}

	return nil
}

func NewBlock(transactions []Transaction, prevBlock Block) (*Block, error) {
	block := &Block{
		Index:        prevBlock.Index + 1,
		Timestamp:    time.Now().String(),
		Transactions: transactions,
		PrevHash:     prevBlock.Hash,
		Difficulty:   4, // Default difficulty
		Nonce:        0,
	}

	// Calculate Merkle root
	merkleRoot, err := GetMerkleRoot(transactions)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate merkle root: %w", err)
	}
	block.MerkleRoot = merkleRoot

	// Mine the block
	block.Hash = MineBlock(block)

	return block, nil
}

func MineBlock(block *Block) string {
	target := strings.Repeat("0", block.Difficulty)
	for {
		hash := CalculateBlockHash(*block)
		if strings.HasPrefix(hash, target) {
			return hash
		}
		block.Nonce++
	}
}
