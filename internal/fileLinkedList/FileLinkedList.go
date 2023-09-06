package fileLinkedList

import (
	"awesomeProject/internal/fileArray"
	"awesomeProject/internal/serialization"
	"fmt"
)

type FileLinkedList struct {
	elementsArray fileArray.FileArray
	indexArray    fileArray.FileArray
}

func NewFileLinkedList[T serialization.Serializer[T]](filename string) (*FileLinkedList, error) {
	fileLinkedList := &FileLinkedList{}

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

func getOffsetFromListID(listID serialization.Offset) (bool, serialization.Offset, error) {

	// We will use the listID as the index for the index array.
	// The index array will contain the offset of the first element in the linked-list.
	// If the index array contains a zero value, then the linked-list is empty.
	// If the index array contains a non-zero value, then the linked-list is not empty.
	// If the ID is valid, return True and the offset.
	// If the ID is invalid, return False and 0.
	// Else return an error.

	//TODO implement me
	panic("implement me")
}

func setOffsetFromListID(listID serialization.Offset, offset serialization.Offset) error {

	//TODO implement me
	panic("implement me")
}

func Add[T serialization.Serializer[T]](l *FileLinkedList, listID serialization.Offset, item T) error {

	// We will lookup the offset of the first element in the linked-list from the index.

	// If the linked-list is empty, then we will append the item to the elements array and update the index.

	// If the linked-list is not empty, then we will append the item to the elements array and update the previous first element's next offset.
	// We will then update the index to point to the new first element.

	//TODO implement me
	panic("implement me")
}

func Remove[T serialization.Serializer[T]](l *FileLinkedList, listID serialization.Offset, item T) error {

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

func Get[T serialization.Serializer[T]](l *FileLinkedList, listID serialization.Offset, index uint64) (T, error) {

	// We will lookup the offset of the first element in the linked-list from the index.
	// We will then walk the linked-list until we reach the specified index.
	// If we reach the specified index, then we will return the item.
	// If we don't reach the specified index, then we will return an error.

	//TODO implement me
	panic("implement me")
}

func Contains[T serialization.Serializer[T]](l *FileLinkedList, listID serialization.Offset, item T) (bool, error) {

	// We will lookup the offset of the first element in the linked-list from the index.
	// We will then walk the linked-list until we find the item.
	// If we find the item, then we will return True.
	// If we don't find the item, then we will return False.
	// We will return an error if the linked-list is empty.

	//TODO implement me
	panic("implement me")
}
