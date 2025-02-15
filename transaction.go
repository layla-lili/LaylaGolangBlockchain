package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// Transaction represents a simple transaction
type Transaction struct {
	TxID            string
	SenderPublicKey []byte // Added field
	SenderAddress   string // Added field
	Receiver        string
	Amount          float64
	Timestamp       time.Time
	Signature       []byte
	Fee             float64
}

// UTXO represents an unspent transaction output
type UTXO struct {
	TxID   string
	Index  int
	Amount float64
}

// Calculate transaction hash (TxID)
func CalculateTxID(tx Transaction) string {
	record := tx.SenderAddress + tx.Receiver + fmt.Sprintf("%f", tx.Amount)
	hash := sha256.Sum256([]byte(record))
	return hex.EncodeToString(hash[:])
}

func CalculateBlockReward(block Block) float64 {
	baseReward := 50.0 // Base mining reward
	totalFees := 0.0

	for _, tx := range block.Transactions {
		totalFees += tx.Fee
	}
	return baseReward + totalFees
}

func ValidateTransaction(tx Transaction, signature []byte) bool {
	// Verify amount
	if tx.Amount <= 0 {
		return false
	}

	// Reconstruct hash that was signed
	txData := fmt.Sprintf("%s%s%f%v",
		tx.SenderAddress,
		tx.Receiver,
		tx.Amount,
		tx.Timestamp,
	)
	txHash := sha256.Sum256([]byte(txData))

	// Convert public key bytes to ECDSA public key
	x, y := elliptic.Unmarshal(elliptic.P256(), tx.SenderPublicKey)
	if x == nil {
		return false
	}

	publicKey := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}

	return ecdsa.VerifyASN1(publicKey, txHash[:], tx.Signature)
}
