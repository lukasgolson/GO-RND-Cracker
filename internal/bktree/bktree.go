package bktree

import (
	"fmt"
	"github.com/edsrzf/mmap-go"
	"os"
)

type BkTree struct {
	Root  *Node
	Nodes mmap.MMap
	Edges mmap.MMap

	NodeFile *os.File
	EdgeFile *os.File
	closing  bool
}

func NewBkTree(filename string) (*BkTree, error) {
	bkTree := &BkTree{}

	nodesFilename := fmt.Sprintf("%s.nodes", filename)
	edgesFilename := fmt.Sprintf("%s.edges", filename)

	var err error

	bkTree.NodeFile, err = os.OpenFile(nodesFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.FileMode(0644))

	bkTree.EdgeFile, err = os.OpenFile(edgesFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.FileMode(0644))

	bkTree.NodeFile.Truncate(1024 * 1024)
	bkTree.EdgeFile.Truncate(1024 * 1024)

	bkTree.Nodes, err = mmap.Map(bkTree.NodeFile, mmap.RDWR, 0)

	bkTree.Edges, err = mmap.Map(bkTree.EdgeFile, mmap.RDWR, 0)

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

	var err error // Declare err variable

	if bkTree.Nodes != nil {
		err = bkTree.Nodes.Unmap()
		err = bkTree.NodeFile.Close()
	}

	if bkTree.Edges != nil {
		err = bkTree.Edges.Unmap()
		err = bkTree.EdgeFile.Close()
	}

	if err != nil {
		return err
	}

	return nil
}
