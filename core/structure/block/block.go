package block

import (
	"blockchain-prototype/core/structure/transaction"
	"blockchain-prototype/core/utils"
	"encoding/binary"
)

type BlockHeader struct {
	PrevBlock utils.HashType
	TxnHash   utils.HashType
	Height    uint
	Nonce     uint
}
type Block struct {
	Header BlockHeader
	// We can create a merkel a tree later on but not now
	Txns []transaction.Transaction
}
type BlockId utils.HashType

const BlockLimit int = 10

func (bh *BlockHeader) Hash() utils.HashType {
	var byteValues []byte

	byteSegment := [utils.HashSize]byte(bh.PrevBlock)
	byteValues = append(byteValues, byteSegment[:]...)

	byteSegment = [utils.HashSize]byte(bh.TxnHash)
	byteValues = append(byteValues, byteSegment[:]...)

	byteNonce := make([]byte, 4)
	binary.LittleEndian.PutUint32(byteNonce, uint32(bh.Nonce))
	byteValues = append(byteValues, byteNonce[:]...)

	byteHeight := make([]byte, 4)
	binary.LittleEndian.PutUint32(byteHeight, uint32(bh.Height))
	byteValues = append(byteValues, byteHeight[:]...)

	return utils.Hash(byteValues)

}

func Create(txns []transaction.Transaction, prevBlockHash utils.HashType, height uint, nonce uint) Block {

	header := BlockHeader{
		PrevBlock: prevBlockHash,
		TxnHash:   transaction.TransactionListHash(txns),
		Nonce:     nonce,
		Height:    height,
	}
	return Block{Header: header, Txns: txns}
}

func CreateGenesis() Block {
	return Create([]transaction.Transaction{}, utils.ZeroHash(), 1, 0)
}

func CreateHeader(txnHash utils.HashType, prevBlockHash utils.HashType, height uint, nonce uint) BlockHeader {
	return BlockHeader{
		PrevBlock: prevBlockHash,
		TxnHash:   txnHash,
		Nonce:     nonce,
		Height:    height,
	}
}

func (b *Block) Hash() utils.HashType {
	return b.Header.Hash()
}

func (b *Block) Id() BlockId { return BlockId(b.Hash()) }
