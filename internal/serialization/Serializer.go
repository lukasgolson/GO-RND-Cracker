package serialization

type Serializer[T any] interface {
	SerializeToBinaryStream(buffer []byte) error
	DeserializeFromBinaryStream(buffer []byte) (T, error)
	StrideLength() Length
	IDByte() byte
}
