package core

import (
	"encoding/binary"
	"fmt"
)

type TransactionId HashType

type OutPoint struct {
	Id TransactionId
	N  uint
}

type TXO struct {
	Reciever Address
	Value    Cookies
}

// No Scripts and no contracts what so ever we are only implementing transfer between addresses
type Transaction struct {
	TxIn   []OutPoint
	Output []TXO
}

func (txin OutPoint) ComputeHash() HashType {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(txin.N))
	bv := append([]byte(txin.Id[:]), []byte(b)...)
	return ComputeHash(bv)
}

func (txo TXO) Hash() HashType {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(txo.Value))
	bv := append(txo.Reciever.ByteMarshaling(), b...)
	return ComputeHash(bv)
}

func (t Transaction) ComputeHash() HashType {
	var byteValues []byte
	for _, input := range t.TxIn {
		hash := input.ComputeHash()
		bv := [HashSize]byte(hash)
		byteValues = append(byteValues, bv[:]...)
	}
	for _, output := range t.TxIn {
		hash := output.ComputeHash()
		bv := [HashSize]byte(hash)
		byteValues = append(byteValues, bv[:]...)
	}
	return ComputeHash(byteValues)
}

func (t *Transaction) Print() {
	hash := t.ComputeHash()
	fmt.Printf("\nName : %x", hash[:2])
	fmt.Printf("\nInputs : ")
	for _, input := range t.TxIn {
		fmt.Printf("( %x, %d),", input.Id[:2], input.N)
	}
	fmt.Printf("\nOutputs : ")
	for _, output := range t.Output {
		fmt.Printf("( %x, %d),", output.Reciever.X, output.Value)
	}
	fmt.Print("\n")
}

func ComputeTransactionListHash(txns []Transaction) HashType {
	var byteValues []byte
	for _, txn := range txns {
		b := [HashSize]byte(txn.ComputeHash())
		byteValues = append(byteValues, b[:]...)
	}
	return ComputeHash(byteValues)
}
func (t Transaction) GetId() TransactionId {
	return TransactionId(t.ComputeHash())
}
