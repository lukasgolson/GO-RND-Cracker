package bktree

import (
	"fmt"
	"os"
	"unsafe"
)

type BkTree struct {
	Root      *Node
	Nodes     []*Node
	Edges     []*Edge
	closing   bool
	nodesFile *os.File
	edgesFile *os.File

	//Cache the size of a single node and edge
	nodeSize int
	edgeSize int
}

func NewBkTree(filename string) (*BkTree, error) {
	bkTree := &BkTree{}

	nodesFilename := fmt.Sprintf("%s.nodes", filename)
	edgesFilename := fmt.Sprintf("%s.edges", filename)

	nodesFile, err := os.OpenFile(nodesFilename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	bkTree.nodesFile = nodesFile

	edgesFile, err := os.OpenFile(edgesFilename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		nodesFile.Close()
		return nil, err
	}
	bkTree.edgesFile = edgesFile

	bkTree.nodeSize = int(unsafe.Sizeof(Node{})) // Cache the size
	bkTree.edgeSize = int(unsafe.Sizeof(Edge{})) // Cache the size

	defer func() {
		if err != nil {
			bkTree.Close()
		}
	}()

	bkTree.Nodes, err = mapNodes(nodesFile.Name(), 1024*8)
	if err != nil {
		return nil, err
	}

	bkTree.Edges, err = mapEdges(edgesFile.Name(), 1024*8)
	if err != nil {
		return nil, err
	}

	return bkTree, nil
}

func (bkTree *BkTree) Close() error {
	if bkTree.closing {
		return nil
	}
	bkTree.closing = true

	if bkTree.nodesFile != nil {
		bkTree.nodesFile.Close()
	}

	if bkTree.edgesFile != nil {
		bkTree.edgesFile.Close()
	}

	return nil
}
