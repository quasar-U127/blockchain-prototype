package transaction

import (
	"blockchain-prototype/core"
	"blockchain-prototype/core/utils"
	"blockchain-prototype/core/wallet"
	"encoding/binary"
	"fmt"
)

type TransactionId utils.HashType

type OutPoint struct {
	Id TransactionId
	N  uint
}

type Output struct {
	Reciever wallet.Address
	Value    core.Cookies
}

// No Scripts and no contracts what so ever we are only implementing transfer between addresses
type Transaction struct {
	Inputs  []OutPoint
	Outputs []Output
}

func (txin OutPoint) ComputeHash() utils.HashType {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(txin.N))
	bv := append([]byte(txin.Id[:]), []byte(b)...)
	return utils.Hash(bv)
}

func (txo Output) Hash() utils.HashType {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(txo.Value))
	bv := append(txo.Reciever.Marshal(), b...)
	return utils.Hash(bv)
}

func (t *Transaction) Hash() utils.HashType {
	var byteValues []byte
	for _, input := range t.Inputs {
		hash := input.ComputeHash()
		bv := [utils.HashSize]byte(hash)
		byteValues = append(byteValues, bv[:]...)
	}
	for _, output := range t.Inputs {
		hash := output.ComputeHash()
		bv := [utils.HashSize]byte(hash)
		byteValues = append(byteValues, bv[:]...)
	}
	return utils.Hash(byteValues)
}

func (t *Transaction) Print() {
	hash := t.Hash()
	fmt.Printf("\nName : %x", hash[:2])
	fmt.Printf("\nInputs : ")
	for _, input := range t.Inputs {
		fmt.Printf("( %x, %d),", input.Id[:2], input.N)
	}
	fmt.Printf("\nOutputs : ")
	for _, output := range t.Outputs {
		fmt.Printf("( %x, %d),", output.Reciever.ShortString(6), output.Value)
	}
	fmt.Print("\n")
}

func TransactionListHash(txns []Transaction) utils.HashType {
	var byteValues []byte
	for _, txn := range txns {
		b := [utils.HashSize]byte(txn.Hash())
		byteValues = append(byteValues, b[:]...)
	}
	return utils.Hash(byteValues)
}
func (t Transaction) GetId() TransactionId {
	return TransactionId(t.Hash())
}
