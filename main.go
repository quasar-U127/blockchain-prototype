package main

import (
	"blockchain-prototype/core"
	"blockchain-prototype/server"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
)

func main() {
	pvKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	ser := server.CreateServer("127.0.0.1:1234")

	for i := 0; i < 10; i++ {
		block := server.RPCBlock{MinedBlock: nil}
		ser.MineBlock((*core.Address)(&pvKey.PublicKey), &block)
		if block.MinedBlock == nil {
			continue
		}
		for _, txn := range block.MinedBlock.Txns {
			txn.Print()
		}
		res := false
		ser.SubmitTransaction(&core.Transaction{
			TxIn:   []core.OutPoint{{Id: core.TransactionId(block.MinedBlock.Txns[0].ComputeHash()), N: 0}},
			Output: []core.TXO{{Reciever: core.Address(pvKey.PublicKey), Value: 45}},
		}, &res)
	}
	// server.Serve(&ser)
}
