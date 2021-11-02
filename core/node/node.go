package node

import (
	"blockchain-prototype/core/blockchain"
	lineagetable "blockchain-prototype/core/lineage-table"
	"blockchain-prototype/core/mempool"
	"blockchain-prototype/core/mining"
	"blockchain-prototype/core/structure/block"
	"blockchain-prototype/core/structure/transaction"
	"blockchain-prototype/core/teller"
	"blockchain-prototype/core/verifier"
	"blockchain-prototype/core/wallet"
)

type Node struct {
	chains []blockchain.BlockChain
	pools  []mempool.MemPool
	ver    verifier.Verifier
	table  lineagetable.LineageTable
	tel    *teller.Teller
	p      int
}

func CreateNode(p int) Node {
	n := Node{
		chains: []blockchain.BlockChain{},
		pools:  []mempool.MemPool{},
		ver:    verifier.CreateVerifier(),
		table:  lineagetable.CreateLineageTable(),
		p:      p,
	}
	for i := 0; i < int(p); i++ {
		n.chains = append(n.chains, blockchain.Create(uint(i)))
		n.pools = append(n.pools, mempool.CreateMemPool())
	}
	return n
}

func (n *Node) GetBlock(index int, id block.BlockId) *block.Block {
	return n.chains[index].GetBlock(id)
}
func (n *Node) GetGenesis(index int) block.BlockId {
	return n.chains[index].Genesis()
}
func (n *Node) GetEnd(index int) block.BlockId {
	return n.chains[index].End()
}
func (n *Node) GetHeight(index int) uint {
	return n.chains[index].Height()
}

func (n *Node) GetUTXO(add *wallet.Address) map[transaction.OutPoint]transaction.Output {
	return n.GetDefaultTeller().GetUTXOs(add)
}

func (n *Node) Mine(index int, add *wallet.Address) bool {
	newBlock := mining.MineBlock(&n.chains[index], add, &n.ver, &n.pools[index])
	if newBlock != nil {
		n.AddBlock(index, newBlock)
		return true
	}
	return false
}

func (n *Node) AddBlock(index int, blk *block.Block) bool {

	if n.ver.VerifyBlock(blk, &n.chains[index]) {
		n.chains[index].AddBlock(blk)
		for _, t := range blk.Txns {
			n.pools[index].RemoveTransaction(&t)
			n.table.StoreTransaction(&t)
		}
		return true
	}
	return false
}

func (n *Node) AddTransaction(index int, txn *transaction.Transaction) bool {
	if n.ver.VerifyTransaction(txn, n.chains[index].GetTeller(n.chains[index].End())) {
		n.pools[index].InsertTransaction(txn)
		return true
	}
	return false
}

func (n *Node) GetFrontier() []block.BlockId {
	frontier := make([]block.BlockId, n.p)
	for i := range frontier {
		frontier[i] = n.chains[i].End()
	}
	return frontier
}

func (n *Node) GetDefaultTeller() *teller.Teller {
	return n.GetTeller(n.GetFrontier())
}
func (n *Node) GetTeller(frontier []block.BlockId) *teller.Teller {

}
