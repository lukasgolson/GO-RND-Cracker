package fileArray

import (
	"math"
	"math/rand"
	"os"
	"testing"
)

func TestNewFileArray(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "test-file")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up

	fA, err := NewFileArray(tmpFile.Name())
	if err != nil {
		t.Fatalf("NewFileArray returned an error: %v", err)
	}
	defer fA.Close() // Clean up

	if fA == nil {
		t.Fatalf("FileArray instance is nil")
	}
}

func generateTestCases(numTestCases int) []struct {
	value uint64
} {
	testCases := make([]struct {
		value uint64
	}, numTestCases)

	for i := 0; i <= numTestCases-2; i++ {
		testCases[i].value = uint64(rand.Intn((i * 10) + 1)) // Generate values algorithmically
	}

	testCases[numTestCases-1].value = math.MaxUint64 // Generate values algorithmically

	return testCases
}

func TestFileArray_Count(t *testing.T) {

	tmpFile, err := os.CreateTemp("", "test-file")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up

	fA, err := NewFileArray(tmpFile.Name())
	if err != nil {
		t.Fatalf("NewFileArray returned an error: %v", err)
	}
	defer fA.Close() // Clean up

	testCases := generateTestCases(200)

	// Test getting length for each test case
	for _, tc := range testCases {

		fA.setCount(tc.value)

		length := fA.Count()
		if length != tc.value {
			t.Fatalf("Expected length: %d, Got: %d", tc.value, length)
		}

		fA.setCount(tc.value)
		if updatedLength := fA.Count(); updatedLength != tc.value {
			t.Fatalf("Expected length: %d, Got: %d", tc.value, updatedLength)
		}
	}
}

func TestFileArray_Close(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-file")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	fA, err := NewFileArray(tmpFile.Name())

	err = fA.Close()
	if err != nil {
		t.Fatalf("Failed to close file array: %v", err)
	}
}
