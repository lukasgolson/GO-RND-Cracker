package fileLinkedList

import (
	"awesomeProject/internal/serialization"
	"encoding/binary"
	"io"
)

type LinkedListNode[T serialization.Serializer[T]] struct {
	NextOffset serialization.Offset
	Item       T
}

func (l LinkedListNode[T]) SerializeToBinaryStream(writer io.Writer) error {

	err := binary.Write(writer, binary.LittleEndian, l.NextOffset)
	if err != nil {
		return err
	}

	err = l.Item.SerializeToBinaryStream(writer)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

func (l LinkedListNode[T]) DeserializeFromBinaryStream(reader io.Reader) (LinkedListNode[T], error) {

	var nextOffset serialization.Offset
	err := binary.Read(reader, binary.LittleEndian, &nextOffset)
	if err != nil {
		return l, err
	}

	var item T
	item, err = T.DeserializeFromBinaryStream(item, reader)
	if err != nil {
		return l, err
	}

	return LinkedListNode[T]{
		NextOffset: nextOffset,
		Item:       item,
	}, nil
}

func (l LinkedListNode[T]) StrideLength() serialization.Length {

	return serialization.Length(binary.Size(l.NextOffset)) + l.Item.StrideLength()
}

func (l LinkedListNode[T]) IDByte() byte {
	return 'L' + 'L' + 'N'
}
