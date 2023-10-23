package tree

import (
	"encoding/binary"
	"github.com/lukasgolson/FileArray/serialization"
)

type edge struct {
	ChildIndex serialization.Offset
	Distance   uint32
}

func newEdge(childIndex serialization.Offset, distance uint32) *edge {
	return &edge{
		ChildIndex: childIndex,
		Distance:   distance,
	}
}

func (e edge) SerializeToBinaryStream(buf []byte) error {

	binary.LittleEndian.PutUint64(buf[0:8], uint64(e.ChildIndex)) // Convert int64 to little-endian binary and put it in the buffer
	binary.LittleEndian.PutUint32(buf[8:12], e.Distance)          // Convert uint32 to little-endian binary and put it in the buffer

	return nil
}

func (e edge) DeserializeFromBinaryStream(buf []byte) (edge, error) {

	e.ChildIndex = serialization.Offset(binary.LittleEndian.Uint64(buf[0:8])) // Read the little-endian binary from the buffer and convert to offset
	e.Distance = binary.LittleEndian.Uint32(buf[8:12])                        // Read the little-endian binary from the buffer and convert to uint32

	return e, nil
}

func (e edge) StrideLength() serialization.Length {
	return 8 + 4
}

func (e edge) IDByte() byte {
	return 'E'
}
