package tree

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

func (e *Edge) SerializeToBinaryStream(writer io.Writer) error {
	err := binary.Write(writer, binary.LittleEndian, e.ParentIndex)
	if err != nil {
		return err
	}

	err = binary.Write(writer, binary.LittleEndian, e.ChildIndex)
	if err != nil {
		return err
	}

	err = binary.Write(writer, binary.LittleEndian, e.Distance)
	if err != nil {
		return err
	}

	return nil
}

func (e *Edge) DeserializeFromBinaryStream(reader io.Reader) error {
	var parentIndex uint32
	err := binary.Read(reader, binary.LittleEndian, &parentIndex)
	if err != nil {
		return err
	}

	var childIndex uint32
	err = binary.Read(reader, binary.LittleEndian, &childIndex)
	if err != nil {
		return err
	}

	var distance uint16
	err = binary.Read(reader, binary.LittleEndian, &distance)
	if err != nil {
		return err
	}

	e.ParentIndex = parentIndex
	e.ChildIndex = childIndex
	e.Distance = distance

	return nil
}

func (e *Edge) SerializedSize() uint64 {
	return uint64(binary.Size(e.ParentIndex) + binary.Size(e.ChildIndex) + binary.Size(e.Distance))
}
