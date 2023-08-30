package bktree

type Node struct {
	ID   uint
	Word []byte
	Seed int32
}

func NewNode(ID uint, word []byte, seed int32) *Node {
	return &Node{
		ID:   ID,
		Word: word,
		Seed: seed,
	}
}
