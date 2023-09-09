package tree

import (
	"awesomeProject/internal/serialization"
	"encoding/binary"
	"io"
)

type edge struct {
	ChildIndex serialization.Offset
	Distance   uint32
}

func NewEdge(parentIndex serialization.Offset, childIndex serialization.Offset, distance uint32) *edge {
	return &edge{
		ChildIndex: childIndex,
		Distance:   distance,
	}
}

func (e edge) SerializeToBinaryStream(writer io.Writer) error {

	err := binary.Write(writer, binary.LittleEndian, e.ChildIndex)
	if err != nil {
		return err
	}

	err = binary.Write(writer, binary.LittleEndian, e.Distance)
	if err != nil {
		return err
	}

	return nil
}

func (e edge) DeserializeFromBinaryStream(reader io.Reader) (edge, error) {

	var childIndex serialization.Offset
	err := binary.Read(reader, binary.LittleEndian, &childIndex)
	if err != nil {
		return e, err
	}

	var distance uint32
	err = binary.Read(reader, binary.LittleEndian, &distance)
	if err != nil {
		return e, err
	}

	e.ChildIndex = childIndex
	e.Distance = distance

	return e, nil
}

func (e edge) StrideLength() serialization.Length {
	return serialization.Length(binary.Size(e.ChildIndex) + binary.Size(e.Distance))
}

func (e edge) IDByte() byte {
	return 'E'
}
