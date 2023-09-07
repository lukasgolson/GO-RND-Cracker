package fileLinkedList

import (
	"awesomeProject/internal/fileArray"
	"awesomeProject/internal/serialization"
	"bytes"
	"fmt"
)

type FileLinkedList[T serialization.Serializer[T]] struct {
	elementsArray fileArray.FileArray
	indexArray    fileArray.FileArray
}

// NewFileLinkedList initializes a new instance of FileLinkedList[T] and its associated file arrays.
// It takes a filename as input and creates two file arrays: one for elements and one for index.
// Returns a pointer to the new FileLinkedList[T] and an error if any.
func NewFileLinkedList[T serialization.Serializer[T]](filename string) (*FileLinkedList[T], error) {
	fileLinkedList := &FileLinkedList[T]{}

	elementsFilename := fmt.Sprintf("%s.elements.bin", filename)
	indexFilename := fmt.Sprintf("%s.index.bin", filename)

	elementsArray, err := fileArray.NewFileArray(LinkedListNode[T]{}, elementsFilename)
	if err != nil {
		return nil, err
	}

	indexArray, err := fileArray.NewFileArray(LinkedListNode[T]{}, indexFilename)
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
func (list *FileLinkedList[T]) getBaseOffsetFromListID(listID serialization.Offset) (bool, serialization.Offset, error) {
	indexCount := list.indexArray.Count()

	if listID >= indexCount {
		return false, indexCount, nil
	}

	offset, err := fileArray.GetItemFromIndex[Index](&list.indexArray, listID)

	if err != nil {
		return false, 0, err
	}

	if offset.offset == serialization.MaxOffset() {
		return false, 0, nil
	}

	return true, offset.offset, nil
}

// setBaseOffsetOnListID sets the base offset for a given list ID.
// If the list ID is beyond the current index, it updates the index as well.
// Returns an error if any.
func (list *FileLinkedList[T]) setBaseOffsetOnListID(listID serialization.Offset, offset serialization.Offset) error {

	IndexCount := list.indexArray.Count()

	if listID > IndexCount {
		for i := IndexCount; i < listID; i++ {
			err := fileArray.SetItemAtIndex(&list.indexArray, NewIndex(i, serialization.MaxOffset()), i)
			if err != nil {
				return err
			}
		}
	}

	err := fileArray.SetItemAtIndex(&list.indexArray, NewIndex(listID, offset), listID)
	if err != nil {
		return err
	}

	return nil
}

// Add appends an item to the specified list ID.
// If the list doesn't exist, it creates a new list.
// Returns an error if any.
func (list *FileLinkedList[T]) Add(listID serialization.Offset, item T) error {
	listExists, listOffset, err := list.getBaseOffsetFromListID(listID)

	if err != nil {
		return err
	}

	var newOffset serialization.Offset

	if !listExists {

		newOffset, err = fileArray.Append(&list.elementsArray, NewLinkedListNode[T](serialization.MaxOffset(), item))
		if err != nil {
			return err
		}

	} else {

		currentHeadNode, err := fileArray.GetItemFromIndex[LinkedListNode[T]](&list.elementsArray, listOffset)
		if err != nil {
			return err
		}

		newOffset, err = fileArray.Append(&list.elementsArray, NewLinkedListNode[T](currentHeadNode.NextOffset, item))
		if err != nil {
			return err
		}
	}

	err = list.setBaseOffsetOnListID(listID, newOffset)
	if err != nil {
		return err
	}

	return nil
}

// Get retrieves an item from the specified list ID at the given index.
// Returns the item and an error if any. If the list or index is out of bounds, it returns an error.
func (list *FileLinkedList[T]) Get(listID serialization.Offset, index serialization.Offset) (T, error) {
	var item T

	listExists, baseOffset, err := list.getBaseOffsetFromListID(listID)

	if err != nil {
		return item, err
	}

	if !listExists {
		return item, fmt.Errorf("list does not exist")
	}

	currentOffset := baseOffset

	for i := serialization.Offset(0); i < index; i++ {
		currentNode, err := fileArray.GetItemFromIndex[LinkedListNode[T]](&list.elementsArray, currentOffset)
		if err != nil {
			return item, err
		}
		currentOffset = currentNode.NextOffset

		if currentOffset == serialization.MaxOffset() {
			return item, fmt.Errorf("index out of bounds")
		}
	}

	currentNode, err := fileArray.GetItemFromIndex[LinkedListNode[T]](&list.elementsArray, currentOffset)
	if err != nil {
		return item, err
	}

	return currentNode.Item, nil
}

// Remove removes an item from the specified list ID that matches the provided indexItem.
// Returns an error if the item is not found or if any other error occurs.
func (list *FileLinkedList[T]) Remove(listID serialization.Offset, indexItem T) error {

	listExists, baseOffset, err := list.getBaseOffsetFromListID(listID)

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

	previousOffset := baseOffset
	currentOffset := baseOffset
	nextOffset := baseOffset

	for nextOffset != serialization.MaxOffset() {

		currentOffset = nextOffset

		item, err := fileArray.GetItemFromIndex[LinkedListNode[T]](&list.elementsArray, nextOffset)
		if err != nil {
			return err
		}

		nextOffset = item.NextOffset

		if nextOffset == serialization.MaxOffset() {
			return fmt.Errorf("item not found")
		}

		var itemBuffer bytes.Buffer
		err = item.Item.SerializeToBinaryStream(&itemBuffer)

		if err != nil {
			return err
		}

		if bytes.Equal(indexBuffer.Bytes(), itemBuffer.Bytes()) {

			if currentOffset == baseOffset {
				err = list.setBaseOffsetOnListID(listID, nextOffset)
			} else {
				previousItem, err := fileArray.GetItemFromIndex[LinkedListNode[T]](&list.elementsArray, previousOffset)
				if err != nil {
					return err
				}

				previousItem.NextOffset = nextOffset
				err = fileArray.SetItemAtIndex[LinkedListNode[T]](&list.elementsArray, previousItem, previousOffset)
			}

			err = fileArray.SetItemAtIndex[LinkedListNode[T]](&list.elementsArray, NewLinkedListNode[T](serialization.MaxOffset(), nil), currentOffset)
			if err != nil {
				return err
			}

		}

		previousOffset = currentOffset

	}

	return nil
}

// Contains checks if the specified list ID contains the given item.
// Returns true if the item is found, false if not, and an error if any.
func (list *FileLinkedList[T]) Contains(listID serialization.Offset, item T) (bool, error) {

	listExists, baseOffset, err := list.getBaseOffsetFromListID(listID)

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

	nextOffset := baseOffset

	for nextOffset != serialization.MaxOffset() {

		item, err := fileArray.GetItemFromIndex[LinkedListNode[T]](&list.elementsArray, nextOffset)
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
