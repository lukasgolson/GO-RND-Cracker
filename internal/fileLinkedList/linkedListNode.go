package fileLinkedList

import (
	"awesomeProject/internal/serialization"
	"encoding/binary"
)

type linkedListNode[T serialization.Serializer[T]] struct {
	NextOffset serialization.Offset
	Item       T
}

func newLinkedListNode[T serialization.Serializer[T]](nextOffset serialization.Offset, item T) linkedListNode[T] {
	return linkedListNode[T]{NextOffset: nextOffset, Item: item}
}

func (l linkedListNode[T]) SerializeToBinaryStream(buf []byte) error {

	binary.LittleEndian.PutUint64(buf[0:8], uint64(l.NextOffset)) // Convert int64 to little-endian binary and put it in the buffer

	err := l.Item.SerializeToBinaryStream(buf[8:])
	if err != nil {
		return err
	}

	return nil
}

func (l linkedListNode[T]) DeserializeFromBinaryStream(buf []byte) (linkedListNode[T], error) {

	var err error

	l.NextOffset = serialization.Offset(binary.LittleEndian.Uint64(buf[0:8]))

	l.Item, err = T.DeserializeFromBinaryStream(l.Item, buf[8:])

	return l, err
}

func (l linkedListNode[T]) StrideLength() serialization.Length {

	return 8 + l.Item.StrideLength()
}

func (l linkedListNode[T]) IDByte() byte {
	return 'L' + 'L' + 'N'
}
