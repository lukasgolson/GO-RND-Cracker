package serialization

import "io"

type Serializer interface {
	SerializeToBinaryStream(writer io.Writer) error
	DeserializeFromBinaryStream(reader io.Reader) error
	SerializedSize() uint64
}
