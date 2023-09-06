package tree

import (
	"awesomeProject/internal/fileArray"
	"encoding/binary"
	"io"
)

type Edge struct {
	ParentIndex fileArray.Offset
	ChildIndex  fileArray.Offset
	Distance    uint32
}

func NewEdge(parentIndex fileArray.Offset, childIndex fileArray.Offset, distance uint32) *Edge {
	return &Edge{
		ParentIndex: parentIndex,
		ChildIndex:  childIndex,
		Distance:    distance,
	}
}

func (e Edge) SerializeToBinaryStream(writer io.Writer) error {
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

func (e Edge) DeserializeFromBinaryStream(reader io.Reader) (Edge, error) {
	var parentIndex fileArray.Offset
	err := binary.Read(reader, binary.LittleEndian, &parentIndex)
	if err != nil {
		return e, err
	}

	var childIndex fileArray.Offset
	err = binary.Read(reader, binary.LittleEndian, &childIndex)
	if err != nil {
		return e, err
	}

	var distance uint32
	err = binary.Read(reader, binary.LittleEndian, &distance)
	if err != nil {
		return e, err
	}

	e.ParentIndex = parentIndex
	e.ChildIndex = childIndex
	e.Distance = distance

	return e, nil
}

func (e Edge) StrideLength() uint64 {
	return uint64(binary.Size(e.ParentIndex) + binary.Size(e.ChildIndex) + binary.Size(e.Distance))
}

func (e Edge) IDByte() byte {
	return 'E'
}
