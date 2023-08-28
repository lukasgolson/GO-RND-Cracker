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
	root *Node
}

func NewBkTree(rootWord []byte, seed int32) *BkTree {
	root := NewNode(math.MaxUint8, rootWord, seed)
	return &BkTree{
		root: root,
	}
}

func (bk *BkTree) Add(word []byte, data int32) {
	distance := MeyersDifferenceAlgorithm(bk.root.Word, word)
	bk.addNode(bk.root, word, data, uint(distance))
}

func (bk *BkTree) addNode(node *Node, word []byte, data int32, distance uint) {
	for {
		foundNode := findNodeWithDistance(node, distance)
		if foundNode != nil {
			node = foundNode
			continue
		}

		newNode := NewNode(distance, word, data)
		node.Children = append(node.Children, newNode)
		break
	}
}

func findNodeWithDistance(node *Node, distance uint) *Node {
	for _, child := range node.Children {
		if child.Distance == distance {
			return child
		}
	}
	return nil
}

func (bk *BkTree) Search(queryWord []byte, tolerance int) []SearchResult {
	result := bk.searchNode(bk.root, queryWord, tolerance)
	return result
}

func (bk *BkTree) searchNode(node *Node, queryWord []byte, tolerance int) []SearchResult {
	var result []SearchResult
	nodeQueue := list.New()
	nodeQueue.PushBack(node)

	for nodeQueue.Len() > 0 {
		currentNode := nodeQueue.Remove(nodeQueue.Front()).(*Node)
		distance := MeyersDifferenceAlgorithm(currentNode.Word, queryWord)

		if distance <= tolerance {
			result = append(result, SearchResult{
				Word:     currentNode.Word,
				Seed:     int(currentNode.Seed),
				Distance: distance,
			})
		}

		for _, currentNodeChild := range currentNode.Children {
			if isWithinTolerance(currentNodeChild.Distance, uint(distance), tolerance) {
				nodeQueue.PushBack(currentNodeChild)
			}
		}
	}

	return result
}

func isWithinTolerance(a, b uint, tolerance int) bool {
	return int(math.Abs(float64(a)-float64(b))) <= tolerance
}

func MeyersDifferenceAlgorithm(s1 []byte, s2 []byte) int {

	if len(s1) == 0 {
		return len(s2) // Return the length of s2 as the score
	}
	if len(s2) == 0 {
		return len(s1) // Return the length of s1 as the score
	}

	score := len(s2)

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
