package fileLinkedList

import (
	"awesomeProject/internal/number"
	"awesomeProject/internal/serialization"
	"os"
	"testing"
)

func TestFileLinkedList_Add(t *testing.T) {

	tmpFile, err := os.CreateTemp("", "test-file")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	defer func(tempFile *os.File) {
		err := tempFile.Close()
		if err != nil {
			t.Fatalf("Failed to close temporary file: %v", err)
		}
	}(tmpFile)
	defer os.Remove(tmpFile.Name())

	// Create a new FileLinkedList and perform the Add operation.
	list, _ := NewFileLinkedList[number.Number](tmpFile.Name(), false)

	err = list.Add(0, *new(number.Number))
	if err != nil {
		t.Errorf("Add() failed: %v", err)
	}

	err = list.Add(0, *new(number.Number))
	if err != nil {
		t.Errorf("Add() failed: %v", err)
	}

	err = list.Add(1, *new(number.Number))
	if err != nil {
		t.Errorf("Add() failed: %v", err)
	}

	err = list.Add(1, *new(number.Number))
	if err != nil {
		t.Errorf("Add() failed: %v", err)
	}
}

func TestFileLinkedList_AddNonContiguous(t *testing.T) {

	tmpFile, err := os.CreateTemp("", "test-file")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	defer func(tempFile *os.File) {
		err := tempFile.Close()
		if err != nil {
			t.Fatalf("Failed to close temporary file: %v", err)
		}
	}(tmpFile)
	defer os.Remove(tmpFile.Name())

	// Create a new FileLinkedList and perform the Add operation.
	list, _ := NewFileLinkedList[number.Number](tmpFile.Name(), false)

	err = list.Add(0, *new(number.Number))
	if err != nil {
		t.Errorf("Add() failed: %v", err)
	}

	err = list.Add(10, *new(number.Number))
	if err != nil {
		t.Errorf("Add() failed: %v", err)
	}

	valid, indexEntry, err := list.getBaseOffsetFromListID(0)
	if err != nil {
		t.Errorf("getBaseOffsetFromListID() failed: %v", err)
	}

	if !valid {
		t.Errorf("getBaseOffsetFromListID() returned incorrect valid value: got %v, expected %v", valid, true)
	}

	if indexEntry.offset != 0 {
		t.Errorf("getBaseOffsetFromListID() returned incorrect baseOffset value: got %v, expected %v", indexEntry.offset, 0)
	}

	for i := 1; i < 10; i++ {
		valid, _, err := list.getBaseOffsetFromListID(serialization.Offset(i))
		if err != nil {
			t.Errorf("getBaseOffsetFromListID() failed: %v", err)
		}

		if valid {
			t.Errorf("getBaseOffsetFromListID() returned incorrect valid value: got %v, expected %v", valid, false)
		}
	}

	valid, indexEntry, err = list.getBaseOffsetFromListID(10)
	if err != nil {
		t.Errorf("getBaseOffsetFromListID() failed: %v", err)
	}

	if !valid {
		t.Errorf("getBaseOffsetFromListID() returned incorrect valid value: got %v, expected %v", valid, true)
	}

	if indexEntry.offset != 1 {
		t.Errorf("getBaseOffsetFromListID() returned incorrect baseOffset value: got %v, expected %v", indexEntry.offset, 1)
	}

}

func TestFileLinkedList_GetOne(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-file")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	defer func(tempFile *os.File) {
		err := tempFile.Close()
		if err != nil {
			t.Fatalf("Failed to close temporary file: %v", err)
		}
	}(tmpFile)
	defer os.Remove(tmpFile.Name())

	// Create a new FileLinkedList and perform the Add operation.
	list, err := NewFileLinkedList[number.Number](tmpFile.Name(), false)
	if err != nil {
		t.Errorf("NewFileLinkedList() failed: %v", err)
	}

	testVal := number.NewNumber(24)

	err = list.Add(0, testVal)
	if err != nil {
		t.Errorf("Add() failed: %v", err)
	}

	item, err := list.Get(0, 0)
	if err != nil {
		t.Errorf("Get() failed: %v", err)
	}

	if item != testVal {
		t.Errorf("Get() returned incorrect item: got %v, expected %v", item, testVal)
	}

}

func TestFileLinkedList_GetMultiple(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-file")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	defer func(tempFile *os.File) {
		err := tempFile.Close()
		if err != nil {
			t.Fatalf("Failed to close temporary file: %v", err)
		}
	}(tmpFile)
	defer os.Remove(tmpFile.Name())

	// Create a new FileLinkedList and perform the Add operation.
	list, _ := NewFileLinkedList[number.Number](tmpFile.Name(), false)

	for i := 0; i < 10; i++ {
		err := list.Add(0, number.NewNumber(int64(i)))
		if err != nil {
			t.Fatalf("Add() failed: %v", err)
		}
	}

	for i := 0; i < 10; i++ {
		item, err := list.Get(0, serialization.Offset(i))

		if err != nil {
			t.Errorf("Get() failed: %v", err)
		}
		expected := number.NewNumber(int64(i))
		if item != expected {
			t.Errorf("Get() returned incorrect item: got %v, expected %v", item, expected)
		}
	}

}

