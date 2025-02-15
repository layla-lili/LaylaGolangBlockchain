package main

import (
	"fmt"
	"sync"
	"time"
)

type Mempool struct {
	transactions map[string]Transaction
	mutex        sync.RWMutex
}

func NewMempool() *Mempool {
	return &Mempool{
		transactions: make(map[string]Transaction),
	}
}

func (m *Mempool) AddTransaction(tx Transaction) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if !ValidateTransaction(tx, []byte(tx.SenderAddress)) {
		return fmt.Errorf("invalid transaction")
	}

	m.transactions[tx.TxID] = tx
	return nil
}

func (m *Mempool) GetTransactions() []Transaction {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	txs := make([]Transaction, 0, len(m.transactions))
	for _, tx := range m.transactions {
		txs = append(txs, tx)
	}
	return txs
}

func (m *Mempool) RemoveTransactions(txs []Transaction) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, tx := range txs {
		delete(m.transactions, tx.TxID)
	}
}

// Cleanup old transactions
func (m *Mempool) CleanupOldTransactions(maxAge time.Duration) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	now := time.Now()
	for txID, tx := range m.transactions {
		if now.Sub(tx.Timestamp) > maxAge {
			delete(m.transactions, txID)
		}
	}
}
