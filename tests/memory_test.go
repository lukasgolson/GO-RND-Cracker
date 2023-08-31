package tests_test

import (
	"awesomeProject/internal/fileArray"
	"awesomeProject/internal/tree"
	"io/ioutil"
	"testing"
)

func TestNewFileArray(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := ioutil.TempFile("", "test-file")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	//defer os.Remove(tmpFile.Name()) // Clean up

	fA, err := fileArray.NewFileArray(tmpFile.Name())
	if err != nil {
		t.Fatalf("NewFileArray returned an error: %v", err)
	}
	defer fA.Close() // Clean up

	if fA == nil {
		t.Fatalf("FileArray instance is nil")
	}
}

func TestFileArray_GetLength(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := ioutil.TempFile("", "test-file")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	//defer os.Remove(tmpFile.Name()) // Clean up

	fA, err := fileArray.NewFileArray(tmpFile.Name())
	if err != nil {
		t.Fatalf("NewFileArray returned an error: %v", err)
	}
	defer fA.Close() // Clean up

	length := fA.Count()
	if length != 0 {
		t.Fatalf("Expected length: 0, Got: %d", length)
	}

	node := tree.NewNode(0, [32]byte{}, 0)

	fileArray.AppendItem(fA, node)

}

// Add more tests as needed...
