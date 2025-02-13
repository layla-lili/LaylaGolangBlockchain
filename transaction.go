package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// Transaction represents a simple transaction
type Transaction struct {
	Sender    string  // Public Key (Address)
	Receiver  string  // Public Key (Address)
	Amount    float64 // Amount of currency
	TxID      string  // Transaction ID (Hash)
}

// UTXO represents an unspent transaction output
type UTXO struct {
	TxID   string
	Index  int
	Amount float64
}

// Calculate transaction hash (TxID)
func CalculateTxID(tx Transaction) string {
	record := tx.Sender + tx.Receiver + fmt.Sprintf("%f", tx.Amount)
	hash := sha256.Sum256([]byte(record))
	return hex.EncodeToString(hash[:])
}