func TestFileLinkedList_SetGetListID(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-file")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	defer func(tempFile *os.File) {
		err := tempFile.Close()
		if err != nil {
			t.Fatalf("Failed to close temporary file: %v", err)
		}
	}(tmpFile)
	defer os.Remove(tmpFile.Name())

	// Create a new FileLinkedList and perform the Add operation.
	list, _ := NewFileLinkedList[number.Number](tmpFile.Name(), false)

	testOffsetValue := serialization.Offset(24)

	for i := 0; i < 15; i++ {
		err = list.setBaseOffsetOnListID(0, testOffsetValue, serialization.Length(i))
		if err != nil {
			return
		}

		valid, indexEntry, err := list.getBaseOffsetFromListID(0)

		if err != nil {
			return
		}

		if !valid {
			t.Errorf("setBaseOffsetOnListID() and getBaseOffsetFromListID() returned incorrect valid value: got %v, expected %v", valid, true)
		}

		if indexEntry.offset != testOffsetValue {
			t.Errorf("setBaseOffsetOnListID() and getBaseOffsetFromListID() returned incorrect offset value: got %v, expected %v", indexEntry.offset, testOffsetValue)
		}

		if indexEntry.length != serialization.Length(i) {
			t.Errorf("setBaseOffsetOnListID() and getBaseOffsetFromListID() returned incorrect length value: got %v, expected %v", indexEntry.length, 1)
		}
	}

}

func TestFileLinkedList_GetInvalidListID(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-file")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	defer func(tempFile *os.File) {
		err := tempFile.Close()
		if err != nil {
			t.Fatalf("Failed to close temporary file: %v", err)
		}
	}(tmpFile)
	defer os.Remove(tmpFile.Name())

	// Create a new FileLinkedList and perform the Add operation.
	list, _ := NewFileLinkedList[number.Number](tmpFile.Name(), false)

	exists, _, err := list.getBaseOffsetFromListID(0)

	if err != nil {
		t.Errorf("getBaseOffsetFromListID() failed: %v", err)
	}

	if exists {
		t.Errorf("getBaseOffsetFromListID() returned incorrect exists value: got %v, expected %v", exists, false)
	}
}

func TestFileLinkedList_Contains(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-file")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	defer func(tempFile *os.File) {
		err := tempFile.Close()
		if err != nil {
			t.Fatalf("Failed to close temporary file: %v", err)
		}
	}(tmpFile)
	defer os.Remove(tmpFile.Name())

	// Create a new FileLinkedList and perform the Add operation.
	list, err := NewFileLinkedList[number.Number](tmpFile.Name(), false)
	if err != nil {
		t.Errorf("NewFileLinkedList() failed: %v", err)
	}

	testItem := number.NewNumber(24)

	for i := 0; i < 15; i++ {
		err := list.Add(0, number.NewNumber(int64(i)))
		if err != nil {
			t.Errorf("Add() failed: %v", err)
		}
	}

	err = list.Add(0, testItem)
	if err != nil {
		t.Errorf("Add() failed: %v", err)
	}

	contains, err := list.Contains(0, testItem)
	if err != nil {
		t.Errorf("Contains() failed: %v", err)
	}

	if contains != true {
		t.Errorf("Contains() returned incorrect value: got %v, expected %v", contains, true)
	}

	contains, err = list.Contains(0, number.NewNumber(100))
	if err != nil {
		t.Errorf("Contains() failed: %v", err)
	}

	if contains != false {
		t.Errorf("Contains() returned incorrect value: got %v, expected %v", contains, false)
	}

	contains, err = list.Contains(1, testItem)

	// This should return an error because the list ID does not exist.
	if err == nil {
		t.Errorf("Contains() did not return an error when it should have")
	}

	if contains != false {
		t.Errorf("Contains() returned incorrect value: got %v, expected %v", contains, false)
	}

}

func TestFileLinkedList_RemoveOne(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-file")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	defer func(tempFile *os.File) {
		err := tempFile.Close()
		if err != nil {
			t.Fatalf("Failed to close temporary file: %v", err)
		}
	}(tmpFile)
	defer os.Remove(tmpFile.Name())

	// Create a new FileLinkedList and perform the Add operation.
	list, err := NewFileLinkedList[number.Number](tmpFile.Name(), false)
	if err != nil {
		t.Errorf("NewFileLinkedList() failed: %v", err)
	}

	testItem := number.NewNumber(24)
	err = list.Add(0, testItem)
	if err != nil {
		return
	}

	contains, err := list.Contains(0, testItem)
	if err != nil {
		return
	}

	if !contains {
		t.Fatalf("Contains() returned incorrect value: got %v, expected %v", contains, true)
	}

	err = list.Remove(0, testItem)
	if err != nil {
		t.Fatalf("Remove() failed: %v", err)
	}

	contains, err = list.Contains(0, testItem)
	if err == nil {
		t.Errorf("Contains() did not return an error when it should have")
	}

	if contains {
		t.Fatalf("Contains() returned incorrect value: got %v, expected %v", contains, false)
	}

}
