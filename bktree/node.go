package bktree

type Node struct {
	Distance uint
	Word     []byte
	Seed     int32
	Children []*Node
}

func NewNode(distance uint, word []byte, seed int32) *Node {
	return &Node{
		Distance: distance,
		Word:     word,
		Seed:     seed,
		Children: []*Node{},
	}
}
