package fileArray

import (
	"awesomeProject/internal/serialization"
	"bytes"
	"fmt"
)

func AppendItem[T serialization.Serializer](fileArray *FileArray, item T) error {
	err := SetItemAtIndex(fileArray, item, fileArray.Count())
	if err != nil {
		return err
	}
	return nil
}

func SetItemAtIndex[T serialization.Serializer](fileArray *FileArray, item T, index uint64) error {
	serializationSize := item.SerializedSize()

	if index > fileArray.Count() {
		return fmt.Errorf("index out of bounds. Max index %d", fileArray.Count())
	}

	var buffer bytes.Buffer
	err := item.SerializeToBinaryStream(&buffer)
	if err != nil {
		return err
	}

	serializedItem := buffer.Bytes()

	slice := fileArray.getSlice()

	arraySize := serializationSize * (index + 1)

	if !fileArray.hasSpace(arraySize) {
		err := fileArray.multiplyFileSize(2)
		if err != nil {
			return err
		}
	}

	memoryLocation := serializationSize * index //<-- Updated calculation for memory location

	copy(slice[memoryLocation:memoryLocation+serializationSize], serializedItem)

	if index >= fileArray.Count() { //<-- Changed the condition to handle index equal to or greater than Count()
		fileArray.setCount(index + 1)
	}

	return nil
}

func GetItemFromIndex[T serialization.Serializer](fileArray *FileArray, index uint64) (serialization.Serializer, error) {
	var err error
	var item T
	var buffer bytes.Buffer

	if index > fileArray.Count() {
		return nil, fmt.Errorf("index out of bounds")
	}

	serializedSize := item.SerializedSize()

	memoryLocation := index * serializedSize

	slice := fileArray.getSlice()

	serializedItem := make([]byte, serializedSize)
	copy(serializedItem, slice[memoryLocation:memoryLocation+serializedSize])
	buffer.Write(serializedItem)

	err = item.DeserializeFromBinaryStream(&buffer)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func ShrinkwrapFile[T serialization.Serializer](fileArray *FileArray) error {
	var sampleItem T

	err := fileArray.shrinkFileSizeToDataSize(sampleItem.SerializedSize())
	if err != nil {
		return err
	}

	return nil
}
