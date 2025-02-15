package main

import (
	// "crypto/ecdh"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

// Wallet represents a cryptocurrency wallet
type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  []byte
	Address    string
	UTXOs      []UTXO
}

// NewWallet creates and returns a new Wallet instance
func NewWallet() (*Wallet, error) {
	fmt.Println("Creating new wallet...")
	privateKey, err := generatePrivateKey()
	if err != nil {
		return nil, err
	}

	publicKey := generatePublicKey(privateKey)
	address := generateAddress(publicKey)

	return &Wallet{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Address:    address,
		UTXOs:      make([]UTXO, 0),
	}, nil
}

// generatePrivateKey creates a new ECDSA private key
func generatePrivateKey() (*ecdsa.PrivateKey, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

// generatePublicKey derives the public key from the private key
func generatePublicKey(privateKey *ecdsa.PrivateKey) []byte {
	return elliptic.Marshal(privateKey.Curve, privateKey.PublicKey.X, privateKey.PublicKey.Y)
}

// generateAddress creates a wallet address from the public key
func generateAddress(publicKey []byte) string {
	// Step 1: SHA-256 hash of public key
	sha256Hash := sha256.Sum256(publicKey)

	// Step 2: RIPEMD-160 hash of the SHA-256 result
	ripemd160Hasher := ripemd160.New()
	_, err := ripemd160Hasher.Write(sha256Hash[:])
	if err != nil {
		return ""
	}
	publicKeyHash := ripemd160Hasher.Sum(nil)

	// Step 3: Add version byte in front (0x00 for mainnet)
	versionedPayload := append([]byte{0x00}, publicKeyHash...)

	// Step 4: Double SHA-256 for checksum
	firstHash := sha256.Sum256(versionedPayload)
	secondHash := sha256.Sum256(firstHash[:])
	checksum := secondHash[:4]

	// Step 5: Concatenate versioned payload and checksum
	finalBytes := append(versionedPayload, checksum...)

	// Step 6: Base58 encode the result
	address := base58.Encode(finalBytes)

	return address
}

func GenerateAddress(publicKey []byte) string {
	// Step 1: SHA-256 hash of public key
	sha256Hash := sha256.Sum256(publicKey)

	// Step 2: RIPEMD-160 hash of the SHA-256 result
	ripemd160Hasher := ripemd160.New()
	if _, err := ripemd160Hasher.Write(sha256Hash[:]); err != nil {
		return ""
	}
	publicKeyHash := ripemd160Hasher.Sum(nil)

	// Step 3: Add version byte in front (0x00 for mainnet)
	versionedPayload := append([]byte{0x00}, publicKeyHash...)

	// Step 4: Double SHA-256 for checksum
	firstHash := sha256.Sum256(versionedPayload)
	secondHash := sha256.Sum256(firstHash[:])

	// Step 5: Add 4-byte checksum to versioned payload
	finalBytes := append(versionedPayload, secondHash[:4]...)

	// Step 6: Base58 encode the result
	return base58.Encode(finalBytes)
}

func (w *Wallet) GetPublicKeyBytes() []byte {
	if w == nil || w.PublicKey == nil {
		return []byte{} // Return empty bytes instead of nil
	}

	pubKeyBytes := w.PublicKey
	if pubKeyBytes == nil {
		return []byte{} // Ensure we never return nil
	}

	return pubKeyBytes
}

func (w *Wallet) GetAddress() string {
	return w.Address
}

// Sign a transaction
func (w *Wallet) SignTransaction(tx *Transaction) error {
	if w == nil || w.PrivateKey == nil {
		return fmt.Errorf("wallet or private key is nil")
	}

	// Set sender information first
	tx.SenderAddress = w.GetAddress()
	tx.SenderPublicKey = w.GetPublicKeyBytes()

	// Calculate transaction hash for signing
	txData := fmt.Sprintf("%s%s%f%v",
		tx.SenderAddress,
		tx.Receiver,
		tx.Amount,
		tx.Timestamp.Unix())

	txHash := sha256.Sum256([]byte(txData))

	// Sign the transaction hash
	signature, err := ecdsa.SignASN1(rand.Reader, w.PrivateKey, txHash[:])
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Set signature and TxID
	tx.Signature = signature
	tx.TxID = hex.EncodeToString(txHash[:])

	return nil
}

// Verify transaction signature
func VerifyTransactionSignature(tx *Transaction, signature string, pubKey ecdsa.PublicKey) bool {
	txHash := sha256.Sum256([]byte(tx.TxID))
	r := new(big.Int)
	s := new(big.Int)
	sigLen := len(signature) / 2
	r.SetString(signature[:sigLen], 16)
	s.SetString(signature[sigLen:], 16)

	return ecdsa.Verify(&pubKey, txHash[:], r, s)
}

func SignData(privateKey *ecdsa.PrivateKey, data []byte) ([]byte, error) {
	hash := sha256.Sum256(data)
	return ecdsa.SignASN1(rand.Reader, privateKey, hash[:])
}

func VerifySignature(publicKeyBytes []byte, data []byte, signature []byte) bool {
	// Improvements:
	// 1. Takes raw bytes for all parameters
	// 2. More flexible - can verify any data, not just transactions
	// 3. Reconstructs public key from bytes
	// 4. Uses ASN.1 encoded signatures

	// Convert public key bytes to ECDSA public key
	x, y := elliptic.Unmarshal(elliptic.P256(), publicKeyBytes)
	if x == nil {
		return false
	}

	publicKey := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}

	return ecdsa.VerifyASN1(publicKey, data, signature)
}

// Verify transaction signature
// func VerifySignature(tx *Transaction, signature string, pubKey ecdsa.PublicKey) bool {
// 	txHash := sha256.Sum256([]byte(tx.TxID))
// 	r := new(big.Int)
// 	s := new(big.Int)
// 	sigLen := len(signature) / 2
// 	r.SetString(signature[:sigLen], 16)
// 	s.SetString(signature[sigLen:], 16)

// 	return ecdsa.Verify(&pubKey, txHash[:], r, s)
// }

// func validateTransaction(tx Transaction, wallet *Wallet) bool {
// 	// Check if the sender has enough funds by validating against the UTXO set
// 	for _, utxo := range wallet.UTXOs {
// 		if utxo.TxID == tx.TxID && utxo.Amount >= tx.Amount {
// 			return true
// 		}
// 	}
// 	return false
// }

// func removeSpentTransactions(tx Transaction) {
// 	// Loop over pending transactions and remove spent ones
// 	for i, t := range pendingTransactions {
// 		if t.TxID == tx.TxID {
// 			// Remove from pending transactions
// 			pendingTransactions = append(pendingTransactions[:i], pendingTransactions[i+1:]...)
// 			break
// 		}
// 	}
// }
