package main

import (
	"blockchain-prototype/core"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
)

func main() {
	pvKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	state := core.CreateBlockChainState()
	for i := 0; i < 10; i++ {
		block := state.MineBlock((*core.Address)(&pvKey.PublicKey))
		fmt.Printf("%v\n", block.ComputeHash())
		fmt.Println(block.Header)
	}

}
