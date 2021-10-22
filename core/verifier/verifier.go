package verifier

import (
	"blockchain-prototype/core"
	"blockchain-prototype/core/blockchain"
	"blockchain-prototype/core/structure/block"
	"blockchain-prototype/core/structure/transaction"
	"blockchain-prototype/core/teller"
	"blockchain-prototype/core/utils"
)

type Verifier struct {
	Reward core.Cookies
}

func CreateVerifier() Verifier {
	return Verifier{Reward: 50}
}

func VerifyDifficulty(hash utils.HashType) bool {
	return hash[0] == 0 && hash[1] < 16
}

func (ver *Verifier) VerifyTransaction(txn *transaction.Transaction, tel *teller.Teller) bool {
	return tel.VerifyTransaction(txn)
}

func (ver *Verifier) VerifyTransactionDelta(txn *transaction.Transaction, spent map[transaction.OutPoint]bool, tel *teller.Teller) bool {

	return tel.VerifyTransactionDelta(txn, spent)

}

func (ver *Verifier) VerifyTransactionList(txns []transaction.Transaction, tel *teller.Teller) bool {

	return tel.VerifyTransactionList(txns)

}

func (ver *Verifier) VerifyBlock(blk *block.Block, chain *blockchain.BlockChain) bool {
	prevId := blk.Header.PrevBlock
	prevBlock := chain.GetBlock(block.BlockId(prevId))
	correct := prevBlock != nil && prevBlock.Header.Height+1 == blk.Header.Height
	if !correct {
		return false
	}
	tel := chain.GetTeller(block.BlockId(prevId))
	if !ver.VerifyTransactionList(blk.Txns[1:], tel) {
		return false
	}
	fees := tel.GetTransactionListFee(blk.Txns[1:])
	coinbaseSum := 0
	for _, o := range blk.Txns[0].Outputs {
		coinbaseSum += int(o.Value)
	}
	if coinbaseSum > int(ver.Reward)+int(fees) {
		return false
	}
	if transaction.TransactionListHash(blk.Txns) != blk.Header.TxnHash {
		return false
	}

	return VerifyDifficulty(blk.Hash())
}
