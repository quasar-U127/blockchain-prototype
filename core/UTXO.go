package core

type TXOSet struct {
	Set map[OutPoint]TXO
}
type UTXOFrontier struct {
	Frontier TXOSet
}

func CreateUTXOFrontier() UTXOFrontier {
	return UTXOFrontier{
		Frontier: TXOSet{Set: map[OutPoint]TXO{}},
	}
}

func (utxof *UTXOFrontier) VerifyTransaction(txn *Transaction) bool {
	var output Cookies = 0
	for _, txo := range txn.Output {
		output += txo.Value
	}
	var input Cookies = 0
	for _, txin := range txn.TxIn {
		v, ok := utxof.Frontier.Set[txin]
		if !ok {
			return false
		}
		input += v.Value
	}
	return input >= output
}

func (utxof *UTXOFrontier) GetTransactionFee(txn *Transaction) Cookies {
	var output Cookies = 0
	for _, txo := range txn.Output {
		output += txo.Value
	}
	var input Cookies = 0
	for _, txin := range txn.TxIn {
		v := utxof.Frontier.Set[txin]
		input += v.Value
	}
	return input - output
}

func (utxof *UTXOFrontier) VerifyTransactionList(txns []Transaction) bool {
	for _, txn := range txns {
		if !utxof.VerifyTransaction(&txn) {
			return false
		}
	}
	return true
}
func (utxof *UTXOFrontier) GetFees(txns []Transaction) Cookies {
	var fees Cookies = 0
	for _, txn := range txns {
		fees += utxof.GetTransactionFee(&txn)
	}
	return fees
}

func (utxof *UTXOFrontier) Update(txns []Transaction) {
	for _, txn := range txns {
		for _, outpoint := range txn.TxIn {
			if outpoint.Id == TransactionId(ZeroHash()) {
				continue
			}
			delete(utxof.Frontier.Set, outpoint)
		}
		hash := txn.ComputeHash()
		for i, txo := range txn.Output {
			utxof.Frontier.Set[OutPoint{Id: TransactionId(hash), N: uint(i)}] = txo
		}
	}
}
