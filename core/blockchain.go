package core

type BlockChain struct {
	Chain map[HashType]Block
	N     uint
	Start HashType
	End   HashType
}

type BlockChainState struct {
	Chain    BlockChain
	Frontier UTXOFrontier
	Pool     MemPool
}

const BlockLimit int = 10

func CreateBlockChainState() *BlockChainState {
	state := BlockChainState{
		Chain: BlockChain{
			Chain: map[HashType]Block{},
			N:     0,
			End:   ZeroHash(),
			Start: ZeroHash(),
		},
		Frontier: UTXOFrontier{
			Frontier: map[OutPoint]Cookies{},
		},
		Pool: MemPool{Txns: map[TransactionId]Transaction{}},
	}
	state.Chain.End = ZeroHash()
	return &state
}

func (state *BlockChainState) MineBlock(address *Address) *Block {
	txns := state.Pool.GetTransactions(BlockLimit-1, &state.Frontier)
	fees := state.Frontier.GetFees(txns)
	coinBaseTxn := Transaction{
		TxIn: []OutPoint{CoinBaseOutpoint()},
		Output: []TXO{
			{Reciever: *address, Value: 50 + fees},
		},
	}
	txns = append(txns, coinBaseTxn)
	prevHash := state.Chain.End
	block := MineBlock(txns, prevHash, state.Chain.N+1)
	if block != nil {

		state.Chain.N++
		state.Chain.Chain[block.ComputeHash()] = *block
		state.Chain.End = block.ComputeHash()
		if state.Chain.N == 1 {
			state.Chain.Start = state.Chain.End
		}
		for _, txn := range txns {
			delete(state.Pool.Txns, txn.GetId())
		}
	}
	return block
}

func (state *BlockChainState) InsertTransaction(txn *Transaction) {
	state.Pool.InsertTransaction(txn)
}
