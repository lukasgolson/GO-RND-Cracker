package fileLinkedList

import (
	"awesomeProject/internal/serialization"
	"encoding/binary"
	"io"
)

type linkedListNode[T serialization.Serializer[T]] struct {
	NextOffset serialization.Offset
	Item       T
}

func newLinkedListNode[T serialization.Serializer[T]](nextOffset serialization.Offset, item T) linkedListNode[T] {
	return linkedListNode[T]{NextOffset: nextOffset, Item: item}
}

func (l linkedListNode[T]) SerializeToBinaryStream(writer io.Writer) error {

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

func (l linkedListNode[T]) DeserializeFromBinaryStream(reader io.Reader) (linkedListNode[T], error) {

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

	return linkedListNode[T]{
		NextOffset: nextOffset,
		Item:       item,
	}, nil
}

func (l linkedListNode[T]) StrideLength() serialization.Length {

	return serialization.Length(binary.Size(l.NextOffset)) + l.Item.StrideLength()
}

func (l linkedListNode[T]) IDByte() byte {
	return 'L' + 'L' + 'N'
}
