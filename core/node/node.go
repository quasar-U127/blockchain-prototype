package node

import (
	"blockchain-prototype/core/blockchain"
	"blockchain-prototype/core/mempool"
	"blockchain-prototype/core/mining"
	"blockchain-prototype/core/structure/block"
	"blockchain-prototype/core/structure/transaction"
	"blockchain-prototype/core/verifier"
	"blockchain-prototype/core/wallet"
)

type Node struct {
	chain blockchain.BlockChain
	pool  mempool.MemPool
	ver   verifier.Verifier
}

func CreateNode() Node {
	return Node{
		chain: blockchain.Create(),
		pool:  mempool.CreateMemPool(),
		ver:   verifier.CreateVerifier(),
	}
}

func (n *Node) GetBlock(id block.BlockId) *block.Block {
	return n.chain.GetBlock(id)
}
func (n *Node) GetGenesis() block.BlockId {
	return n.chain.Genesis()
}
func (n *Node) GetEnd() block.BlockId {
	return n.chain.End()
}
func (n *Node) GetHeight() uint {
	return n.chain.Height()
}

func (n *Node) GetUTXO(add *wallet.Address) map[transaction.OutPoint]transaction.Output {
	return n.chain.GetTeller(n.chain.End()).GetUTXOs(add)
}

func (n *Node) Mine(add *wallet.Address) bool {
	newBlock := mining.MineBlock(&n.chain, add, &n.ver, &n.pool)
	if newBlock != nil {
		n.AddBlock(newBlock)
		return true
	}
	return false
}

func (n *Node) AddBlock(blk *block.Block) bool {

	if n.ver.VerifyBlock(blk, &n.chain) {
		n.chain.AddBlock(blk)
		for _, t := range blk.Txns {
			n.pool.RemoveTransaction(&t)
		}
		return true
	}
	return false
}

func (n *Node) AddTransaction(txn *transaction.Transaction) bool {
	if n.ver.VerifyTransaction(txn, n.chain.GetTeller(n.chain.End())) {
		n.pool.InsertTransaction(txn)
		return true
	}
	return false
}
