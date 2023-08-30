package bktree

import "awesomeProject/internal/algorithms"

func (bkTree *BkTree) Add(word []byte, data int32) error {
	if len(bkTree.Nodes) == 0 {
		root := NewNode(0, word, data)

		bkTree.AddNode(root)

		return nil
	}

	currentNode := bkTree.Root

	for currentNode != nil {
		wordDistance := algorithms.MeyersDifferenceAlgorithm(currentNode.Word, word)

		if wordDistance == 0 {
			return nil
		}

		childNode := bkTree.findChildWithDistance(currentNode, wordDistance)

		if childNode == nil {
			newNode := NewNode(uint32(len(bkTree.Nodes)), word, data)
			edge := NewEdge(currentNode.ID, newNode.ID, wordDistance)

			bkTree.AddNode(newNode)

			bkTree.AddEdge(edge)

			return nil
		}

		currentNode = childNode
	}

	return nil
}
