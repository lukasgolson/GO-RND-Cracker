package bktree

import (
	"encoding/binary"
	"io"
)

type Edge struct {
	ParentIndex uint32
	ChildIndex  uint32
	Distance    uint16
}

func NewEdge(parentIndex uint32, childIndex uint32, distance uint16) *Edge {
	return &Edge{
		ParentIndex: parentIndex,
		ChildIndex:  childIndex,
		Distance:    distance,
	}
}

func (edge *Edge) SerializeToBinaryStream(writer io.Writer) error {
	parentBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(parentBytes, edge.ParentIndex)
	_, err := writer.Write(parentBytes)
	if err != nil {
		return err
	}

	childBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(childBytes, edge.ChildIndex)
	_, err = writer.Write(childBytes)
	if err != nil {
		return err
	}

	distanceBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(distanceBytes, edge.Distance)
	_, err = writer.Write(distanceBytes)
	if err != nil {
		return err
	}

	return nil
}
