package mining

import (
	"blockchain-prototype/core"
	"blockchain-prototype/core/blockchain"
	"blockchain-prototype/core/mempool"
	"blockchain-prototype/core/structure/block"
	"blockchain-prototype/core/structure/transaction"
	"blockchain-prototype/core/teller"
	"blockchain-prototype/core/utils"
	"blockchain-prototype/core/verifier"
	"blockchain-prototype/core/wallet"
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
func TuneBlock(blk *block.Block, ver *verifier.Verifier) bool {
	tries := uint(100000)
	for i := uint(0); i < tries; i++ {
		blk.Header.Nonce = i
		hash := blk.Hash()
		if verifier.VerifyDifficulty(hash) {
			return true
		}
	}
	return false
}

func MineBlock(chain *blockchain.BlockChain, add *wallet.Address, ver *verifier.Verifier, pool *mempool.MemPool) *block.Block {
	tel := chain.GetTeller(chain.End())
	txns := []transaction.Transaction{}
	minedTxns := MineTransactions(pool, tel, block.BlockLimit-1)
	fees := tel.GetTransactionListFee(minedTxns)
	coinBaseTxn := transaction.Transaction{
		Inputs: []transaction.OutPoint{{Id: transaction.TransactionId(utils.ZeroHash()), N: chain.Height() + 1}},
		Outputs: []transaction.Output{
			{Reciever: *add, Value: 50 + fees},
		},
	}
	txns = append(txns, coinBaseTxn)
	txns = append(txns, minedTxns...)
	prevHash := chain.End()
	newBlock := block.Create(txns, utils.HashType(prevHash), chain.Height()+1, 0)
	tuned := TuneBlock(&newBlock, ver)
	// block := MineBlock(txns, prevHash, state.Chain.N+1)
	if tuned {
		return &newBlock
	} else {
		return nil
	}
}
