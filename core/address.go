package core

import (
	"crypto/ecdsa"
)

// Address 0x0 will be the miner source address
type Address ecdsa.PublicKey

func (a Address) ByteMarshaling() []byte {
	return append(a.X.Bytes(), a.Y.Bytes()...)
}

func (a Address) ComputeHash() HashType {
	return ComputeHash(a.ByteMarshaling())
}
