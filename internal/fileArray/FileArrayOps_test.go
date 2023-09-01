package fileArray

import (
	"awesomeProject/internal/serialization"
	"os"
	"testing"
)

func TestAppendItemSpace(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-file")
	defer os.Remove(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	fA, err := NewFileArray(tmpFile.Name())

	num := serialization.NewNumber(42)

	err = fA.expandMemoryMapSize(int64(num.SerializedSize()))
	if err != nil {
		t.Fatalf("Failed to expand memory map size: %v", err)
	}

	err = AppendItem(fA, &num)
	if err != nil {
		t.Fatalf("Failed to append item when space is available: %v", err)
	}

}

func TestAppendItemNoSpace(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-file")
	defer os.Remove(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	fA, err := NewFileArray(tmpFile.Name())

	num := serialization.NewNumber(42)

	err = AppendItem(fA, &num)

	if r := recover(); r != nil {
		t.Fatalf("Failed to append item when space is not available: %v", err)
	}

	if err != nil {
		t.Fatalf("Failed to append item when space is not available: %v", err)
	}
}

func TestSetItemAtIndex(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-file")
	defer os.Remove(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	fA, err := NewFileArray(tmpFile.Name())

	num := serialization.NewNumber(42)

	err = SetItemAtIndex(fA, &num, 0)
	if err != nil {
		t.Fatalf("Failed to set item at index: %v", err)
	}
}

func TestSetItemAtIndexWithIndexOutOfBounds(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-file")
	defer os.Remove(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	fA, err := NewFileArray(tmpFile.Name())

	num := serialization.NewNumber(42)

	err = SetItemAtIndex(fA, &num, 1)
	if err == nil {
		t.Fatalf("SetItemAtIndex did not fail with index out of bounds")
	}
}

func TestGetItemFromIndex(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-file")
	defer os.Remove(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	fA, err := NewFileArray(tmpFile.Name())

	num := serialization.NewNumber(42)

	err = SetItemAtIndex(fA, &num, 0)
	if err != nil {
		t.Fatalf("Failed to set item at index: %v", err)
	}

	retrievedNumber, err := GetItemFromIndex[serialization.Number](fA, 0)

	if err != nil {
		t.Fatalf("Failed to get item from index: %v", err)
	}

	if retrievedNumber != num {
		t.Fatalf("Retrieved number does not match the original number: %v, expected %v", retrievedNumber, num)
	}

	// Now you can use numberItem as a serialization.Number

}
