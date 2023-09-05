package tree

import (
	"awesomeProject/internal/algorithms"
	"awesomeProject/internal/util"
	"math"
)

func (tree *Tree) AddToBKTree(rootIndex uint32, word [NodeWordSize]byte, seed int32) error {
	// Step 1: Check if the tree is empty, if so, create a root node.
	if tree.isEmpty() {

		if _, err := tree.addNode(word, seed); err != nil {
			return err
		}

		return nil
	}

	currentNodeIndex := rootIndex

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

			if _, err := tree.AddEdge(currentNodeIndex, id, editDistance); err != nil {
				return err
			}

			return nil
		}

		currentNodeIndex = childNodeIndex
	}
}

func (tree *Tree) FindClosestElement(rootIndex uint32, word [NodeWordSize]byte, maxDistance uint32) (*Node, uint32) {
	if tree.Nodes.Count() == 0 {
		return nil, math.MaxUint32
	}

	S := make([]uint32, 0)      // Set of nodes to process
	S = append(S, rootIndex)    // Insert the root node into S
	bestWord := Node{}          // Best matching element
	bestDistance := maxDistance // Best matching distance, initialized to maxDistance

	for len(S) != 0 {
		u := S[len(S)-1] // Pop the last node from S
		S = S[:len(S)-1]

		n, err := tree.getNodeByIndex(u)

		if err != nil {
			return nil, math.MaxUint32
		}

		dU := algorithms.MeyersDifferenceAlgorithm(n.Word[:], word[:])

		if dU < bestDistance {

			bestWord, err = tree.getNodeByIndex(u)
			bestDistance = dU
		}

		for _, edge := range tree.getEgressArcs(u) {
			v := edge.ChildIndex
			dUV := uint32(util.Abs(int32(edge.Distance) - int32(dU)))

			if dUV < bestDistance {
				S = append(S, v) // Insert v into S
			}
		}
	}

	if bestDistance == maxDistance {
		return nil, math.MaxUint32
	}

	return &bestWord, bestDistance
}
