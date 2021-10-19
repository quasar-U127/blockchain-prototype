package core

type MemPool struct {
	Txns map[TransactionId]Transaction
}

// func (pool *MemPool) InsertTransaction(txn *Transaction, utxo *UTXOFrontier) {
func (pool *MemPool) InsertTransaction(txn *Transaction) {
	pool.Txns[TransactionId(txn.ComputeHash())] = *txn
}

func (pool *MemPool) RemoveTransaction(txn *Transaction) {
	delete(pool.Txns, txn.GetId())
}
