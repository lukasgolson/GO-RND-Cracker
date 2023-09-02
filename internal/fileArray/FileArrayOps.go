package fileArray

import (
	"awesomeProject/internal/serialization"
	"bytes"
	"fmt"
)

func AppendItem[T serialization.Serializer[T]](fileArray *FileArray, item T) error {
	err := SetItemAtIndex(fileArray, item, fileArray.Count())

	if err != nil {
		return err
	}
	return nil
}

func SetItemAtIndex[T serialization.Serializer[T]](fileArray *FileArray, item T, index uint64) error {
	serializationSize := (item).SerializedSize()

	if index > fileArray.Count() {
		return fmt.Errorf("index out of bounds. Max index %d", fileArray.Count())
	}

	var buffer bytes.Buffer
	err := (item).SerializeToBinaryStream(&buffer)
	if err != nil {
		return err
	}

	serializedItem := buffer.Bytes()

	arraySize := serializationSize * (index + 1)

	if !fileArray.hasSpace(arraySize) {
		err := fileArray.multiplyMemoryMapSize(2)
		if err != nil {
			return err
		}
	}

	memoryLocation := serializationSize * index //<-- Updated calculation for memory location

	slice := fileArray.getDataSlice()

	copy(slice[memoryLocation:memoryLocation+serializationSize], serializedItem)

	if index >= fileArray.Count() { //<-- Changed the condition to handle index equal to or greater than Count()
		fileArray.setCount(index + 1)
	}

	return nil
}

func GetItemFromIndex[T serialization.Serializer[T]](fileArray *FileArray, index uint64) (T, error) {
	var err error
	var item T

	if index >= fileArray.Count() {
		return item, fmt.Errorf("index out of bounds")
	}

	var buffer bytes.Buffer

	serializedSize := item.SerializedSize()

	memoryLocation := index * serializedSize

	slice := fileArray.getDataSlice()

	serializedItem := make([]byte, serializedSize)
	copy(serializedItem, slice[memoryLocation:memoryLocation+serializedSize])
	buffer.Write(serializedItem)

	item, err = item.DeserializeFromBinaryStream(&buffer)
	if err != nil {
		return item, err
	}

	return item, nil
}

func ShrinkWrapFileArray[T serialization.Serializer[T]](fileArray *FileArray) error {
	var sampleItem T

	err := fileArray.shrinkFileSizeToDataSize(sampleItem.SerializedSize())
	if err != nil {
		return err
	}

	return nil
}
