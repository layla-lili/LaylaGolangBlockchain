package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"fmt"
	"sync"
)

type Consensus struct {
	state     *BlockchainState
	chainLock sync.RWMutex
}

func NewConsensus(state *BlockchainState) *Consensus {
	return &Consensus{
		state: state,
	}
}

func (c *Consensus) HandleChainSync(receivedChain []Block) bool {
	c.chainLock.Lock()
	defer c.chainLock.Unlock()

	if !c.ValidateChain(receivedChain) {
		return false
	}

	currentChain := c.state.GetChain()
	if len(receivedChain) <= len(currentChain) {
		return false
	}

	// Verify work on all blocks
	for _, block := range receivedChain {
		if !c.ValidateProofOfWork(block) {
			return false
		}
	}

	// Replace our chain
	c.state.ReplaceChain(receivedChain)
	return true
}

func (c *Consensus) ValidateChain(chain []Block) bool {
	if len(chain) == 0 {
		return false
	}

	// Validate each block
	for i := 1; i < len(chain); i++ {
		block := chain[i]
		prevBlock := chain[i-1]

		// Validate block hash and previous hash
		if block.PrevHash != prevBlock.Hash {
			fmt.Printf("❌ Invalid previous hash at block %d\n", block.Index)
			return false
		}

		calculatedHash := CalculateBlockHash(block)
		if calculatedHash != block.Hash {
			fmt.Printf("❌ Invalid block hash at block %d\n", block.Index)
			return false
		}

		// Validate block transactions
		for _, tx := range block.Transactions {
			if !c.ValidateTransaction(tx) {
				fmt.Printf("❌ Invalid transaction in block %d\n", block.Index)
				return false
			}
		}
	}

	return true
}

func (c *Consensus) ValidateTransaction(tx Transaction) bool {
	// Reconstruct public key from bytes
	x, y := elliptic.Unmarshal(elliptic.P256(), tx.SenderPublicKey)
	if x == nil {
		fmt.Printf("❌ Invalid public key format\n")
		return false
	}

	pubKey := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}

	// Calculate transaction hash using the same method as signing
	txData := fmt.Sprintf("%s%s%f%v",
		tx.SenderAddress,
		tx.Receiver,
		tx.Amount,
		tx.Timestamp.Unix())

	txHash := sha256.Sum256([]byte(txData))

	// Verify signature
	if !ecdsa.VerifyASN1(pubKey, txHash[:], tx.Signature) {
		fmt.Printf("❌ Invalid signature for transaction %s\n", tx.TxID)
		return false
	}

	// Verify address matches public key
	derivedAddr := GenerateAddress(tx.SenderPublicKey)
	if derivedAddr != tx.SenderAddress {
		fmt.Printf("❌ Address mismatch: derived=%s, tx=%s\n", derivedAddr, tx.SenderAddress)
		return false
	}

	return true
}

func (c *Consensus) ValidateProofOfWork(block Block) bool {
	prefix := make([]byte, block.Difficulty)
	for i := 0; i < block.Difficulty; i++ {
		prefix[i] = '0'
	}

	hash := CalculateBlockHash(block)
	return bytes.HasPrefix([]byte(hash), prefix)
}
