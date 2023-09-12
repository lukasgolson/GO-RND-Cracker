package fileArray

import (
	"awesomeProject/internal/serialization"
	"fmt"
)

// Append appends an item of type T to a FileArray.
// It serializes the item and adds it to the end of the array.
//
// Parameters:
//   - item: The item of type T to be appended.
//
// Returns:
//   - An error if the operation fails, nil otherwise.
func (fileArray *FileArray[T]) Append(item T) (serialization.Offset, error) {

	id := fileArray.Count()
	err := fileArray.SetItemAtIndex(item, id)

	if err != nil {
		return 0, err
	}
	return id, nil
}

// SetItemAtIndex sets an item of type T at a specific index within a FileArray.
// It serializes the item and stores it at the given index.
//
// Parameters:
//   - item: The item of type T to be set.
//   - index: The index where the item will be placed.
//
// Returns:
//   - An error if the operation fails, nil otherwise.
func (fileArray *FileArray[T]) SetItemAtIndex(item T, index serialization.Offset) error {
	serializationSize := (item).StrideLength()

	if index > fileArray.Count() {
		return fmt.Errorf("index out of bounds. Max index %d", fileArray.Count())
	}

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

	err := (item).SerializeToBinaryStream(slice[memoryLocation : memoryLocation+serializationSize])
	if err != nil {
		return err
	}

	if index >= fileArray.Count() {
		fileArray.setCount(index + 1)
	}

	return nil
}

// GetItemFromIndex retrieves an item of type T from a FileArray at the specified index.
// It deserializes the item and returns it along with any potential errors.
//
// Parameters:
//   - index: The index of the item to retrieve.
//
// Returns:
//   - The item of type T at the specified index.
//   - An error if the operation fails, nil otherwise.
func (fileArray *FileArray[T]) GetItemFromIndex(index serialization.Offset) (T, error) {
	var err error
	var item T

	if index >= fileArray.Count() {
		return item, fmt.Errorf("index out of bounds")
	}

	serializedSize := item.StrideLength()

	memoryLocation := index * serializedSize

	slice := fileArray.getDataSlice()

	item, err = item.DeserializeFromBinaryStream(slice[memoryLocation : memoryLocation+serializedSize])
	if err != nil {
		return item, err
	}

	return item, nil
}

// ShrinkWrap reduces the size of a FileArray to match the actual data size.
// It is used to optimize disk usage by resizing the array to fit its contents.
//
// Returns:
//   - An error if the operation fails, nil otherwise.
func (fileArray *FileArray[T]) ShrinkWrap() error {
	var sampleItem T

	err := fileArray.shrinkFileSizeToDataSize(sampleItem.StrideLength())
	if err != nil {
		return err
	}

	return nil
}
