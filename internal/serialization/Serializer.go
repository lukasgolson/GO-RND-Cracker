package serialization

import "io"

type Serializer[T any] interface {
	SerializeToBinaryStream(writer io.Writer) error
	DeserializeFromBinaryStream(reader io.Reader) (T[T], error)
	SerializedSize() uint64
}
