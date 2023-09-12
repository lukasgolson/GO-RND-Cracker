package number

import (
	"awesomeProject/internal/serialization"
	"encoding/binary"
)

type Number struct {
	Value int64
}

// NewNumber creates a new num instance with the provided value.
func NewNumber(val int64) Number {
	return Number{Value: val}
}

// SerializeToBinaryStream serializes the num struct to a binary stream.
func (number Number) SerializeToBinaryStream(buf []byte) error {
	binary.LittleEndian.PutUint64(buf, uint64(number.Value)) // Convert int64 to little-endian binary and put it in the buffer

	return nil
}

// DeserializeFromBinaryStream deserializes the num struct from a binary stream.
func (number Number) DeserializeFromBinaryStream(buf []byte) (Number, error) {

	number.Value = int64(binary.LittleEndian.Uint64(buf)) // Read the little-endian binary from the buffer and convert to int64
	return number, nil

}

func (number Number) StrideLength() serialization.Length {
	return 8
}

func (number Number) IDByte() byte {
	return 'N'
}
