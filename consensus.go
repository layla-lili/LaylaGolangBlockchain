package main

import (
	"bytes"
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
	for i := 1; i < len(chain); i++ {
		currentBlock := chain[i]
		previousBlock := chain[i-1]

		if !bytes.Equal([]byte(currentBlock.PrevHash), []byte(previousBlock.Hash)) {
			return false
		}

		if !c.ValidateBlockTransactions(currentBlock) {
			return false
		}
	}
	return true
}

func (c *Consensus) ValidateBlockTransactions(block Block) bool {
	for _, tx := range block.Transactions {
		if !c.ValidateTransaction(tx) {
			return false
		}
	}
	return true
}

func (c *Consensus) ValidateTransaction(tx Transaction) bool {
	wallet := c.state.GetWallet()
	publicKeyBytes := wallet.GetPublicKeyBytes()
	return ValidateTransaction(tx, publicKeyBytes)
}

func (c *Consensus) ValidateProofOfWork(block Block) bool {
	prefix := make([]byte, block.Difficulty)
	for i := 0; i < block.Difficulty; i++ {
		prefix[i] = '0'
	}

	hash := CalculateBlockHash(block)
	return bytes.HasPrefix([]byte(hash), prefix)
}
