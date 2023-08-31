package memory

import (
	"awesomeProject/internal/interfaces"
	"bytes"
	"fmt"
)

func AppendItem[T interfaces.Serializer[any]](fileArray *FileArray, item T) error {
	serializationSize := item.SerializedSize()

	var buffer bytes.Buffer
	err := item.SerializeToBinaryStream(&buffer)
	if err != nil {
		return err
	}

	serializedItem := buffer.Bytes()

	currentLength := fileArray.GetLength()

	slice := fileArray.getSlice()

	indexToOverwrite := currentLength * serializationSize

	copy(slice[indexToOverwrite:indexToOverwrite+serializationSize], serializedItem)

	fileArray.incrementLength()

	return nil
}

func SetItemAtIndex[T interfaces.Serializer[any]](fileArray *FileArray, item T, index uint64) error {
	serializationSize := item.SerializedSize()

	var buffer bytes.Buffer
	err := item.SerializeToBinaryStream(&buffer)
	if err != nil {
		return err
	}

	serializedItem := buffer.Bytes()

	slice := fileArray.getSlice()

	memoryLocation := index * serializationSize

	copy(slice[memoryLocation:memoryLocation+serializationSize], serializedItem)

	fileArray.incrementLength()

	return nil
}

func GetItemFromIndex[T interfaces.Serializer[any]](fileArray *FileArray, index uint64) (interfaces.Serializer[any], error) {
	var item T
	var buffer bytes.Buffer

	if index > fileArray.GetLength() {
		return nil, fmt.Errorf("index out of bounds")
	}

	serializedSize := item.SerializedSize()

	memoryLocation := index * serializedSize

	slice := fileArray.getSlice()

	serializedItem := make([]byte, serializedSize)
	copy(serializedItem, slice[memoryLocation:memoryLocation+serializedSize])
	buffer.Write(serializedItem)

	itemTemp, err := item.DeserializeFromBinaryStream(&buffer)
	if err != nil {
		return nil, err
	}

	item = itemTemp.(T)

	return item, nil
}

func ShrinkwrapFile[T interfaces.Serializer[any]](fileArray *FileArray) error {
	var sampleItem T

	err := fileArray.shrinkFileSizeToDataSize(sampleItem.SerializedSize())
	if err != nil {
		return err
	}

	return nil
}
