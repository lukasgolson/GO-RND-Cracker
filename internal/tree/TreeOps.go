package tree

import (
	"awesomeProject/internal/algorithms"
	"awesomeProject/internal/serialization"
	"awesomeProject/internal/util"
	"math"
)

type SearchResult struct {
	Word     [NodeWordSize]byte
	Seed     int32
	Distance uint32
}

func (tree *Tree) Add(word [NodeWordSize]byte, seed int32) error {
	if tree.isEmpty() {

		if _, err := tree.addNode(word, seed); err != nil {
			return err
		}

		return nil
	}

	currentNodeIndex := serialization.Offset(0)

	for {
		currentNode, err := tree.getNodeByIndex(currentNodeIndex)
		if err != nil {
			return err
		}
		editDistance := algorithms.MeyersDifferenceAlgorithm(currentNode.Word[:], word[:])

		if editDistance == 0 {
			return nil
		}

		childNodeIndex, found := tree.findChildNodeWithDistance(currentNodeIndex, editDistance)

		if !found {

			id, err := tree.addNode(word, seed)

			if err != nil {
				return err
			}

			if err := tree.addEdge(currentNodeIndex, id, editDistance); err != nil {
				return err
			}

			return nil
		}

		currentNodeIndex = childNodeIndex
	}
}



func (tree *Tree) FindClosestElement(word [NodeWordSize]byte, maxDistance uint32) SearchResult {
	if tree.nodes.Count() == 0 {
		return SearchResult{[NodeWordSize]byte{}, 0, math.MaxUint32}
	}

	// Removed preallocation of nodesArray to simplify
	nodes := []serialization.Offset{0}  // Added root node directly
	var bestWord node
	bestDistance := maxDistance

	// Combined the node popping and nodeID retrieval into one step
	for len(nodes) > 0 {
		nodeID, nodes := nodes[len(nodes)-1], nodes[:len(nodes)-1]
		
		n, err := tree.getNodeByIndex(nodeID)
		if err != nil {
			return SearchResult{[NodeWordSize]byte{}, 0, math.MaxUint32}
		}

		dU := algorithms.MeyersDifferenceAlgorithm(n.Word[:], word[:])

		// Removed unnecessary block
		if dU < bestDistance {
			bestWord, bestDistance = n, dU
			if bestDistance == 0 {
				break
			}
		}

		// Replaced egressArcs with inline function call
		for _, edge := range tree.getEgressArcs(nodeID) {
			dUV := uint32(util.Abs(int32(edge.Distance) - int32(dU)))
			if dUV <= bestDistance {
				nodes = append(nodes, edge.ChildIndex)
			}
		}
	}

	// Removed unnecessary check for bestDistance
	if bestDistance == maxDistance {
		return SearchResult{[NodeWordSize]byte{}, 0, math.MaxUint32}
	}

	return SearchResult{bestWord.Word, bestWord.Seed, bestDistance}
}
