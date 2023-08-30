package bktree

import "io"

type StreamSerializer interface {
	SerializeToBinaryStream(writer io.Writer) error
}
