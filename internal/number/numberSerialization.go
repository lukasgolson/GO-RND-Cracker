package number

import (
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
	err := binary.Write(writer, binary.LittleEndian, number.Value)
	if err != nil {
		return err
	}
	return nil
}

// DeserializeFromBinaryStream deserializes the num struct from a binary stream.
func (number Number) DeserializeFromBinaryStream(reader io.Reader) (Number, error) {

	err := binary.Read(reader, binary.LittleEndian, &number.Value)

	if err != nil {
		return number, err
	}
	return number, nil

}

func (number Number) StrideLength() uint64 {
	return uint64(binary.Size(number.Value))
}

func (number Number) IDByte() byte {
	return 'N'
}
