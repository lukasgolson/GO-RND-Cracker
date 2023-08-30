package bktree

import (
	"golang.org/x/exp/mmap"
	"math"
	"os"
	"reflect"
	"unsafe"
)

func mapEdges(filename string, size int) ([]*Edge, error) {
	mappedData, err := mmap.Open(filename)
	if err != nil {
		return nil, err
	}
	defer mappedData.Close()

	numEdges := size / int(unsafe.Sizeof(Edge{}))

	edgesSlice := (*[math.MaxInt32 / unsafe.Sizeof(Edge{})]*Edge)(unsafe.Pointer(&mappedData))
	edges := edgesSlice[:numEdges:numEdges]

	return edges, nil
}

func mapNodes(filename string, size int) ([]*Node, error) {
	mappedData, err := mmap.Open(filename)
	if err != nil {
		return nil, err
	}
	defer mappedData.Close()

	numNodes := size / int(unsafe.Sizeof(Node{}))

	nodesSlice := (*[math.MaxInt32 / unsafe.Sizeof(Node{})]*Node)(unsafe.Pointer(&mappedData))
	nodes := nodesSlice[:numNodes:numNodes]

	return nodes, nil
}

func (bkTree *BkTree) expandFile(file *os.File, items interface{}, newSize int) error {
	itemSize := bkTree.nodeSize
	if reflect.TypeOf(items) == reflect.TypeOf([]*Edge{}) {
		itemSize = bkTree.edgeSize
	}

	currentSize := reflect.ValueOf(items).Len() * itemSize
	if newSize <= currentSize {
		return nil
	}

	err := file.Truncate(int64(newSize))
	if err != nil {
		return err
	}

	switch items.(type) {
	case []*Node:
		bkTree.Nodes, err = mapNodes(file.Name(), newSize)
	case []*Edge:
		bkTree.Edges, err = mapEdges(file.Name(), newSize)
	}

	if err != nil {
		return err
	}

	return nil
}

func (bkTree *BkTree) expandNodesFile(newSize int) error {
	return bkTree.expandFile(bkTree.nodesFile, bkTree.Nodes, newSize)
}

func (bkTree *BkTree) expandEdgesFile(newSize int) error {
	return bkTree.expandFile(bkTree.edgesFile, bkTree.Edges, newSize)
}
