package teller

import (
	"blockchain-prototype/core"
	"blockchain-prototype/core/structure/transaction"
)

func (tel *Teller) VerifyTransaction(txn *transaction.Transaction) bool {
	var output core.Cookies = 0
	for _, txo := range txn.Outputs {
		output += txo.Value
	}
	var input core.Cookies = 0
	for _, outpoint := range txn.Inputs {
		v, ok := tel.utxo[outpoint]
		if !ok {
			return false
		}
		input += v.Value
	}
	return input >= output
}

func (tel *Teller) VerifyTransactionDelta(txn *transaction.Transaction, spent map[transaction.OutPoint]bool) bool {
	if !tel.VerifyTransaction(txn) {
		return false
	}
	for _, input := range txn.Inputs {
		if _, found := spent[input]; found {
			return false
		}
	}
	return true

}

func (tel *Teller) VerifyTransactionList(txns []transaction.Transaction) bool {

	spent := map[transaction.OutPoint]bool{}
	for _, txn := range txns {
		if tel.VerifyTransactionDelta(&txn, spent) {
			for _, input := range txn.Inputs {
				spent[input] = true
			}
		} else {
			return false
		}
	}
	return true

}

func (tel *Teller) CommitTransaction(txn *transaction.Transaction) {
	for _, input := range txn.Inputs {
		delete(tel.utxo, input)
	}
	id := txn.GetId()
	for i, output := range txn.Outputs {
		tel.utxo[transaction.OutPoint{Id: id, N: uint(i)}] = output
	}
}

func (tel *Teller) CommitTransactionList(txns []transaction.Transaction) {
	for _, txn := range txns {
		tel.CommitTransaction(&txn)
	}
}
