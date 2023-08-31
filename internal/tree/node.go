package tree

type Node struct {
	ID   uint32
	Word []byte
	Seed int32
}

func NewNode(ID uint32, word []byte, seed int32) *Node {
	return &Node{
		ID:   ID,
		Word: word,
		Seed: seed,
	}
}
