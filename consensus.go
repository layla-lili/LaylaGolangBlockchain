package main

import (
	"bytes"
	"sync"
)

type Consensus struct {
	chainLock sync.RWMutex
}

func (c *Consensus) HandleChainSync(receivedChain []Block) bool {
	c.chainLock.Lock()
	defer c.chainLock.Unlock()

	if !ValidateChain(receivedChain) {
		return false
	}

	if len(receivedChain) <= len(Blockchain) {
		return false
	}

	// Verify work on all blocks
	for _, block := range receivedChain {
		if !ValidateProofOfWork(block) {
			return false
		}
	}

	// Replace our chain
	Blockchain = receivedChain
	return true
}

func ValidateChain(chain []Block) bool {
	for i := 1; i < len(chain); i++ {
		currentBlock := chain[i]
		previousBlock := chain[i-1]

		if !bytes.Equal([]byte(currentBlock.PrevHash), []byte(previousBlock.Hash)) {
			return false
		}

		if !ValidateBlockTransactions(currentBlock) {
			return false
		}
	}
	return true
}

func ValidateBlockTransactions(block Block) bool {
	for _, tx := range block.Transactions {
		if !ValidateTransaction(tx, []byte(tx.Sender)) {
			return false
		}
	}
	return true
}

func ValidateProofOfWork(block Block) bool {
	prefix := make([]byte, block.Difficulty)
	for i := 0; i < block.Difficulty; i++ {
		prefix[i] = '0'
	}

	hash := CalculateBlockHash(block)
	return bytes.HasPrefix([]byte(hash), prefix)
}
