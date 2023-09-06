package fileIndex

import (
	"encoding/binary"
	"io"
)

type Index struct {
	itemID int64
	offset int64
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

func (i Index) StrideLength() uint64 {
	return uint64(binary.Size(i.itemID) + binary.Size(i.offset))
}

func (i Index) IDByte() byte {
	return 'I'
}
