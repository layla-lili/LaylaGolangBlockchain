package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"crypto/ecdsa"
	"crypto/elliptic"
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

func ValidateTransaction(tx Transaction, pubKeyBytes []byte) bool {
	// Reconstruct public key
	x, y := elliptic.Unmarshal(elliptic.P256(), pubKeyBytes)
	if x == nil {
		return false
	}

	pubKey := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}

	// Verify address matches public key
	derivedAddr := GenerateAddress(pubKeyBytes)
	if derivedAddr != tx.SenderAddress {
		return false
	}

	// Calculate transaction hash (same as signing)
	txData := fmt.Sprintf("%s%s%f%v",
		tx.SenderAddress,
		tx.Receiver,
		tx.Amount,
		tx.Timestamp.Unix())

	txHash := sha256.Sum256([]byte(txData))

	// Verify signature
	return ecdsa.VerifyASN1(pubKey, txHash[:], tx.Signature)
}

// func ValidateTransaction(tx Transaction, signature []byte) bool {
// 	// Verify amount
// 	if tx.Amount <= 0 {
// 		return false
// 	}

// 	// Reconstruct hash that was signed
// 	txData := fmt.Sprintf("%s%s%f%v",
// 		tx.SenderAddress,
// 		tx.Receiver,
// 		tx.Amount,
// 		tx.Timestamp,
// 	)
// 	txHash := sha256.Sum256([]byte(txData))

// 	// Convert public key bytes to ECDSA public key
// 	x, y := elliptic.Unmarshal(elliptic.P256(), tx.SenderPublicKey)
// 	if x == nil {
// 		return false
// 	}

// 	publicKey := &ecdsa.PublicKey{
// 		Curve: elliptic.P256(),
// 		X:     x,
// 		Y:     y,
// 	}

// 	return ecdsa.VerifyASN1(publicKey, txHash[:], tx.Signature)
// }

// func (c *Consensus) ValidateTransaction(tx Transaction) bool {
// 	// 1. Check public key presence
// 	if len(tx.SenderPublicKey) == 0 {
// 		return false
// 	}

// 	// 2. Verify signature
// 	txHash := CalculateTxID(tx)
// 	txHashBytes, _ := hex.DecodeString(txHash)
// 	if !VerifySignature(tx.SenderPublicKey, txHashBytes, tx.Signature) {
// 		return false
// 	}

// 	// 3. Verify address derivation
// 	if GenerateAddress(tx.SenderPublicKey) != tx.SenderAddress {
// 		return false
// 	}

// 	return true
// }
