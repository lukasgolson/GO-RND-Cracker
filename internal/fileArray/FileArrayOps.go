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
func Append[T serialization.Serializer[T]](fileArray *FileArray, item T) error {
	err := SetItemAtIndex(fileArray, item, fileArray.Count())

	if err != nil {
		return err
	}
	return nil
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
func SetItemAtIndex[T serialization.Serializer[T]](fileArray *FileArray, item T, index uint64) error {
	serializationSize := (item).StrideLength()

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
func GetItemFromIndex[T serialization.Serializer[T]](fileArray *FileArray, index uint64) (T, error) {
	var err error
	var item T

	if index >= fileArray.Count() {
		return item, fmt.Errorf("index out of bounds")
	}

	var buffer bytes.Buffer

	serializedSize := item.StrideLength()

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

	err := fileArray.shrinkFileSizeToDataSize(sampleItem.StrideLength())
	if err != nil {
		return err
	}

	return nil
}
