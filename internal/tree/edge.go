package tree

type Edge struct {
	ParentIndex uint32
	ChildIndex  uint32
	Distance    uint16
}

func NewEdge(parentIndex uint32, childIndex uint32, distance uint16) *Edge {
	return &Edge{
		ParentIndex: parentIndex,
		ChildIndex:  childIndex,
		Distance:    distance,
	}
}
