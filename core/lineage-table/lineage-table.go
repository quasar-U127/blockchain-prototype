package lineagetable

import "blockchain-prototype/core/structure/transaction"

type LineageTable struct {
}

func CreateLineageTable() LineageTable {
	return LineageTable{}
}

func (lt *LineageTable) StoreTransaction(txn *transaction.Transaction) bool {
	return false
}

func (lt *LineageTable) FetchTransaction(txnHash transaction.TransactionId) transaction.Transaction
func (lt *LineageTable) FetchLatestTransaction(txnHash transaction.TransactionId) transaction.Transaction
