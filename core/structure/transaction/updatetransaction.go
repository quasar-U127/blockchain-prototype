package transaction

import (
	"blockchain-prototype/core/utils"
	"fmt"
)

// No Scripts and no contracts what so ever we are only implementing transfer between addresses
type UpdateTransaction struct {
	Original   TransactionId
	Outputs    []Output
	ChainIndex int
}

func (ut *UpdateTransaction) Hash() utils.HashType {
	var byteValues []byte
	byteValues = append(byteValues, ut.Original[:]...)
	for _, output := range ut.Outputs {
		hash := output.Hash()
		bv := [utils.HashSize]byte(hash)
		byteValues = append(byteValues, bv[:]...)
	}
	return utils.Hash(byteValues)
}

func (ut *UpdateTransaction) Print() {
	hash := ut.Hash()
	fmt.Printf("\nName : %x", hash[:2])
	fmt.Printf("\nInputs : ")
	fmt.Printf("\nOutputs : ")
	for _, output := range ut.Outputs {
		fmt.Printf("( %x, %d),", output.Reciever.ShortString(6), output.Value)
	}
	fmt.Print("\n")
}

func UpdateTransactionListHash(utxns []Transaction) utils.HashType {
	var hashs []utils.HashType
	for _, utxn := range utxns {
		hashs = append(hashs, utxn.Hash())
	}
	return utils.MerkleHash(hashs)
}
func (ut UpdateTransaction) GetId() TransactionId {
	return TransactionId(ut.Hash())
}
