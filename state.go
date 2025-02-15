package main

import (
	"fmt"
	"sync"

	"github.com/libp2p/go-libp2p/core/host"
)

// BlockchainState encapsulates all blockchain state
type BlockchainState struct {
	chain      []Block
	pendingTxs []Transaction
	mempool    *Mempool
	wallet     *Wallet
	p2pHost    host.Host
	consensus  *Consensus

	// Mutexes for thread safety
	chainMutex sync.RWMutex
	txMutex    sync.RWMutex
}

func (bs *BlockchainState) ReplaceChain(newChain []Block) {
	bs.chainMutex.Lock()
	defer bs.chainMutex.Unlock()
	bs.chain = newChain
}

// NewBlockchainState initializes a new blockchain state
func NewBlockchainState() *BlockchainState {
	fmt.Println("🔧 Creating new blockchain state...")
	state := &BlockchainState{
		chain:      make([]Block, 0),
		pendingTxs: make([]Transaction, 0),
		mempool:    NewMempool(),
		consensus:  &Consensus{},
	}

	fmt.Println("✨ Blockchain state created successfully")
	return state
}

// Chain operations
// Fix the AddBlock method to avoid deadlock
func (s *BlockchainState) AddBlock(block Block) error {
	fmt.Printf("📦 Adding block %d to chain\n", block.Index)

	// Special case for genesis block
	if len(s.chain) == 0 && block.Index == 0 {
		s.chainMutex.Lock()
		defer s.chainMutex.Unlock()
		s.chain = append(s.chain, block)
		fmt.Println("🌟 Genesis block added successfully")
		return nil
	}

	// Normal block addition logic
	lastBlock := s.GetLastBlock()
	if err := ValidateBlock(block, lastBlock); err != nil {
		return fmt.Errorf("invalid block: %w", err)
	}

	s.chainMutex.Lock()
	defer s.chainMutex.Unlock()

	if len(s.chain) > 0 && s.chain[len(s.chain)-1].Hash != lastBlock.Hash {
		return fmt.Errorf("chain changed during validation")
	}

	s.chain = append(s.chain, block)
	fmt.Printf("✅ Block %d added successfully\n", block.Index)
	return nil
}

func (s *BlockchainState) GetChain() []Block {
	s.chainMutex.RLock()
	defer s.chainMutex.RUnlock()
	return s.chain
}

func (s *BlockchainState) GetLastBlock() Block {
	s.chainMutex.RLock()
	defer s.chainMutex.RUnlock()
	if len(s.chain) == 0 {
		return Block{} // Return genesis block or handle empty chain
	}
	return s.chain[len(s.chain)-1]
}

// Transaction operations
func (s *BlockchainState) AddTransaction(tx Transaction) error {
	s.txMutex.Lock()
	defer s.txMutex.Unlock()

	if err := s.mempool.AddTransaction(tx); err != nil {
		return fmt.Errorf("failed to add transaction: %w", err)
	}
	return nil
}

func (s *BlockchainState) GetPendingTransactions() []Transaction {
	return s.mempool.GetTransactions()
}

// P2P operations
func (s *BlockchainState) SetP2PHost(h host.Host) {
	s.p2pHost = h
}

func (s *BlockchainState) GetP2PHost() host.Host {
	return s.p2pHost
}

// Wallet operations
func (s *BlockchainState) SetWallet(w *Wallet) {
	s.wallet = w
}

func (s *BlockchainState) GetWallet() *Wallet {
	return s.wallet
}

func (s *BlockchainState) GetConsensus() *Consensus {
	return NewConsensus(s)
}
