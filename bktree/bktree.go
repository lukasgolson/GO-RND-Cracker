package bktree

import (
	"container/list"
	"math"
)

type SearchResult struct {
	Word     []byte
	Seed     int
	Distance int
}

type BkTree struct {
	Root  *Node
	Nodes []*Node
	Edges []*Edge
}

func NewBkTree(rootWord []byte, seed int32) *BkTree {
	root := NewNode(0, rootWord, seed)
	nodes := []*Node{root}
	var edges []*Edge
	return &BkTree{
		Root:  root,
		Nodes: nodes,
		Edges: edges,
	}
}

func (bk *BkTree) Add(word []byte, data int32) {
	if len(bk.Nodes) == 0 {
		root := NewNode(0, word, data)
		bk.Nodes = append(bk.Nodes, root)
		return
	}

	currentNode := bk.Root // Start from the root node

	for currentNode != nil {
		wordDistance := MeyersDifferenceAlgorithm(currentNode.Word, word)

		if wordDistance == 0 {
			return // Node already exists, do nothing
		}

		v := bk.findChildWithDistance(currentNode, wordDistance)

		if v == nil {
			newNode := NewNode(uint32(len(bk.Nodes)), word, data)
			edge := NewEdge(currentNode.ID, newNode.ID, wordDistance)
			bk.Nodes = append(bk.Nodes, newNode)
			bk.Edges = append(bk.Edges, edge)
			return
		}

		currentNode = v // Move down the tree to the next node
	}
}

func (bk *BkTree) findChildWithDistance(node *Node, distance uint16) *Node {
	for _, edge := range bk.Edges {
		if edge.ParentIndex == node.ID && edge.Distance == distance {
			return bk.Nodes[edge.ChildIndex]
		}
	}
	return nil
}

func (bk *BkTree) Search(queryWord []byte, tolerance int) []SearchResult {
	result := bk.searchNodes(bk.Root, queryWord, tolerance)
	return result
}

func (bk *BkTree) searchNodes(node *Node, queryWord []byte, tolerance int) []SearchResult {
	var result []SearchResult
	nodeQueue := list.New()
	nodeQueue.PushBack(node)

	for nodeQueue.Len() > 0 {
		currentNode := nodeQueue.Remove(nodeQueue.Front()).(*Node)
		distance := MeyersDifferenceAlgorithm(currentNode.Word, queryWord)

		if distance <= uint16(tolerance) {
			result = append(result, SearchResult{
				Word:     currentNode.Word,
				Seed:     int(currentNode.Seed),
				Distance: int(distance),
			})
		}

		for _, edge := range bk.Edges {
			if edge.ParentIndex == currentNode.ID && isWithinTolerance(edge.Distance, distance, tolerance) {
				nodeQueue.PushBack(bk.Nodes[edge.ChildIndex])
			}
		}
	}

	return result
}

func isWithinTolerance(a, b uint16, tolerance int) bool {
	return int(math.Abs(float64(a)-float64(b))) <= tolerance
}

func MeyersDifferenceAlgorithm(s1 []byte, s2 []byte) uint16 {

	if len(s1) == 0 {
		return uint16(len(s2)) // Return the length of s2 as the score
	}
	if len(s2) == 0 {
		return uint16(len(s1)) // Return the length of s1 as the score
	}

	score := uint16(len(s2))

	peq := make([]int64, 256)
	var i int

	for i = 0; i < len(peq); i++ {
		peq[i] = 0
	}

	for i = 0; i < len(s2); i++ {
		peq[s2[i]] |= int64(1) << uint(i)
	}

	var mv int64 = 0
	var pv int64 = -1
	var last = int64(1) << uint(len(s2)-1)

	for i = 0; i < len(s1); i++ {
		eq := peq[s1[i]]

		xv := eq | mv
		xh := (((eq & pv) + pv) ^ pv) | eq

		ph := mv | ^(xh | pv)
		mh := pv & xh

		if (ph & last) != 0 {
			score++
		}
		if (mh & last) != 0 {
			score--
		}

		ph = (ph << 1) | 1
		mh = mh << 1

		pv = mh | ^(xv | ph)
		mv = ph & xv
	}

	return score
}
