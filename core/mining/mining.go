package mining

import (
	"blockchain-prototype/core"
	"blockchain-prototype/core/mempool"
	"blockchain-prototype/core/structure/block"
	"blockchain-prototype/core/structure/transaction"
	"blockchain-prototype/core/teller"
	"blockchain-prototype/core/utils"
	"blockchain-prototype/core/verifier"
	"sort"
)

type feePair struct {
	fee core.Cookies
	id  transaction.TransactionId
}

type sortByValue []feePair

func (a sortByValue) Len() int           { return len(a) }
func (a sortByValue) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortByValue) Less(i, j int) bool { return a[i].fee > a[j].fee }

func MineTransactions(pool *mempool.MemPool, tel *teller.Teller, n int) []transaction.Transaction {
	spent := map[transaction.OutPoint]bool{}
	txlist := []transaction.Transaction{}

	feeTxn := []feePair{}

	for _, txn := range pool.Txns {
		if tel.VerifyTransaction(&txn) {
			feeTxn = append(feeTxn, feePair{fee: tel.GetTransactionFee(&txn), id: txn.GetId()})
		}
	}

	sort.Sort(sortByValue(feeTxn))

	for _, txnid := range feeTxn {
		txn := pool.Txns[txnid.id]
		if len(txlist) >= n {
			break
		}

		if tel.VerifyTransactionDelta(&txn, spent) {

			txlist = append(txlist, txn)
			for _, outpoint := range txn.Inputs {
				spent[outpoint] = true
			}
		}
	}
	return txlist
}
func MineBlock(txns []transaction.Transaction, prevBlockHash utils.HashType, height uint) *block.Block {
	tries := uint(100000)
	block := block.Create(txns, prevBlockHash, height, 0)
	for i := uint(0); i < tries; i++ {
		block.Header.Nonce = i
		hash := block.Hash()
		if verifier.VerifyDifficulty(hash) {
			return &block
		}
	}
	return nil
}
