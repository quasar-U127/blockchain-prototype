package wallet

import (
	"crypto/ecdsa"
	"crypto/rand"
)

type Wallet struct {
	keys      map[string]*ecdsa.PrivateKey
	keyToName map[Address]string
}

func CreateWallet() Wallet {
	return Wallet{keys: map[string]*ecdsa.PrivateKey{}, keyToName: map[Address]string{}}
}

func (w *Wallet) Add(addName string) {
	pvKey, _ := ecdsa.GenerateKey(curve, rand.Reader)
	w.keys[addName] = pvKey
	w.keyToName[Address{x: pvKey.X, y: pvKey.Y}] = addName
}

func (w *Wallet) Get(addName string) Address {
	pbKey := w.keys[addName].PublicKey
	return Address{x: pbKey.X, y: pbKey.Y}
}

func (w *Wallet) GetAddressNames() []string {
	names := []string{}
	for name := range w.keys {
		names = append(names, name)
	}
	return names
}
