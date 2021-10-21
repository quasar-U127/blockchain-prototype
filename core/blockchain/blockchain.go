package blockchain

import (
	"blockchain-prototype/core/structure/block"
	"blockchain-prototype/core/verifier"
)

type BlockChain struct {
	chain   map[block.BlockId]block.Block
	height  uint
	genesis block.BlockId
	end     block.BlockId
}

// type BlockChainState struct {
// 	Chain    BlockChain
// 	Frontier UTXOFrontier
// 	Pool     MemPool
// 	Mined    MemPool
// }

func (chain *BlockChain) Height() uint           { return chain.height }
func (chain *BlockChain) Genesis() block.BlockId { return chain.genesis }
func (chain *BlockChain) End() block.BlockId     { return chain.end }
func (chain *BlockChain) AddBlock(block *block.Block) {
	id := block.Id()
	chain.chain[id] = *block
	if block.Header.Height > chain.height {
		chain.end = id
		chain.height = block.Header.Height
	}
}
func (chain *BlockChain) VerifyBlock(b *block.Block, ver *verifier.Verifier) bool {
	verified := ver.VerifyBlock(b)
	if verified {
		prevId := b.Header.PrevBlock
		prevBlock, found := chain.chain[block.BlockId(prevId)]
		return found && prevBlock.Header.Height+1 == b.Header.Height
	}
	return false
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
