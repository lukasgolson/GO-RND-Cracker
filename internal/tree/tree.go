package tree

import (
	"awesomeProject/internal/fileArray"
	"awesomeProject/internal/fileLinkedList"
	"awesomeProject/internal/serialization"
	"fmt"
)

type Tree struct {
	Nodes *fileArray.FileArray[Node]
	Edges *fileLinkedList.FileLinkedList[Edge]

	closing bool
}

func New(filename string) (*Tree, error) {
	bkTree := &Tree{}

	nodesFilename := fmt.Sprintf("%s.nodes.bin", filename)
	edgesFilename := fmt.Sprintf("%s.edges", filename)

	var err error

	bkTree.Nodes, err = fileArray.NewFileArray[Node](nodesFilename)
	bkTree.Edges, err = fileLinkedList.NewFileLinkedList[Edge](edgesFilename)

	if err != nil {
		return nil, err
	}

	return bkTree, nil
}

func (tree *Tree) getFileNames() (string, string, string) {

	file1, file2 := tree.Edges.GetFileName()

	return file1, file2, tree.Nodes.GetFileName()
}

func (tree *Tree) isEmpty() bool {
	return tree.Nodes.Count() == 0
}

func (tree *Tree) addNode(data [32]byte, seed int32) (serialization.Offset, error) {

	id := tree.Nodes.Count()

	_, err := tree.Nodes.Append(*NewNode(id, data, seed))

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (tree *Tree) addEdge(parentIndex, childIndex serialization.Offset, distance uint32) error {
	newEdge := NewEdge(parentIndex, childIndex, distance)

	err := tree.Edges.Add(parentIndex, *newEdge)
	if err != nil {
		return err
	}

	return err
}

func (tree *Tree) getNodeByIndex(index serialization.Offset) (Node, error) {
	return tree.Nodes.GetItemFromIndex(index)
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
