package verifier

import (
	"blockchain-prototype/core/structure/block"
	"blockchain-prototype/core/structure/transaction"
	"blockchain-prototype/core/teller"
	"blockchain-prototype/core/utils"
)

type Verifier struct {
	tel *teller.Teller
}

func VerifyDifficulty(hash utils.HashType) bool {
	return hash[0] == 0 && hash[1] < 16
}

func (ver *Verifier) VerifyTransaction(txn *transaction.Transaction) bool {
	return ver.tel.VerifyTransaction(txn)
}

func (ver *Verifier) VerifyTransactionDelta(txn *transaction.Transaction, spent map[transaction.OutPoint]bool) bool {

	return ver.tel.VerifyTransactionDelta(txn, spent)

}

func (ver *Verifier) VerifyTransactionList(txns []transaction.Transaction) bool {

	return ver.tel.VerifyTransactionList(txns)

}

func (ver *Verifier) VerifyBlock(blk *block.Block) bool {
	if !ver.VerifyTransactionList(blk.Txns) {
		return false
	}
	if transaction.TransactionListHash(blk.Txns) != blk.Header.TxnHash {
		return false
	}

	return VerifyDifficulty(blk.Hash())
}

func (ver *Verifier) CommitTransaction(txn *transaction.Transaction) {
	ver.tel.CommitTransaction(txn)
}

func (ver *Verifier) CommitTransactionList(txns []transaction.Transaction) {
	ver.tel.CommitTransactionList(txns)
}

func (ver *Verifier) CommitBlock(blk *block.Block) {
	ver.CommitTransactionList(blk.Txns)
}
