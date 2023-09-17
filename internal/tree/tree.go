package tree

import (
	"awesomeProject/internal/fileArray"
	"awesomeProject/internal/fileLinkedList"
	"awesomeProject/internal/serialization"
	"fmt"
)

type Tree struct {
	nodes *fileArray.FileArray[node]
	edges *fileLinkedList.FileLinkedList[edge]

	closing bool
}

func NewOrLoad(filename string) (*Tree, error) {
	bkTree := &Tree{}

	nodesFilename := fmt.Sprintf("%s.nodes.bin", filename)
	edgesFilename := fmt.Sprintf("%s.edges", filename)

	var err error

	bkTree.nodes, err = fileArray.NewFileArray[node](nodesFilename)
	bkTree.edges, err = fileLinkedList.NewFileLinkedList[edge](edgesFilename)

	if err != nil {
		return nil, err
	}

	return bkTree, nil
}

func (tree *Tree) getFileNames() (string, string, string) {

	file1, file2 := tree.edges.GetFileName()

	return file1, file2, tree.nodes.GetFileName()
}

func (tree *Tree) isEmpty() bool {
	return tree.nodes.Count() == 0
}

func (tree *Tree) addNode(data [32]byte, seed int32) (serialization.Offset, error) {

	id := tree.nodes.Count()

	_, err := tree.nodes.Append(*NewNode(id, data, seed))

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (tree *Tree) addEdge(parentIndex, childIndex serialization.Offset, distance uint32) error {
	newEdge := newEdge(childIndex, distance)

	err := tree.edges.Add(parentIndex, *newEdge)
	if err != nil {
		return err
	}

	return err
}

func (tree *Tree) getNodeByIndex(index serialization.Offset) (node, error) {
	return tree.nodes.GetItemFromIndex(index)
}

func (tree *Tree) Close() error {
	if tree.closing {
		return nil
	}
	tree.closing = true

	var err error // Declare err variable

	err = tree.nodes.Close()
	err = tree.edges.Close()

	if err != nil {
		return err
	}

	return nil
}

func (tree *Tree) ShrinkWrap() error {
	err := tree.nodes.ShrinkWrap()
	if err != nil {
		return err
	}

	err = tree.edges.ShrinkWrap()
	if err != nil {
		return err
	}

	return nil
}

func (tree *Tree) Length() serialization.Length {
	return tree.nodes.Count()
}

func (tree *Tree) PreExpand(length serialization.Length) error {

	err := tree.nodes.Expand(length)
	if err != nil {
		return err
	}
	err = tree.edges.ExpandIndex(length)
	if err != nil {
		return err
	}

	err = tree.edges.ExpandElements(length)
	if err != nil {
		return err
	}

	return nil

}
