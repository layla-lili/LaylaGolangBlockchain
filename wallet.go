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
	PublicKey  string
}

// Create a new wallet
func NewWallet() *Wallet {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	address := sha256.Sum256(publicKey)

	return &Wallet{
		PrivateKey: privateKey,
		PublicKey:  hex.EncodeToString(address[:]),
	}
}

// Sign a transaction
func (w *Wallet) SignTransaction(tx *Transaction) string {
	txHash := sha256.Sum256([]byte(tx.TxID))
	r, s, _ := ecdsa.Sign(rand.Reader, w.PrivateKey, txHash[:])
	return fmt.Sprintf("%s%s", r.Text(16), s.Text(16))
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
