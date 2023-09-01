package serialization

import (
	"encoding/binary"
	"io"
)

type Num struct {
	value int64
}

// Number creates a new num instance with the provided value.
func Number(val int64) Num {
	return Num{value: val}
}

// SerializeToBinaryStream serializes the num struct to a binary stream.
func (number *Num) SerializeToBinaryStream(writer io.Writer) error {
	err := binary.Write(writer, binary.LittleEndian, number.value) // <-- Use "number.value" instead of "number"
	if err != nil {
		return err
	}
	return nil
}

// DeserializeFromBinaryStream deserializes the num struct from a binary stream.
func (number *Num) DeserializeFromBinaryStream(reader io.Reader) error {
	err := binary.Read(reader, binary.LittleEndian, &number.value) // <-- Use "&number.value" instead of "number"
	if err != nil {
		return err
	}
	return nil
}

// SerializedSize returns the size of the serialized num struct.
func (number *Num) SerializedSize() uint64 {
	return uint64(binary.Size(number.value)) // <-- Use "number.value" instead of "number"
}
