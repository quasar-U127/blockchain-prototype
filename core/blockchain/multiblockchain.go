package blockchain

import (
	"blockchain-prototype/core/structure/block"
	"blockchain-prototype/core/structure/transaction"
	"blockchain-prototype/core/teller"
)

type MultiBlockchain struct {
	store  block.Store
	chains []BlockChain
}

func CreateMultiBlockchain(p uint) MultiBlockchain {
	n := MultiBlockchain{
		store:  block.CreateStore(),
		chains: []BlockChain{},
	}
	for i := 0; i < int(p); i++ {
		n.chains = append(n.chains, Create(uint(i)))
	}
	return n
}

func (mchain *MultiBlockchain) Height() []uint {
	heights := make([]uint, len(mchain.chains))
	for i := range heights {
		heights[i] = mchain.chains[i].height
	}
	return heights
}
func (mchain *MultiBlockchain) Genesis() []block.BlockId {
	gen := []block.BlockId{}
	for _, bc := range mchain.chains {
		gen = append(gen, bc.Genesis())
	}
	return gen
}
func (mchain *MultiBlockchain) End() []block.BlockId {
	ends := []block.BlockId{}
	for _, bc := range mchain.chains {
		ends = append(ends, bc.End())
	}
	return ends
}
func (mchain *MultiBlockchain) AddBlock(blk *block.Block) {
	id := blk.Id()
	mchain.blockStore.AddBlock(blk)
	// fmt.Printf("Proposed block of height %d", blk.Header.Height)
	// mchain.mchain[id] = *block
	if blk.Header.Height > mchain.height {
		mchain.end = id
		mchain.height = blk.Header.Height
		mchain.endTeller.Update([]transaction.Transaction{}, blk.Txns)
	}
}

func (mchain *MultiBlockchain) GetTeller(blId block.BlockId) *teller.Teller {
	if mchain.end == blId {
		return &mchain.endTeller
	} else {
		targetId := blId
		targetBlock := mchain.GetBlock(blId)
		if targetBlock == nil {
			return nil
		}
		curId := mchain.end
		curBlock := mchain.GetBlock(curId)
		removed := []*block.Block{}
		added := []*block.Block{}
		for targetId != curId {
			if targetBlock.Header.Height > curBlock.Header.Height {
				added = append(added, targetBlock)
				targetId = block.BlockId(targetBlock.Header.PrevBlock)
				targetBlock = mchain.GetBlock(targetId)
			} else if targetBlock.Header.Height < curBlock.Header.Height {
				removed = append(removed, curBlock)
				curId = block.BlockId(curBlock.Header.PrevBlock)
				curBlock = mchain.GetBlock(curId)
			} else {
				added = append(added, targetBlock)
				targetId = block.BlockId(targetBlock.Header.PrevBlock)
				targetBlock = mchain.GetBlock(targetId)

				removed = append(removed, curBlock)
				curId = block.BlockId(curBlock.Header.PrevBlock)
				curBlock = mchain.GetBlock(curId)
			}
		}
		removedTxn := []transaction.Transaction{}
		addedTxn := []transaction.Transaction{}

		for _, b := range removed {
			removedTxn = append(removedTxn, b.Txns...)
		}
		for _, b := range added {
			addedTxn = append(addedTxn, b.Txns...)
		}
		newTeller := mchain.endTeller.Copy()

		newTeller.Update(removedTxn, addedTxn)
		return &newTeller

	}
}

func (mchain *MultiBlockchain) GetBlock(id block.BlockId) *block.Block {

	return mchain.blockStore.GetBlock(id)
}
