package bktree

import (
	"awesomeProject/internal/interfaces"
	"os"
)

func StreamDataToBinaryFile(filename string, dataSlices []interfaces.Serializer) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	for _, data := range dataSlices {
		err := data.SerializeToBinaryStream(file)
		if err != nil {
			return err
		}
	}

	return nil
}

func (bkTree *BkTree) SerializeTree(nodeFilename string, edgeFilename string) error {
	var nodeDataSlices []interfaces.Serializer
	for _, node := range bkTree.Nodes {
		nodeDataSlices = append(nodeDataSlices, node)
	}

	err := StreamDataToBinaryFile(nodeFilename, nodeDataSlices)
	if err != nil {
		return err
	}

	var edgeDataSlices []interfaces.Serializer
	for _, edge := range bkTree.Edges {
		edgeDataSlices = append(edgeDataSlices, edge)
	}

	err = StreamDataToBinaryFile(edgeFilename, edgeDataSlices)
	if err != nil {
		return err
	}

	return nil
}
