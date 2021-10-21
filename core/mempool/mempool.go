package mempool

import "blockchain-prototype/core/structure/transaction"

type MemPool struct {
	Txns map[transaction.TransactionId]transaction.Transaction
	// OutPointToTxn map[transaction.OutPoint]map[transaction.TransactionId]bool
}

// func (pool *MemPool) InsertTransaction(txn *transaction.Transaction, utxo *UTXOFrontier) {
func (pool *MemPool) InsertTransaction(txn *transaction.Transaction) {
	pool.Txns[transaction.TransactionId(txn.Hash())] = *txn
	// for _, input := range txn.Inputs {
	// 	if _, ok := pool.OutPointToTxn[input]; !ok {
	// 		pool.OutPointToTxn[input] = map[transaction.TransactionId]bool{}
	// 	}
	// 	pool.OutPointToTxn[input][txn.GetId()] = true
	// }
}

func (pool *MemPool) RemoveTransaction(txn *transaction.Transaction) {
	delete(pool.Txns, txn.GetId())
}

func (pool *MemPool) RegisterMined(txns []*transaction.Transaction) {
	for _, txn := range txns {
		pool.RemoveTransaction(txn)
	}
}
