package bktree

import (
	"os"
)

func StreamDataToBinaryFile(filename string, dataSlices []StreamSerializer) error {
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

func (bk *BkTree) SerializeTree(nodeFilename string, edgeFilename string) error {
	var nodeDataSlices []StreamSerializer
	for _, node := range bk.Nodes {
		nodeDataSlices = append(nodeDataSlices, node)
	}

	err := StreamDataToBinaryFile(nodeFilename, nodeDataSlices)
	if err != nil {
		return err
	}

	var edgeDataSlices []StreamSerializer
	for _, edge := range bk.Edges {
		edgeDataSlices = append(edgeDataSlices, edge)
	}

	err = StreamDataToBinaryFile(edgeFilename, edgeDataSlices)
	if err != nil {
		return err
	}

	return nil
}
