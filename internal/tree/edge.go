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

func NewEdge(childIndex serialization.Offset, distance uint32) *edge {
	return &edge{
		ChildIndex: childIndex,
		Distance:   distance,
	}
}

func (e edge) SerializeToBinaryStream(writer io.Writer) error {

	buf := make([]byte, 8+4) // Create a buffer for int64 (8 bytes) and uint32 (4 bytes)

	binary.LittleEndian.PutUint64(buf[0:8], uint64(e.ChildIndex)) // Convert int64 to little-endian binary and put it in the buffer
	binary.LittleEndian.PutUint32(buf[8:12], e.Distance)          // Convert uint32 to little-endian binary and put it in the buffer

	_, err := writer.Write(buf)
	return err
}

func (e edge) DeserializeFromBinaryStream(reader io.Reader) (edge, error) {

	buf := make([]byte, 8+4) // Create a buffer for int64 (8 bytes) and uint32 (4 bytes)

	_, err := io.ReadFull(reader, buf)
	if err != nil {
		return e, err
	}

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
