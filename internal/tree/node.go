package tree

import (
	"awesomeProject/internal/serialization"
	"encoding/binary"
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

func (n node) SerializeToBinaryStream(buf []byte) error {

	binary.LittleEndian.PutUint64(buf[0:8], uint64(n.ID))                               // Convert int64 to little-endian binary and put it in the buffer
	copy(buf[8:8+NodeWordSize], n.Word[:])                                              // Copy the word into the buffer
	binary.LittleEndian.PutUint32(buf[8+NodeWordSize:8+NodeWordSize+4], uint32(n.Seed)) // Convert int32 to little-endian binary and put it in the buffer

	return nil
}

func (n node) DeserializeFromBinaryStream(buf []byte) (node, error) {

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
