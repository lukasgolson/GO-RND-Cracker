package tree

import (
	"awesomeProject/internal/serialization"
)

func (tree *Tree) getNextNodeID() uint32 {
	return uint32(tree.Nodes.Count())
}

func (tree *Tree) findChildNodeWithDistance(parentNodeID serialization.Offset, distance uint32) (serialization.Offset, bool) {

	egressArcs := tree.getEgressArcs(parentNodeID)

	for i, arc := range egressArcs {
		if arc.Distance == distance {
			return egressArcs[i].ChildIndex, true
		}
	}

	return 0, false
}

func (tree *Tree) getEgressArcs(parentNodeID serialization.Offset) []Edge {
	// Create a slice to store egress arcs
	egressArcs := make([]Edge, 0)

	valid, count, err := tree.Edges.Count(parentNodeID)

	if err != nil {
		return nil
	}

	if !valid {
		return egressArcs
	}

	for i := serialization.Length(0); i < count; i++ {
		edge, err := tree.Edges.Get(parentNodeID, i)
		if err != nil {
			continue
		}

		egressArcs = append(egressArcs, edge)
	}

	return egressArcs
}
