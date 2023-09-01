package fileArray

import (
	"awesomeProject/internal/serialization"
	"os"
	"testing"
)

func TestAppendItemSpace(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-file")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	fA, err := NewFileArray(tmpFile.Name())

	num := serialization.Number(42)

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
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	fA, err := NewFileArray(tmpFile.Name())

	num := serialization.Number(42)

	err = AppendItem(fA, &num)
	if err != nil {
		t.Fatalf("Failed to append item when space is not available: %v", err)
	}
}
