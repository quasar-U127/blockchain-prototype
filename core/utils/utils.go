package utils

import "crypto/sha256"

const HashSize = sha256.Size

type HashType [HashSize]byte

func Hash(b []byte) HashType {
	return sha256.Sum256(b)
}

func ZeroHash() HashType {
	hash := new(HashType)
	for i := 0; i < HashSize; i++ {
		hash[i] = 0
	}
	return *hash
}

type Null struct{}

// type Hashable interface {
// 	Hash() HashType
// }

func MerkleHash(hashs []HashType) HashType {
	var byteValues []byte
	for _, hash := range hashs {
		byteValues = append(byteValues, hash[:]...)
	}
	return Hash(byteValues)
}
