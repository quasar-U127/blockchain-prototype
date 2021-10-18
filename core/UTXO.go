package core

type UTXOFrontier struct {
	Frontier map[OutPoint]Cookies
}

func (utxof *UTXOFrontier) VerifyTransaction(txn *Transaction) bool {
	var output Cookies = 0
	for _, txo := range txn.Output {
		output += txo.Value
	}
	var input Cookies = 0
	for _, txin := range txn.TxIn {
		v, ok := utxof.Frontier[txin]
		if !ok {
			return false
		}
		input += v
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
		v := utxof.Frontier[txin]
		input += v
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
