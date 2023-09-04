package serialization

import "io"

type Serializer[T any] interface {
	SerializeToBinaryStream(writer io.Writer) error
	DeserializeFromBinaryStream(reader io.Reader) (T, error)
	StrideLength() uint64
	IDByte() byte
}
