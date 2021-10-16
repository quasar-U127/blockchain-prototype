package core

type TransactionId uint32

type TransactionInput struct {
	Id TransactionId
	N  uint
}

type TXO struct {
	Reciever Address
	Value    Cookies
}

// No Scripts and no contracts what so ever we are only implementing transfer between addresses
type Transaction struct {
	TxIn   []TransactionInput
	Output []TXO
}

func (t Transaction) GetHash() uint32 {
	return 0
}

func (t Transaction) GetId() TransactionId {

}
