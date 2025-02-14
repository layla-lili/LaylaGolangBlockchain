package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

// Wallet structure
type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  []byte
	Address    string
	UTXOs      []UTXO // <- Adding this field to keep track of UTXOs
}

// Create a new wallet
func NewWallet() (*Wallet, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate wallet: %v", err)
	}

	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	address := generateAddress(publicKey)

	return &Wallet{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Address:    address,
		UTXOs:      []UTXO{}, // Initialize the UTXOs slice
	}, nil
}

func generateAddress(publicKey []byte) string {
	hash := sha256.Sum256(publicKey)
	return hex.EncodeToString(hash[:])
}

// Sign a transaction
func (w *Wallet) SignTransaction(tx *Transaction) error {
	txHash := sha256.Sum256([]byte(tx.TxID))
	signature, err := ecdsa.SignASN1(rand.Reader, w.PrivateKey, txHash[:])
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %v", err)
	}
	tx.Signature = signature
	return nil
}

// Verify transaction signature
func VerifySignature(tx *Transaction, signature string, pubKey ecdsa.PublicKey) bool {
	txHash := sha256.Sum256([]byte(tx.TxID))
	r := new(big.Int)
	s := new(big.Int)
	sigLen := len(signature) / 2
	r.SetString(signature[:sigLen], 16)
	s.SetString(signature[sigLen:], 16)

	return ecdsa.Verify(&pubKey, txHash[:], r, s)
}

func validateTransaction(tx Transaction, wallet *Wallet) bool {
	// Check if the sender has enough funds by validating against the UTXO set
	for _, utxo := range wallet.UTXOs {
		if utxo.TxID == tx.TxID && utxo.Amount >= tx.Amount {
			return true
		}
	}
	return false
}

func removeSpentTransactions(tx Transaction) {
	// Loop over pending transactions and remove spent ones
	for i, t := range pendingTransactions {
		if t.TxID == tx.TxID {
			// Remove from pending transactions
			pendingTransactions = append(pendingTransactions[:i], pendingTransactions[i+1:]...)
			break
		}
	}
}
