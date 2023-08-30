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
	Nodes []*Node
	Edges []*Edge
}

func NewBkTree(rootWord []byte, seed int32) *BkTree {
	root := NewNode(0, rootWord, seed)
	nodes := []*Node{root}
	return &BkTree{
		Nodes: nodes,
		Edges: nil,
	}
}

func (bk *BkTree) Add(word []byte, data int32) {
	distance := MeyersDifferenceAlgorithm(bk.Nodes[0].Word, word)
	newNode := NewNode(uint(len(bk.Nodes)), word, data)
	bk.Nodes = append(bk.Nodes, newNode)
	bk.addEdge(0, newNode.ID, distance)
}

func (bk *BkTree) addEdge(parentID uint, childID uint, distance uint) {
	edge := NewEdge(parentID, childID, distance)
	bk.Edges = append(bk.Edges, edge)
}

func (bk *BkTree) Search(queryWord []byte, tolerance int) []SearchResult {
	result := bk.searchNode(bk.Nodes[0], queryWord, tolerance)
	return result
}

func (bk *BkTree) searchNode(node *Node, queryWord []byte, tolerance int) []SearchResult {
	var result []SearchResult
	nodeQueue := list.New()
	nodeQueue.PushBack(node)

	for nodeQueue.Len() > 0 {
		currentNode := nodeQueue.Remove(nodeQueue.Front()).(*Node)
		distance := MeyersDifferenceAlgorithm(currentNode.Word, queryWord)

		if distance <= uint(tolerance) {
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

func isWithinTolerance(a, b uint, tolerance int) bool {
	return int(math.Abs(float64(a)-float64(b))) <= tolerance
}

func MeyersDifferenceAlgorithm(s1 []byte, s2 []byte) uint {

	if len(s1) == 0 {
		return uint(len(s2)) // Return the length of s2 as the score
	}
	if len(s2) == 0 {
		return uint(len(s1)) // Return the length of s1 as the score
	}

	score := uint(len(s2))

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
	var last int64 = int64(1) << uint(len(s2)-1)

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
