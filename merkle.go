package main

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/cbergoon/merkletree"
)

// Make Transaction implement merkletree.Content interface
func (tx Transaction) CalculateHash() ([]byte, error) {
	// Use existing TxID calculation
	hash := sha256.Sum256([]byte(tx.TxID))
	return hash[:], nil
}

func (tx Transaction) Equals(other merkletree.Content) (bool, error) {
	otherTx, ok := other.(Transaction)
	if !ok {
		return false, errors.New("invalid content type for comparison")
	}
	return tx.TxID == otherTx.TxID, nil
}

// NewMerkleTree creates a new Merkle tree from transactions
func NewMerkleTree(transactions []Transaction) (*merkletree.MerkleTree, error) {
	if len(transactions) == 0 {
		return nil, errors.New("cannot create Merkle tree with no transactions")
	}

	var list []merkletree.Content
	for _, tx := range transactions {
		list = append(list, tx)
	}

	tree, err := merkletree.NewTree(list)
	if err != nil {
		return nil, err
	}

	return tree, nil
}

// VerifyTransaction verifies if a transaction is part of the block
func VerifyTransaction(block Block, tx Transaction) (bool, error) {
	tree, err := NewMerkleTree(block.Transactions)
	if err != nil {
		return false, err
	}

	return tree.VerifyContent(tx)
}

// GetMerkleRoot returns the Merkle root of transactions
func GetMerkleRoot(transactions []Transaction) ([]byte, error) {
	tree, err := NewMerkleTree(transactions)
	if err != nil {
		return nil, err
	}

	return tree.MerkleRoot(), nil
}

// VerifyTransactionInBlock verifies if a transaction is included in the block's Merkle tree
func VerifyTransactionInBlock(block Block, tx Transaction) (bool, error) {
	// Create a new Merkle tree from block transactions
	tree, err := NewMerkleTree(block.Transactions)
	if err != nil {
		return false, fmt.Errorf("failed to create Merkle tree: %w", err)
	}

	// Verify the transaction is in the tree
	return tree.VerifyContent(tx)
}
