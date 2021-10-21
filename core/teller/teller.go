package teller

import (
	"blockchain-prototype/core"
	"blockchain-prototype/core/structure/transaction"
	"blockchain-prototype/core/wallet"
)

type Teller struct {
	utxo map[transaction.OutPoint]transaction.Output
}

func (tel *Teller) GetUTXOs(add *wallet.Address) map[transaction.OutPoint]transaction.Output {
	utxo := map[transaction.OutPoint]transaction.Output{}
	for outpoint, output := range tel.utxo {
		utxo[outpoint] = output
	}
	return utxo
}

func (tel *Teller) GetBalance(add *wallet.Address) core.Cookies {
	balance := 0
	for _, output := range tel.GetUTXOs(add) {
		balance += int(output.Value)
	}
	return core.Cookies(balance)
}

func (tel *Teller) GetTransactionFee(txn *transaction.Transaction) core.Cookies {
	var output core.Cookies = 0
	for _, txo := range txn.Outputs {
		output += txo.Value
	}
	var input core.Cookies = 0
	for _, outpoint := range txn.Inputs {
		input += tel.utxo[outpoint].Value
	}
	return input - output
}
func (tel *Teller) GetTransactionListFee(txns []transaction.Transaction) core.Cookies {
	fees := 0
	for _, txn := range txns {
		fees += int(tel.GetTransactionFee(&txn))
	}
	return core.Cookies(fees)
}
