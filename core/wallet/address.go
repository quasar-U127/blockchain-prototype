package wallet

import (
	"blockchain-prototype/core/utils"
)

func (a *Address) Marshal() []byte {
	return append(a.x.Bytes(), a.y.Bytes()...)
}

func (a *Address) ComputeHash() utils.HashType {
	return utils.Hash(a.Marshal())
}

func (a *Address) ShortString(l int) string {
	return a.String()[:l]
}

func (a *Address) String() string {
	return a.x.Text(35)
}

func (a *Address) Verify() bool {
	return true
}
