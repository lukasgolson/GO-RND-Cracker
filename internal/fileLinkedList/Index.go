package fileLinkedList

import (
	"awesomeProject/internal/serialization"
	"encoding/binary"
	"io"
)

type Index struct {
	itemID serialization.Offset
	offset serialization.Offset
	count  serialization.Length
}

func (i Index) SerializeToBinaryStream(writer io.Writer) error {
	err := binary.Write(writer, binary.LittleEndian, i.itemID)

	if err != nil {
		return err
	}

	err = binary.Write(writer, binary.LittleEndian, i.offset)

	if err != nil {
		return err
	}

	return nil
}

func (i Index) DeserializeFromBinaryStream(reader io.Reader) (Index, error) {

	err := binary.Read(reader, binary.LittleEndian, &i.itemID)

	if err != nil {
		return i, err
	}

	err = binary.Read(reader, binary.LittleEndian, &i.offset)

	if err != nil {
		return i, err
	}

	return i, nil
}

func (i Index) StrideLength() serialization.Length {
	return serialization.Length(binary.Size(i.itemID) + binary.Size(i.offset))
}

func (i Index) IDByte() byte {
	return 'I'
}
