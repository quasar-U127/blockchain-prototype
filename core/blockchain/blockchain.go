package blockchain

import (
	"blockchain-prototype/core/structure/block"
	"blockchain-prototype/core/structure/transaction"
	"blockchain-prototype/core/teller"
)

type BlockChain struct {
	blockStore block.Store
	height     uint
	genesis    block.BlockId
	end        block.BlockId
	endTeller  teller.Teller
}

// type BlockChainState struct {
// 	Chain    BlockChain
// 	Frontier UTXOFrontier
// 	Pool     MemPool
// 	Mined    MemPool
// }

func Create() BlockChain {
	chain := BlockChain{
		blockStore: block.CreateStore(),
		height:     1,
		endTeller:  teller.CreateTeller(),
	}
	genesis := block.CreateGenesis()
	chain.blockStore.AddBlock(&genesis)
	chain.genesis = genesis.Id()
	chain.end = genesis.Id()
	return chain
}

func (chain *BlockChain) Height() uint           { return chain.height }
func (chain *BlockChain) Genesis() block.BlockId { return chain.genesis }
func (chain *BlockChain) End() block.BlockId     { return chain.end }
func (chain *BlockChain) AddBlock(blk *block.Block) {
	id := blk.Id()
	chain.blockStore.AddBlock(blk)
	// fmt.Printf("Proposed block of height %d", blk.Header.Height)
	// chain.chain[id] = *block
	if blk.Header.Height > chain.height {
		chain.end = id
		chain.height = blk.Header.Height
		chain.endTeller.Update([]transaction.Transaction{}, blk.Txns)
	}
}

func (chain *BlockChain) GetTeller(blId block.BlockId) *teller.Teller {
	if chain.end == blId {
		return &chain.endTeller
	} else {
		targetId := blId
		targetBlock := chain.GetBlock(blId)
		if targetBlock == nil {
			return nil
		}
		curId := chain.end
		curBlock := chain.GetBlock(curId)
		removed := []*block.Block{}
		added := []*block.Block{}
		for targetId != curId {
			if targetBlock.Header.Height > curBlock.Header.Height {
				added = append(added, targetBlock)
				targetId = block.BlockId(targetBlock.Header.PrevBlock)
				targetBlock = chain.GetBlock(targetId)
			} else if targetBlock.Header.Height < curBlock.Header.Height {
				removed = append(removed, curBlock)
				curId = block.BlockId(curBlock.Header.PrevBlock)
				curBlock = chain.GetBlock(curId)
			} else {
				added = append(added, targetBlock)
				targetId = block.BlockId(targetBlock.Header.PrevBlock)
				targetBlock = chain.GetBlock(targetId)

				removed = append(removed, curBlock)
				curId = block.BlockId(curBlock.Header.PrevBlock)
				curBlock = chain.GetBlock(curId)
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
		newTeller := chain.endTeller.Copy()

		newTeller.Update(removedTxn, addedTxn)
		return &newTeller

	}
}

func (chain *BlockChain) GetBlock(id block.BlockId) *block.Block {

	return chain.blockStore.GetBlock(id)
}

// func CreateBlockChainState() *BlockChainState {
// 	state := BlockChainState{
// 		Chain: BlockChain{
// 			Chain: map[utils.HashType]Block{},
// 			N:     0,
// 			End:   utils.ZeroHash(),
// 			Start: utils.ZeroHash(),
// 		},
// 		Frontier: CreateUTXOFrontier(),
// 		Pool:     MemPool{Txns: map[TransactionId]Transaction{}},
// 	}
// 	return &state
// }

// func (state *BlockChainState) MineBlock(address *Address) *Block {
// 	txns := []Transaction{}
// 	minedTxns := state.GetMiningTransactions(BlockLimit - 1)
// 	fees := state.Frontier.GetFees(minedTxns)
// 	coinBaseTxn := Transaction{
// 		TxIn: []OutPoint{{Id: TransactionId(utils.ZeroHash()), N: state.Chain.N + 1}},
// 		Output: []TXO{
// 			{Reciever: *address, Value: 50 + fees},
// 		},
// 	}
// 	txns = append(txns, coinBaseTxn)
// 	txns = append(txns, minedTxns...)
// 	prevHash := state.Chain.End
// 	block := MineBlock(txns, prevHash, state.Chain.N+1)
// 	if block != nil {

// 		state.Chain.N++
// 		state.Chain.Chain[block.Hash()] = *block
// 		state.Chain.End = block.Hash()
// 		if state.Chain.N == 1 {
// 			state.Chain.Start = state.Chain.End
// 		}
// 		for _, txn := range txns {
// 			delete(state.Pool.Txns, txn.GetId())
// 		}
// 		state.Frontier.Update(block.Txns)
// 	}
// 	return block
// }
