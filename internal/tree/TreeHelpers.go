package tree

import (
	"awesomeProject/internal/serialization"
)

func (tree *Tree) getNextNodeID() uint32 {
	return uint32(tree.Nodes.Count())
}

func (tree *Tree) findChildNodeWithDistance(parentIndex serialization.Offset, distance uint32) (serialization.Offset, bool) {
	for i := serialization.Offset(0); i < tree.Edges.Count(); i++ {
		edge, err := tree.getEdgeByIndex(i)
		if err != nil {
			continue
		}

		if edge.ParentIndex == parentIndex && edge.Distance == distance {
			return edge.ChildIndex, true
		}
	}

	return 0, false
}

func (tree *Tree) getEgressArcs(u serialization.Offset) []Edge {
	// Create a slice to store egress arcs
	egressArcs := make([]Edge, 0)

	// Iterate through the edges in the tree
	for i := serialization.Offset(0); i < tree.Edges.Count(); i++ {

		edge, err := tree.getEdgeByIndex(i)
		if err != nil {
			continue
		}

		// Check if the edge's parent index matches the given node index
		if edge.ParentIndex == u {
			// Append the edge to the egressArcs slice
			egressArcs = append(egressArcs, edge)
		}
	}

	// Return the egress arcs for the specified node
	return egressArcs
}
