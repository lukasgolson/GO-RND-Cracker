package bktree

import (
	"awesomeProject/internal/algorithms"
	"awesomeProject/internal/tree"
)

var nodes []tree.Node
var edges []tree.Edge

func addToBKTree(rootIndex uint32, word [tree.NodeWordSize]byte, seed int32) {
	if nodes[rootIndex].Word == [tree.NodeWordSize]byte{} {
		nodes[rootIndex].Word = word
		nodes[rootIndex].Seed = seed
		return
	}

	distance := algorithms.MeyersDifferenceAlgorithm(nodes[rootIndex].Word[:], word[:])
	if edges[rootIndex].ChildIndex == 0 {
		newNode := tree.Node{
			ID:   uint32(len(nodes)),
			Word: word,
			Seed: seed,
		}
		nodes = append(nodes, newNode)
		edges[rootIndex].ChildIndex = newNode.ID
		edges = append(edges, tree.Edge{ParentIndex: rootIndex, ChildIndex: newNode.ID, Distance: distance})
	} else {
		addToBKTree(edges[rootIndex].ChildIndex, word, seed)
	}
}

func searchInBKTree(rootIndex uint32, maxDist uint16, query [tree.NodeWordSize]byte) []tree.Node {
	var result []tree.Node

	if nodes[rootIndex].Word == [tree.NodeWordSize]byte{} {
		return result
	}

	distance := algorithms.MeyersDifferenceAlgorithm(nodes[rootIndex].Word[:], query[:])
	if distance <= maxDist {
		result = append(result, nodes[rootIndex])
	}

	for _, edge := range edges {
		if edge.ParentIndex == rootIndex && distance-edge.Distance <= maxDist {
			result = append(result, searchInBKTree(edge.ChildIndex, maxDist, query)...)
		}
	}

	return result
}
