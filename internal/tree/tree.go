package tree

import (
	"awesomeProject/internal/fileArray"
	"awesomeProject/internal/serialization"
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

func (tree *Tree) getFileNames() (string, string) {
	return tree.Nodes.GetFileName(), tree.Edges.GetFileName()
}

func (tree *Tree) isEmpty() bool {
	return tree.Nodes.Count() == 0
}

func (tree *Tree) addNode(data [32]byte, seed int32) (serialization.Offset, error) {

	id := tree.Nodes.Count()

	_, err := fileArray.Append[Node](tree.Nodes, *NewNode(id, data, seed))

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (tree *Tree) AddEdge(parentIndex, childIndex serialization.Offset, distance uint32) (serialization.Offset, error) {
	newEdge := NewEdge(parentIndex, childIndex, distance)
	id, err := fileArray.Append(tree.Edges, *newEdge)
	return id, err
}

func (tree *Tree) getNodeByIndex(index serialization.Offset) (Node, error) {
	return fileArray.GetItemFromIndex[Node](tree.Nodes, index)
}

func (tree *Tree) getEdgeByIndex(index serialization.Offset) (Edge, error) {
	return fileArray.GetItemFromIndex[Edge](tree.Edges, index)
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
