package bktree

import "awesomeProject/internal/algorithms"

func (bkTree *BkTree) Add(word []byte, data int32) error {
	if len(bkTree.Nodes) == 0 {
		root := NewNode(0, word, data)
		bkTree.Nodes = append(bkTree.Nodes, root)
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
			newNodesSize := (len(bkTree.Nodes) + 1) * bkTree.nodeSize
			newEdgesSize := (len(bkTree.Edges) + 1) * bkTree.edgeSize
			if newNodesSize > len(bkTree.Nodes)*2 {
				err := bkTree.expandNodesFile(newNodesSize)
				if err != nil {
					return err
				}
			}
			if newEdgesSize > len(bkTree.Edges)*2 {
				err := bkTree.expandEdgesFile(newEdgesSize)
				if err != nil {
					return err
				}
			}

			newNode := NewNode(uint32(len(bkTree.Nodes)), word, data)
			edge := NewEdge(currentNode.ID, newNode.ID, wordDistance)
			bkTree.Nodes = append(bkTree.Nodes, newNode)
			bkTree.Edges = append(bkTree.Edges, edge)
			return nil
		}

		currentNode = childNode
	}

	return nil
}
