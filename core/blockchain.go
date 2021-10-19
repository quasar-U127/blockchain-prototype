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
	Mined    MemPool
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
		Frontier: CreateUTXOFrontier(),
		Pool:     MemPool{Txns: map[TransactionId]Transaction{}},
	}
	return &state
}

func (state *BlockChainState) MineBlock(address *Address) *Block {
	txns := []Transaction{}
	minedTxns := state.GetMiningTransactions(BlockLimit - 1)
	fees := state.Frontier.GetFees(minedTxns)
	coinBaseTxn := Transaction{
		TxIn: []OutPoint{{Id: TransactionId(ZeroHash()), N: state.Chain.N + 1}},
		Output: []TXO{
			{Reciever: *address, Value: 50 + fees},
		},
	}
	txns = append(txns, coinBaseTxn)
	txns = append(txns, minedTxns...)
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
		state.Frontier.Update(block.Txns)
	}
	return block
}

func (state *BlockChainState) GetMiningTransactions(n int) []Transaction {
	spent := map[OutPoint]bool{}
	txlist := []Transaction{}
	for _, txn := range state.Pool.Txns {
		if len(txlist) >= n {
			break
		}
		if state.Frontier.VerifyTransaction(&txn) {
			valid := true
			for _, outpoint := range txn.TxIn {
				if _, ok := spent[outpoint]; ok {
					valid = false
				}
			}
			if valid {

				txlist = append(txlist, txn)
				for _, outpoint := range txn.TxIn {
					spent[outpoint] = true
				}
			}
		}
	}
	return txlist
}

func (state *BlockChainState) InsertTransaction(txn *Transaction) {
	if state.Frontier.VerifyTransaction(txn) {
		state.Pool.InsertTransaction(txn)
	}
}
