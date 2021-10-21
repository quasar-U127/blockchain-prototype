package interfaces

import "blockchain-prototype/core/utils"

type Hashable interface {
	Hash() utils.HashType
}

type Marshalable interface {
	Marshal() []byte
	Unmarshal([]byte)
}
