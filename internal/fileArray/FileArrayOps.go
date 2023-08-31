package fileArray

import (
	"awesomeProject/internal/interfaces"
	"bytes"
	"fmt"
)

func AppendItem[T interfaces.Serializer[T]](fileArray *FileArray, item T) error {
	err := SetItemAtIndex(fileArray, item, fileArray.Count())
	if err != nil {
		return err
	}
	return nil
}

func SetItemAtIndex[T interfaces.Serializer[T]](fileArray *FileArray, item T, index uint64) error {
	serializationSize := item.SerializedSize()

	if index > fileArray.Count() {
		return fmt.Errorf("index out of bounds")
	}

	var buffer bytes.Buffer
	err := item.SerializeToBinaryStream(&buffer)
	if err != nil {
		return err
	}

	serializedItem := buffer.Bytes()

	slice := fileArray.getSlice()

	arraySize := serializationSize * index

	if index == fileArray.Count() {
		arraySize += serializationSize
	}

	if fileArray.hasSpace(arraySize) {
		err := fileArray.adjustFileSize(2)
		if err != nil {
			return err
		}
	}

	memoryLocation := index * serializationSize

	copy(slice[memoryLocation:memoryLocation+serializationSize], serializedItem)

	if index == fileArray.Count() {
		fileArray.setCount(index + 1)
	}

	return nil
}

func GetItemFromIndex[T interfaces.Serializer[T]](fileArray *FileArray, index uint64) (interfaces.Serializer[T], error) {
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

	item, err = item.DeserializeFromBinaryStream(&buffer)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func ShrinkwrapFile[T interfaces.Serializer[T]](fileArray *FileArray) error {
	var sampleItem T

	err := fileArray.shrinkFileSizeToDataSize(sampleItem.SerializedSize())
	if err != nil {
		return err
	}

	return nil
}
