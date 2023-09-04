package tree

import (
	"awesomeProject/internal/fileArray"
	"fmt"
)

type Tree struct {
	Nodes *fileArray.FileArray
	Edges *fileArray.FileArray

	closing bool
}

func NewTree(filename string) (*Tree, error) {
	bkTree := &Tree{}

	nodesFilename := fmt.Sprintf("%s.nodes.bin", filename)
	edgesFilename := fmt.Sprintf("%s.edges.bin", filename)

	var err error

	bkTree.Nodes, err = fileArray.NewFileArray(Node{}, nodesFilename)
	bkTree.Edges, err = fileArray.NewFileArray(Edge{}, edgesFilename)

	if err != nil {
		return nil, err
	}

	return bkTree, nil
}

func (tree *Tree) Close() error {
	if tree.closing {
		return nil
	}
	tree.closing = true

	var err error // Declare err variable

	err = tree.Nodes.Close()
	err = tree.Edges.Close()

	if err != nil {
		return err
	}

	return nil
}
