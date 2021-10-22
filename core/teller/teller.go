package teller

import (
	"blockchain-prototype/core"
	"blockchain-prototype/core/structure/transaction"
	"blockchain-prototype/core/wallet"
)

type Teller struct {
	utxo      map[transaction.OutPoint]transaction.Output
	committed map[transaction.TransactionId]transaction.Transaction
}

func CreateTeller() Teller {
	return Teller{
		utxo:      map[transaction.OutPoint]transaction.Output{},
		committed: map[transaction.TransactionId]transaction.Transaction{},
	}
}

func (tel *Teller) Update(removed []transaction.Transaction, added []transaction.Transaction) {
	// add := map[transaction.OutPoint]utils.Null{}
	delta := map[transaction.OutPoint]int{}
	for _, txn := range removed {
		id := txn.GetId()
		for i := range txn.Outputs {
			delta[transaction.OutPoint{Id: id, N: uint(i)}]--
		}
		for _, op := range txn.Inputs {
			delta[op]++
		}
		delete(tel.committed, id)
	}
	for _, t := range added {
		id := t.GetId()
		for _, op := range t.Inputs {
			delta[op]--
		}
		for i := range t.Outputs {
			delta[transaction.OutPoint{Id: id, N: uint(i)}]++
		}
		tel.committed[id] = t
	}
	for op, v := range delta {
		switch v {
		case 0:
			continue
		case -1:
			delete(tel.utxo, op)
		case 1:
			{
				t := tel.committed[op.Id]
				tel.utxo[op] = t.Outputs[op.N]
			}
		}
	}
}

func (tel *Teller) Copy() Teller {
	newTeller := Teller{
		utxo:      make(map[transaction.OutPoint]transaction.Output),
		committed: make(map[transaction.TransactionId]transaction.Transaction),
	}
	for outpoint, output := range tel.utxo {
		newTeller.utxo[outpoint] = output
	}
	for txid, txn := range tel.committed {
		newTeller.committed[txid] = txn
	}

	return newTeller
}

func (tel *Teller) GetUTXOs(add *wallet.Address) map[transaction.OutPoint]transaction.Output {
	utxo := map[transaction.OutPoint]transaction.Output{}
	for outpoint, output := range tel.utxo {
		if add == nil || output.Reciever == *add {
			utxo[outpoint] = output
		}
	}
	return utxo
}

func (tel *Teller) GetBalance(add *wallet.Address) core.Cookies {
	var balance core.Cookies = 0
	for _, output := range tel.GetUTXOs(add) {
		balance += output.Value
	}
	return core.Cookies(balance)
}

func (tel *Teller) GetTransaction(txnId transaction.TransactionId) transaction.Transaction {
	return tel.committed[txnId]
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
