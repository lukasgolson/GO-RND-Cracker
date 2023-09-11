package tree

import (
	"awesomeProject/internal/serialization"
	"encoding/binary"
	"io"
)

const NodeWordSize = 32

type node struct {
	ID   serialization.Offset
	Word [NodeWordSize]byte
	Seed int32
}

func NewNode(ID serialization.Offset, word [NodeWordSize]byte, seed int32) *node {
	return &node{
		ID:   ID,
		Word: word,
		Seed: seed,
	}
}

func (n node) SerializeToBinaryStream(writer io.Writer) error {

	buf := make([]byte, 8+NodeWordSize+4) // Create a buffer for int64 (8 bytes), word (32 bytes) and int32 (4 bytes)

	binary.LittleEndian.PutUint64(buf[0:8], uint64(n.ID))                               // Convert int64 to little-endian binary and put it in the buffer
	copy(buf[8:8+NodeWordSize], n.Word[:])                                              // Copy the word into the buffer
	binary.LittleEndian.PutUint32(buf[8+NodeWordSize:8+NodeWordSize+4], uint32(n.Seed)) // Convert int32 to little-endian binary and put it in the buffer

	_, err := writer.Write(buf)
	return err
}

func (n node) DeserializeFromBinaryStream(reader io.Reader) (node, error) {
	buf := make([]byte, 8+NodeWordSize+4) // Create a buffer for int64 (8 bytes), word (32 bytes) and int32 (4 bytes)

	_, err := io.ReadFull(reader, buf)
	if err != nil {
		return n, err
	}

	n.ID = serialization.Offset(binary.LittleEndian.Uint64(buf[0:8])) // Read the little-endian binary from the buffer and convert to offset

	copy(n.Word[:], buf[8:8+NodeWordSize])

	n.Seed = int32(binary.LittleEndian.Uint32(buf[8+NodeWordSize : 8+NodeWordSize+4])) // Read the little-endian binary from the buffer and convert to int32

	return n, nil
}

func (n node) StrideLength() serialization.Length {
	return 8 + NodeWordSize + 4
}

func (n node) IDByte() byte {
	return 'N'
}
