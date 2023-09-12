package fileLinkedList

import (
	"awesomeProject/internal/serialization"
	"encoding/binary"
)

type indexEntry struct {
	itemID serialization.Offset
	offset serialization.Offset
	length serialization.Length
}

func newIndexEntry(itemID serialization.Offset, offset serialization.Offset, length serialization.Length) indexEntry {
	return indexEntry{itemID: itemID, offset: offset, length: length}
}

func (i indexEntry) SerializeToBinaryStream(buf []byte) error {

	binary.LittleEndian.PutUint64(buf[0:8], uint64(i.itemID))   // Convert int64 to little-endian binary and put it in the buffer
	binary.LittleEndian.PutUint64(buf[8:16], uint64(i.offset))  // Convert int64 to little-endian binary and put it in the buffer
	binary.LittleEndian.PutUint64(buf[16:24], uint64(i.length)) // Convert int64 to little-endian binary and put it in the buffer

	return nil
}

func (i indexEntry) DeserializeFromBinaryStream(buf []byte) (indexEntry, error) {

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
