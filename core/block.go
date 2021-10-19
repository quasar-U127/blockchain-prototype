package core

import (
	"encoding/binary"
)

type BlockHeader struct {
	PrevBlock HashType
	TxnHash   HashType
	Height    uint
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

	byteHeight := make([]byte, 4)
	binary.LittleEndian.PutUint32(byteHeight, uint32(bh.Height))
	byteValues = append(byteValues, byteHeight[:]...)

	return ComputeHash(byteValues)

}

func CreateBlock(txns []Transaction, prevBlockHash HashType, height uint, nonce uint) Block {

	header := BlockHeader{
		PrevBlock: prevBlockHash,
		TxnHash:   ComputeTransactionListHash(txns),
		Nonce:     nonce,
		Height:    height,
	}
	return Block{Header: header, Txns: txns}
}

func (b Block) ComputeHash() HashType {
	return b.Header.ComputeHash()
}

func ValidDifficulty(hash HashType) bool {
	return hash[0] == 0 && hash[1] < 16
}

func MineBlock(txns []Transaction, prevBlockHash HashType, height uint) *Block {
	tries := uint(100000)
	block := CreateBlock(txns, prevBlockHash, height, 0)
	for i := uint(0); i < tries; i++ {
		block.Header.Nonce = i
		hash := block.ComputeHash()
		if ValidDifficulty(hash) {
			return &block
		}
	}
	return nil
}
