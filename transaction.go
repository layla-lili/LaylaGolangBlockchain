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
	Sender    string  // Public Key (Address)
	Receiver  string  // Public Key (Address)
	Amount    float64 // Amount of currency
	TxID      string  // Transaction ID (Hash)
	Signature []byte
	Timestamp time.Time
	Fee       float64
}

// type Transaction struct {
// 	Sender    string `json:"sender"`
// 	Receiver string `json:"receiver"`
// 	Amount    float64 `json:"amount"`
// 	Signature string `json:"signature"`
// }

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

func CalculateBlockReward(block Block) float64 {
	baseReward := 50.0 // Base mining reward
	totalFees := 0.0

	for _, tx := range block.Transactions {
		totalFees += tx.Fee
	}
	return baseReward + totalFees
}


func ValidateTransaction(tx Transaction, publicKey []byte) bool {
	if tx.Amount <= 0 {
		return false
	}

	txHash := sha256.Sum256([]byte(tx.TxID))
	x, y := elliptic.Unmarshal(elliptic.P256(), publicKey)
	if x == nil {
		return false
	}

	publicKeyECDSA := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}

	return ecdsa.VerifyASN1(publicKeyECDSA, txHash[:], tx.Signature)
}
