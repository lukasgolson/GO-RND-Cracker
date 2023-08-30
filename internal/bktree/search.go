package bktree

import (
	"awesomeProject/internal/algorithms"
	"container/list"
	"fmt"
	"math"
)

type SearchResult struct {
	Word     []byte
	Seed     int
	Distance int
}

func (bkTree *BkTree) Search(queryWord []byte, tolerance int) []SearchResult {
	result := bkTree.searchNodes(bkTree.Root, queryWord, tolerance)
	return result
}

func (bkTree *BkTree) searchNodes(node *Node, queryWord []byte, tolerance int) []SearchResult {
	var result []SearchResult
	nodeQueue := list.New()
	nodeQueue.PushBack(node)

	for nodeQueue.Len() > 0 {
		currentNode := nodeQueue.Remove(nodeQueue.Front()).(*Node)
		distance := algorithms.MeyersDifferenceAlgorithm(currentNode.Word, queryWord)

		if distance <= uint16(tolerance) {
			result = append(result, SearchResult{
				Word:     currentNode.Word,
				Seed:     int(currentNode.Seed),
				Distance: int(distance),
			})
		}

		edgeIterator := NewEdgeIterator(bkTree)

		for edgeIterator.HasNext() {
			edge, err := edgeIterator.Next()
			if err != nil {
				fmt.Println("Error:", err)
				break
			}

			if edge.ParentIndex == currentNode.ID && isWithinTolerance(edge.Distance, distance, tolerance) {
				nodeQueue.PushBack(bkTree.Nodes[edge.ChildIndex])
			}

		}
	}

	return result
}

func isWithinTolerance(a, b uint16, tolerance int) bool {
	return int(math.Abs(float64(a)-float64(b))) <= tolerance
}

func (bkTree *BkTree) findChildWithDistance(node *Node, distance uint16) *Node {

	edgeIterator := NewEdgeIterator(bkTree)

	for edgeIterator.HasNext() {
		edge, err := edgeIterator.Next()
		if err != nil {
			fmt.Println("Error:", err)
			break
		}

		if edge.ParentIndex == node.ID && edge.Distance == distance {
			node, err = bkTree.GetNodeAtIndex(int(edge.ChildIndex))
			return node
		}

	}

	return nil
}
