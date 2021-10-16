package core

type Block struct {
	// We can create a merkel a tree later on but not now
	Txns []Transaction
}
