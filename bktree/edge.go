package bktree

type Edge struct {
	ParentIndex uint
	ChildIndex  uint
	Distance    uint
}

func NewEdge(parentIndex uint, childIndex uint, distance uint) *Edge {
	return &Edge{
		ParentIndex: parentIndex,
		ChildIndex:  childIndex,
		Distance:    distance,
	}
}
