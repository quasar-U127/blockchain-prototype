package core

import "crypto/sha256"

const HashSize = sha256.Size

type HashType [HashSize]byte

func ComputeHash(b []byte) HashType {
	return sha256.Sum256(b)
}
