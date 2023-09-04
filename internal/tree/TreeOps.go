package tree

import (
	"awesomeProject/internal/algorithms"
	"awesomeProject/internal/fileArray"
)

func (tree *Tree) addToBKTree(rootIndex uint32, word [NodeWordSize]byte, seed int32) (*Node, error) {

	// Step 1: Check if the current node is empty (no word). If so, add the word to the node and return.
	if tree.Nodes.Count() == 0 {
		// Create a root node if the tree is empty.
		rootNode := Node{
			ID:   0,
			Word: word,
			Seed: seed,
		}

		// Add the root node to the tree.
		err := fileArray.SetItemAtIndex(tree.Nodes, rootNode, 0)
		if err != nil {
			return nil, err
		}

		return &rootNode, nil
	}

	currentNodeIndex := rootIndex

	for {
		// Step 2: Calculate the edit distance between the current node's word and the word to be added.
		currentNode, err := fileArray.GetItemFromIndex[Node](tree.Nodes, uint64(currentNodeIndex))
		if err != nil {
			return nil, err
		}
		editDistance := algorithms.MeyersDifferenceAlgorithm(currentNode.Word[:], word[:])

		// Step 3: If the edit distance is 0, return the current node.
		if editDistance == 0 {
			return &currentNode, nil
		}

		// Step 4: Find the child node (if it exists) with the same edit distance.
		var childNodeIndex uint32
		found := false

		for edgeIndex := uint64(0); edgeIndex < tree.Edges.Count(); edgeIndex++ {
			edge, err := fileArray.GetItemFromIndex[Edge](tree.Edges, edgeIndex)
			if err != nil {
				return nil, err
			}
			if edge.ParentIndex == currentNodeIndex {
				if edge.Distance == editDistance {
					childNodeIndex = edge.ChildIndex
					found = true
					break
				}
			}
		}

		// Step 5: If the child node is not found, create a new node and add the corresponding edge.
		if !found {
			newNode := Node{
				ID:   uint32(tree.Nodes.Count()), // Assign a new ID to the node.
				Word: word,
				Seed: seed,
			}

			// Add the new node to the tree.
			err := fileArray.SetItemAtIndex(tree.Nodes, newNode, uint64(newNode.ID))
			if err != nil {
				return nil, err
			}

			// Create the edge between the current node and the new node and store it in the edge file array.
			newEdge := Edge{ParentIndex: currentNodeIndex, ChildIndex: newNode.ID, Distance: editDistance}
			err = fileArray.Append(tree.Edges, newEdge)
			if err != nil {
				return nil, err
			}

			return &newNode, nil
		}

		// Step 6: Set the current node to the found child node and continue the loop.
		currentNodeIndex = childNodeIndex
	}
}
