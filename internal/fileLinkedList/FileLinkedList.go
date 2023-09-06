package fileLinkedList

import (
	"awesomeProject/internal/fileArray"
	"awesomeProject/internal/serialization"
	"fmt"
)

type FileLinkedList[T serialization.Serializer[T]] struct {
	elementsArray fileArray.FileArray
	indexArray    fileArray.FileArray
}

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

func (list *FileLinkedList[T]) getOffsetFromListID(listID serialization.Offset) (bool, serialization.Offset, error) {

	// We will use the listID as the index for the index array.
	// The index array will contain the offset of the first element in the linked-list.
	// If the index array contains a zero value, then the linked-list is empty.
	// If the index array contains a non-zero value, then the linked-list is not empty.
	// If the ID is valid, return True and the offset.
	// If the ID is invalid, return False and 0.

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

func (list *FileLinkedList[T]) setOffsetFromListID(listID serialization.Offset, offset serialization.Offset) error {

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

func (list *FileLinkedList[T]) Add(listID serialization.Offset, item T) error {

	// We will lookup the offset of the first element in the linked-list from the index.

	// If the linked-list is empty, then we will append the item to the elements array and update the index.

	// If the linked-list is not empty, then we will append the item to the elements array and update the previous first element's next offset.
	// We will then update the index to point to the new first element.

	listExists, listOffset, err := list.getOffsetFromListID(listID)

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

	err = list.setOffsetFromListID(listID, newOffset)
	if err != nil {
		return err
	}

	return nil
}

func (list *FileLinkedList[T]) Get(listID serialization.Offset, index serialization.Offset) (T, error) {

	// We will lookup the offset of the first element in the linked-list from the index.
	// We will then walk the linked-list until we reach the specified index.
	// If we reach the specified index, then we will return the item.
	// If we don't reach the specified index, then we will return an error.

	var item T

	listExists, listOffset, err := list.getOffsetFromListID(listID)

	if err != nil {
		return item, err
	}

	if !listExists {
		return item, fmt.Errorf("list does not exist")
	}

	currentOffset := listOffset

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

func (list *FileLinkedList[T]) Remove(listID serialization.Offset, index serialization.Offset) error {

	// We will lookup the offset of the first element in the linked-list from the index.
	// Then we will iterate through the linked-list until we find the item.
	// If we find the item, then we will update the previous element's next offset to point to the next element.
	// We will then update the index to point to the new first element.
	// We will then remove the element from the elements array.
	// If the linked-list is empty, then we will update the index to point to 0.
	// Eventually we will maintain a free list to reuse deleted elements.
	// If we don't find the item, then we will return an error.

	//TODO implement me
	panic("implement me")
}

func (list *FileLinkedList[T]) Contains(listID serialization.Offset, item T) (bool, error) {

	// We will lookup the offset of the first element in the linked-list from the index.
	// We will then walk the linked-list until we find the item.
	// If we find the item, then we will return True.
	// If we don't find the item, then we will return False.
	// We will return an error if the linked-list is empty.

	//TODO implement me
	panic("implement me")
}
