package fileLinkedList

import (
	"awesomeProject/internal/serialization"
	"encoding/binary"
	"io"
)

type indexEntry struct {
	itemID serialization.Offset
	offset serialization.Offset
	length serialization.Length
}

func newIndexEntry(itemID serialization.Offset, offset serialization.Offset, length serialization.Length) indexEntry {
	return indexEntry{itemID: itemID, offset: offset, length: length}
}

func (i indexEntry) SerializeToBinaryStream(writer io.Writer) error {
	err := binary.Write(writer, binary.LittleEndian, i.itemID)

	if err != nil {
		return err
	}

	err = binary.Write(writer, binary.LittleEndian, i.offset)

	if err != nil {
		return err
	}

	err = binary.Write(writer, binary.LittleEndian, i.length)

	if err != nil {
		return err
	}

	return nil
}

func (i indexEntry) DeserializeFromBinaryStream(reader io.Reader) (indexEntry, error) {

	err := binary.Read(reader, binary.LittleEndian, &i.itemID)

	if err != nil {
		return i, err
	}

	err = binary.Read(reader, binary.LittleEndian, &i.offset)

	if err != nil {
		return i, err
	}

	err = binary.Read(reader, binary.LittleEndian, &i.length)

	if err != nil {
		return i, err
	}

	return i, nil
}

func (i indexEntry) StrideLength() serialization.Length {
	return serialization.Length(binary.Size(i.itemID) + binary.Size(i.offset) + binary.Size(i.length))
}

func (i indexEntry) IDByte() byte {
	return 'I'
}
