package number

import (
	"awesomeProject/internal/serialization"
	"encoding/binary"
	"io"
)

type Number struct {
	Value int64
}

// NewNumber creates a new num instance with the provided value.
func NewNumber(val int64) Number {
	return Number{Value: val}
}

// SerializeToBinaryStream serializes the num struct to a binary stream.
func (number Number) SerializeToBinaryStream(writer io.Writer) error {
	buf := make([]byte, 8)                                   // Create a buffer for int64 (8 bytes)
	binary.LittleEndian.PutUint64(buf, uint64(number.Value)) // Convert int64 to little-endian binary and put it in the buffer

	_, err := writer.Write(buf)

	return err
}

// DeserializeFromBinaryStream deserializes the num struct from a binary stream.
func (number Number) DeserializeFromBinaryStream(reader io.Reader) (Number, error) {
	buf := make([]byte, 8) // Create a buffer for int64 (8 bytes)

	_, err := io.ReadFull(reader, buf)
	if err != nil {
		return number, err
	}

	number.Value = int64(binary.LittleEndian.Uint64(buf)) // Read the little-endian binary from the buffer and convert to int64
	return number, nil

}

func (number Number) StrideLength() serialization.Length {
	return serialization.Length(binary.Size(number.Value))
}

func (number Number) IDByte() byte {
	return 'N'
}
