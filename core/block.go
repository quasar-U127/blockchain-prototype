package core

import (
	"encoding/binary"
	"fmt"
)

type BlockHeader struct {
	PrevBlock HashType
	TxnHash   HashType
	Nonce     uint
}
type Block struct {
	Header BlockHeader
	// We can create a merkel a tree later on but not now
	Txns []Transaction
}

func (bh BlockHeader) ComputeHash() HashType {
	var byteValues []byte

	byteSegment := [HashSize]byte(bh.PrevBlock)
	byteValues = append(byteValues, byteSegment[:]...)

	byteSegment = [HashSize]byte(bh.TxnHash)
	byteValues = append(byteValues, byteSegment[:]...)

	byteNonce := make([]byte, 4)
	binary.LittleEndian.PutUint32(byteNonce, uint32(bh.Nonce))
	byteValues = append(byteValues, byteNonce[:]...)

	return ComputeHash(byteValues)

}

func CreateBlock(txns []Transaction, prevBlockHash HashType, nonce uint) Block {

	header := BlockHeader{
		PrevBlock: prevBlockHash,
		TxnHash:   ComputeTransactionListHash(txns),
		Nonce:     nonce,
	}
	return Block{Header: header, Txns: txns}
}

func (b Block) ComputeHash() HashType {
	return b.Header.ComputeHash()
}

func MineBlock(txns []Transaction, prevBlockHash HashType) {
	tries := uint(1000)
	header := BlockHeader{
		PrevBlock: prevBlockHash,
		TxnHash:   ComputeTransactionListHash(txns),
		Nonce:     0,
	}
	for i := uint(0); i < tries; i++ {
		header.Nonce = i
		hash := header.ComputeHash()
		if hash[0] == 0 {
			fmt.Printf("%x\n", [HashSize]byte(hash))
		}
	}
}
