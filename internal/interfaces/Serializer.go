package interfaces

import "io"

type Serializer interface {
	SerializeToBinaryStream(writer io.Writer) error
}
