package block

type Store struct {
	storage map[BlockId]Block
}

func CreateStore() Store {
	return Store{
		storage: map[BlockId]Block{},
	}
}
func (store *Store) GetBlock(id BlockId) *Block {
	blk := store.storage[id]
	return &blk
}

func (store *Store) AddBlock(blk *Block) {
	store.storage[blk.Id()] = *blk
}
