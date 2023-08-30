package bktree

import (
	"bytes"
	"fmt"
)

func (bkTree *BkTree) GetNodeAtIndex(index int) (*Node, error) {
	offset := index * nodeByteSize // Use cached nodeSize
	if offset >= len(bkTree.Nodes) {
		return nil, fmt.Errorf("index out of bounds")
	}

	nodeData := bkTree.Nodes[offset : offset+nodeByteSize]

	nodeReader := bytes.NewReader(nodeData)

	node, err := DeserializeNodeFromBinaryStream(nodeReader)
	if err != nil {
		return nil, err
	}

	return node, err
}

func (bkTree *BkTree) AddNode(node *Node) (int, error) {
	nodeData := make([]byte, nodeByteSize)
	writer := bytes.NewBuffer(nodeData)

	err := node.SerializeToBinaryStream(writer)
	if err != nil {
		return -1, err
	}

	index := len(bkTree.Nodes) / nodeByteSize
	bkTree.Nodes = append(bkTree.Nodes, nodeData...)

	return index, nil
}

func (bkTree *BkTree) GetEdgeAtIndex(index int) (*Edge, error) {
	offset := index * edgeByteSize
	if offset >= len(bkTree.Edges) {
		return nil, fmt.Errorf("index out of bounds")
	}

	edgeData := bkTree.Edges[offset : offset+edgeByteSize]

	reader := bytes.NewReader(edgeData)

	edge, err := DeserializeEdgeFromBinaryStream(reader)
	if err != nil {
		return nil, err
	}

	return edge, err
}

func (bkTree *BkTree) AddEdge(edge *Edge) (int, error) {
	nodeData := make([]byte, nodeByteSize)
	writer := bytes.NewBuffer(nodeData)

	err := edge.SerializeToBinaryStream(writer)
	if err != nil {
		return -1, err
	}

	index := len(bkTree.Nodes) / nodeByteSize
	bkTree.Nodes = append(bkTree.Nodes, nodeData...)

	return index, nil
}
