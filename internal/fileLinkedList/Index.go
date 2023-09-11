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

	buf := make([]byte, 8+8+8) // Create a buffer for itemID, offset and length

	binary.LittleEndian.PutUint64(buf[0:8], uint64(i.itemID))   // Convert int64 to little-endian binary and put it in the buffer
	binary.LittleEndian.PutUint64(buf[8:16], uint64(i.offset))  // Convert int64 to little-endian binary and put it in the buffer
	binary.LittleEndian.PutUint64(buf[16:24], uint64(i.length)) // Convert int64 to little-endian binary and put it in the buffer

	_, err := writer.Write(buf)
	if err != nil {
		return err
	}

	return nil
}

func (i indexEntry) DeserializeFromBinaryStream(reader io.Reader) (indexEntry, error) {

	buf := make([]byte, 8+8+8) // Create a buffer for itemID, offset and length
	_, err := io.ReadFull(reader, buf)
	if err != nil {
		return i, err
	}

	i.itemID = serialization.Offset(binary.LittleEndian.Uint64(buf[0:8]))   // Read the little-endian binary from the buffer and convert to offset
	i.offset = serialization.Offset(binary.LittleEndian.Uint64(buf[8:16]))  // Read the little-endian binary from the buffer and convert to offset
	i.length = serialization.Length(binary.LittleEndian.Uint64(buf[16:24])) // Read the little-endian binary from the buffer and convert to offset

	return i, nil
}

func (i indexEntry) StrideLength() serialization.Length {
	return 8 + 8 + 8
}

func (i indexEntry) IDByte() byte {
	return 'I'
}
