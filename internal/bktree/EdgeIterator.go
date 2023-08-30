package bktree

import (
	"bytes"
	"fmt"
)

type EdgeIterator struct {
	bkTree       *BkTree
	currentIndex int
}

func NewEdgeIterator(bkTree *BkTree) *EdgeIterator {
	return &EdgeIterator{
		bkTree:       bkTree,
		currentIndex: 0,
	}
}

func (it *EdgeIterator) Next() (*Edge, error) {
	if it.currentIndex >= len(it.bkTree.Edges) {
		return nil, fmt.Errorf("iterator has reached the end")
	}

	offset := it.currentIndex * edgeByteSize
	edgeData := it.bkTree.Edges[offset : offset+edgeByteSize]

	reader := bytes.NewReader(edgeData)

	edge, err := DeserializeEdgeFromBinaryStream(reader)
	if err != nil {
		return nil, err
	}

	it.currentIndex++
	return edge, nil
}

func (it *EdgeIterator) HasNext() bool {
	return it.currentIndex < len(it.bkTree.Edges)
}
