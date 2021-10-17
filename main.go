package main

import (
	"blockchain-prototype/core"
)

func main() {
	var prevHash [core.HashSize]byte
	var txns []core.Transaction
	core.MineBlock(txns, core.HashType(prevHash))

}
