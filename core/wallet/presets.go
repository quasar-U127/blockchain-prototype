package wallet

import (
	"crypto/elliptic"
	"math/big"
)

// Address 0x0 will be the miner source address

type Address struct{ x, y *big.Int }

var curve elliptic.Curve = elliptic.P256()
