package core

type MemPool struct {
	Txns         map[TransactionId]Transaction
	FeeToTxnList map[Cookies]([]Transaction)
}

// func (pool *MemPool) InsertTransaction(txn *Transaction, utxo *UTXOFrontier) {
func (pool *MemPool) InsertTransaction(txn *Transaction) {
	pool.Txns[TransactionId(txn.ComputeHash())] = *txn

	// fee := utxo.GetTransactionFee(txn)
	// pool.FeeToTxnList[fee] = append(pool.FeeToTxnList[fee], *txn)
}

func (pool *MemPool) GetTransactions(n int, utxo *UTXOFrontier) []Transaction {
	spent := map[OutPoint]bool{}
	txlist := []Transaction{}
	for _, txn := range pool.Txns {
		if len(txlist) >= n {
			break
		}
		if utxo.VerifyTransaction(&txn) {
			valid := true
			for _, outpoint := range txn.TxIn {
				if _, ok := spent[outpoint]; ok {
					valid = false
				}
			}
			if valid {

				txlist = append(txlist, txn)
				for _, outpoint := range txn.TxIn {
					spent[outpoint] = true
				}
			}
		}
	}
	return txlist
}
