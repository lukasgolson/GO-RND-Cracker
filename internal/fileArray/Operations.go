package fileArray

import (
	"awesomeProject/internal/serialization"
	"bytes"
	"fmt"
)

// Append appends an item of type T to a FileArray.
// It serializes the item and adds it to the end of the array.
//
// Parameters:
//   - fileArray: The target FileArray to which the item is appended.
//   - item: The item of type T to be appended.
//
// Returns:
//   - An error if the operation fails, nil otherwise.
func Append[T serialization.Serializer[T]](fileArray *FileArray, item T) (serialization.Offset, error) {

	id := fileArray.Count()
	err := SetItemAtIndex[T](fileArray, item, serialization.Offset(id))

	if err != nil {
		return 0, err
	}
	return serialization.Offset(id), nil
}

// SetItemAtIndex sets an item of type T at a specific index within a FileArray.
// It serializes the item and stores it at the given index.
//
// Parameters:
//   - fileArray: The target FileArray in which the item is set.
//   - item: The item of type T to be set.
//   - index: The index where the item will be placed.
//
// Returns:
//   - An error if the operation fails, nil otherwise.
func SetItemAtIndex[T serialization.Serializer[T]](fileArray *FileArray, item T, index serialization.Offset) error {
	serializationSize := serialization.Offset((item).StrideLength())

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

	stallCounter := 0
	for !fileArray.hasSpace(uint64(arraySize)) {
		stallCounter++
		err := fileArray.multiplyMemoryMapSize(2)
		if err != nil {
			return err
		}

		if stallCounter > 32 {
			return fmt.Errorf("stall counter exceeded")
		}
	}

	memoryLocation := serializationSize * index

	slice := fileArray.getDataSlice()

	copy(slice[memoryLocation:memoryLocation+serializationSize], serializedItem)

	if index >= fileArray.Count() {
		fileArray.setCount(index + 1)
	}

	return nil
}

// GetItemFromIndex retrieves an item of type T from a FileArray at the specified index.
// It deserializes the item and returns it along with any potential errors.
//
// Parameters:
//   - fileArray: The source FileArray from which the item is retrieved.
//   - index: The index of the item to retrieve.
//
// Returns:
//   - The item of type T at the specified index.
//   - An error if the operation fails, nil otherwise.
func GetItemFromIndex[T serialization.Serializer[T]](fileArray *FileArray, index serialization.Offset) (T, error) {
	var err error
	var item T

	if index >= fileArray.Count() {
		return item, fmt.Errorf("index out of bounds")
	}

	serializedSize := serialization.Offset(item.StrideLength())

	memoryLocation := index * serializedSize

	slice := fileArray.getDataSlice()

	var buffer bytes.Buffer
	serializedItem := make([]byte, serializedSize)
	copy(serializedItem, slice[memoryLocation:memoryLocation+serializedSize])
	buffer.Write(serializedItem)

	item, err = item.DeserializeFromBinaryStream(&buffer)
	if err != nil {
		return item, err
	}

	return item, nil
}

// ShrinkWrapFileArray reduces the size of a FileArray to match the actual data size.
// It is used to optimize disk usage by resizing the array to fit its contents.
//
// Parameters:
//   - fileArray: The FileArray to be shrink-wrapped.
//
// Returns:
//   - An error if the operation fails, nil otherwise.
func ShrinkWrapFileArray[T serialization.Serializer[T]](fileArray *FileArray) error {
	var sampleItem T

	err := fileArray.shrinkFileSizeToDataSize(serialization.Length(sampleItem.StrideLength()))
	if err != nil {
		return err
	}

	return nil
}
