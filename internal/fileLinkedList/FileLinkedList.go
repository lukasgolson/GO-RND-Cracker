package fileLinkedList

import (
	"awesomeProject/internal/fileArray"
	"awesomeProject/internal/serialization"
	"bytes"
	"fmt"
)

type FileLinkedList[T serialization.Serializer[T]] struct {
	elementsArray fileArray.FileArray[LinkedListNode[T]]
	indexArray    fileArray.FileArray[Index]
}

// NewFileLinkedList initializes a new instance of FileLinkedList[T] and its associated file arrays.
// It takes a filename as input and creates two file arrays: one for elements and one for index.
// Returns a pointer to the new FileLinkedList[T] and an error if any.
func NewFileLinkedList[T serialization.Serializer[T]](filename string) (*FileLinkedList[T], error) {
	fileLinkedList := &FileLinkedList[T]{}

	elementsFilename := fmt.Sprintf("%s.elements.bin", filename)
	indexFilename := fmt.Sprintf("%s.index.bin", filename)

	elementsArray, err := fileArray.NewFileArray[LinkedListNode[T]](elementsFilename)
	if err != nil {
		return nil, err
	}

	indexArray, err := fileArray.NewFileArray[Index](indexFilename)
	if err != nil {
		return nil, err
	}

	fileLinkedList.elementsArray = *elementsArray
	fileLinkedList.indexArray = *indexArray

	return fileLinkedList, nil
}

// getBaseOffsetFromListID retrieves the base offset associated with a list ID.
// It checks if the list ID is within bounds and returns a boolean indicating existence,
// the base offset, and an error if any.
func (list *FileLinkedList[T]) getBaseOffsetFromListID(listID serialization.Offset) (bool, Index, error) {
	indexCount := list.indexArray.Count()

	if indexCount == 0 {
		return false, *new(Index), nil
	}

	if listID >= indexCount {
		return false, *new(Index), nil
	}

	indexEntry, err := list.indexArray.GetItemFromIndex(listID)

	if err != nil {
		return false, *new(Index), err
	}

	if indexEntry.offset == serialization.MaxOffset() {
		return false, *new(Index), nil
	}

	return true, indexEntry, nil
}

// setBaseOffsetOnListID sets the base offset for a given list ID.
// If the list ID is beyond the current index, it updates the index as well.
// Returns an error if any.
func (list *FileLinkedList[T]) setBaseOffsetOnListID(listID serialization.Offset, offset serialization.Offset, length serialization.Length) error {
	indexCount := list.indexArray.Count()

	if listID >= indexCount {
		numItemsToAdd := listID - indexCount + 1

		for i := indexCount; i < indexCount+numItemsToAdd; i++ {
			newIndex := NewIndex(i, serialization.MaxOffset(), 0)

			err := list.indexArray.SetItemAtIndex(newIndex, i)
			if err != nil {
				return err
			}
		}
	}

	err := list.indexArray.SetItemAtIndex(NewIndex(listID, offset, length), listID)
	if err != nil {
		return err
	}

	return nil
}

// Add appends an item to the specified list ID.
// If the list doesn't exist, it creates a new list.
// Returns an error if any.
func (list *FileLinkedList[T]) Add(listID serialization.Offset, item T) error {
	listExists, indexEntry, err := list.getBaseOffsetFromListID(listID)

	if err != nil {
		return err
	}

	var newOffset serialization.Offset

	if !listExists {

		newOffset, err = list.elementsArray.Append(NewLinkedListNode[T](serialization.MaxOffset(), item))
		if err != nil {
			return err
		}

	} else {

		//currentHeadNode, err := fileArray.GetItemFromIndex[LinkedListNode[T]](&list.elementsArray, indexEntry.offset)

		newOffset, err = list.elementsArray.Append(NewLinkedListNode[T](indexEntry.offset, item))
		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	err = list.setBaseOffsetOnListID(listID, newOffset, indexEntry.length+1)
	if err != nil {
		return err
	}

	return nil
}

// Get retrieves an item from the specified list ID at the given index.
// Returns the item and an error if any. If the list or index is out of bounds, it returns an error.
func (list *FileLinkedList[T]) Get(listID serialization.Offset, index serialization.Offset) (T, error) {
	var item T
	listExists, indexEntry, err := list.getBaseOffsetFromListID(listID)

	if err != nil {
		return item, fmt.Errorf("list index error" + err.Error())
	}

	if !listExists {
		return item, fmt.Errorf("list does not exist")
	}

	index = indexEntry.length - index

	indexCounter := serialization.Offset(0)

	for nextOffset := indexEntry.offset; nextOffset != serialization.MaxOffset(); {
		currentNode, err := list.elementsArray.GetItemFromIndex(nextOffset)
		if err != nil {
			return item, err
		}

		nextOffset = currentNode.NextOffset

		indexCounter++

		if indexCounter == index {
			return currentNode.Item, nil
		}
	}

	return item, fmt.Errorf("index out of bounds test")
}

// Remove removes an item from the specified list ID that matches the provided indexItem.
// Returns an error if the item is not found or if any other error occurs.
func (list *FileLinkedList[T]) Remove(listID serialization.Offset, indexItem T) error {

	listExists, indexEntry, err := list.getBaseOffsetFromListID(listID)

	if err != nil {
		return err
	}

	if !listExists {
		return fmt.Errorf("list does not exist")
	}

	var indexBuffer bytes.Buffer

	err = indexItem.SerializeToBinaryStream(&indexBuffer)
	if err != nil {
		return err
	}

	previousOffset := indexEntry.offset
	currentOffset := indexEntry.offset
	nextOffset := indexEntry.offset

	for nextOffset != serialization.MaxOffset() {

		currentOffset = nextOffset

		item, err := list.elementsArray.GetItemFromIndex(nextOffset)
		if err != nil {
			return err
		}

		nextOffset = item.NextOffset

		var itemBuffer bytes.Buffer
		err = item.Item.SerializeToBinaryStream(&itemBuffer)

		if err != nil {
			return err
		}

		if bytes.Equal(indexBuffer.Bytes(), itemBuffer.Bytes()) {

			if currentOffset == indexEntry.offset {

				if nextOffset == serialization.MaxOffset() {
					err = list.setBaseOffsetOnListID(listID, serialization.MaxOffset(), 0)
					if err != nil {
						return err
					}
					return nil
				} else {
					err = list.setBaseOffsetOnListID(listID, nextOffset, indexEntry.length-1)
					if err != nil {
						return err
					}
				}
			} else {
				previousItem, err := list.elementsArray.GetItemFromIndex(previousOffset)
				if err != nil {
					return err
				}

				previousItem.NextOffset = nextOffset
				err = list.elementsArray.SetItemAtIndex(previousItem, previousOffset)
			}

			err = list.elementsArray.SetItemAtIndex(NewLinkedListNode[T](serialization.MaxOffset(), *new(T)), currentOffset)
			if err != nil {
				return err
			}

		}

		previousOffset = currentOffset

		if nextOffset == serialization.MaxOffset() {
			return fmt.Errorf("item not found")
		}

	}

	return nil
}

// Contains checks if the specified list ID contains the given item.
// Returns true if the item is found, false if not, and an error if any.
func (list *FileLinkedList[T]) Contains(listID serialization.Offset, item T) (bool, error) {

	listExists, indexEntry, err := list.getBaseOffsetFromListID(listID)

	if err != nil {
		return false, err
	}

	if !listExists {
		return false, fmt.Errorf("list does not exist")
	}

	itemBuffer := bytes.Buffer{}
	err = item.SerializeToBinaryStream(&itemBuffer)
	if err != nil {
		return false, err
	}

	nextOffset := indexEntry.offset

	for nextOffset != serialization.MaxOffset() {

		item, err := list.elementsArray.GetItemFromIndex(nextOffset)
		if err != nil {
			return false, err
		}

		nextOffset = item.NextOffset

		currentItemBuffer := bytes.Buffer{}
		err = item.Item.SerializeToBinaryStream(&currentItemBuffer)
		if err != nil {
			return false, err
		}

		if bytes.Equal(itemBuffer.Bytes(), currentItemBuffer.Bytes()) {
			return true, nil
		}

	}

	return false, nil
}
