package fileArray

import (
	"awesomeProject/internal/number"
	"awesomeProject/internal/serialization"
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

	fA, err := NewFileArray[number.Number](tmpFile.Name(), false)
	if err != nil {
		t.Fatalf("NewFileArray returned an error: %v", err)
	}
	defer fA.Close() // Clean up

	if fA == nil {
		t.Fatalf("FileArray instance is nil")
	}
}

func TestOpenAndInitializeFile(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer func(tempFile *os.File) {
		err := tempFile.Close()
		if err != nil {
			t.Fatalf("Failed to close temporary file: %v", err)
		}
	}(tempFile)
	defer os.Remove(tempFile.Name())

	// Call the function being tested
	file, err := openAndInitializeFile[number.Number](tempFile.Name())

	// Check for errors
	if err != nil {
		t.Fatalf("openAndInitializeFile returned an error: %v", err)
	}

	// Verify the file exists
	if _, err := os.Stat(tempFile.Name()); os.IsNotExist(err) {
		t.Fatalf("The file should exist but doesn't.")
	}

	// Clean up: close the file
	file.Close()
}

func TestOpenMmap(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	// Create a sample data to write into the file
	data := []byte("Hello, World!")

	// Write the data to the temporary file
	_, err = tempFile.Write(data)
	if err != nil {
		t.Fatalf("Failed to write data to the file: %v", err)
	}

	// Open the file using mmap
	memoryMap, err := openReadWriteMmap(tempFile)

	// Check for errors
	if err != nil {
		t.Fatalf("openReadWriteMmap returned an error: %v", err)
	}
	defer memoryMap.Unmap() // Ensure we unmap the memory

	// Check if the mapped data matches the original data
	if string(memoryMap) != string(data) {
		t.Fatalf("Mapped data does not match original data.")
	}
}

func generateTestCases(numTestCases int) []struct {
	value serialization.Offset
} {
	testCases := make([]struct {
		value serialization.Offset
	}, numTestCases)

	for i := 0; i <= numTestCases-2; i++ {
		testCases[i].value = serialization.Offset(rand.Intn((i * 10) + 1)) // Generate values algorithmically
	}

	testCases[numTestCases-1].value = serialization.MaxOffset() // Generate values algorithmically

	return testCases
}

func TestFileArray_Count(t *testing.T) {

	tmpFile, err := os.CreateTemp("", "test-file")
	defer os.Remove(tmpFile.Name()) // Clean up

	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	fA, err := NewFileArray[number.Number](tmpFile.Name(), false)
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
	defer os.Remove(tmpFile.Name()) // Clean up

	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	fA, err := NewFileArray[number.Number](tmpFile.Name(), false)

	err = fA.Close()
	if err != nil {
		t.Fatalf("Failed to close file array: %v", err)
	}
}

func TestFileArrayCountEmptyMemoryMap(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-file")
	defer os.Remove(tmpFile.Name()) // Clean up

	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	fileArray, err := NewFileArray[number.Number](tmpFile.Name(), false)

	if err != nil {
		t.Fatalf("Failed to create file array: %v", err)
	}

	count := fileArray.Count()
	if count != 0 {
		t.Fatalf("Count() returned %d, expected 0", count)
	}
}

func TestFileArraySetCount(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-file")

	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	fileArray, err := NewFileArray[number.Number](tmpFile.Name(), false)

	if err != nil {
		t.Fatalf("Failed to create file array: %v", err)
	}

	// Test the setCount() method
	expectedCount := serialization.Offset(42)
	fileArray.setCount(expectedCount)
	count := fileArray.Count()
	if count != expectedCount {
		t.Fatalf("setCount() did not set the count correctly. Got %d, expected %d", count, expectedCount)
	}
}

func TestFileArrayIncrementCount(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-file")
	defer os.Remove(tmpFile.Name()) // Clean up

	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	fileArray, err := NewFileArray[number.Number](tmpFile.Name(), false)

	if err != nil {
		t.Fatalf("Failed to create file array: %v", err)
	}

	// Test the incrementCount() method
	expectedCount := serialization.Offset(42)
	fileArray.setCount(expectedCount)
	fileArray.incrementCount()
	count := fileArray.Count()
	if count != expectedCount+1 {
		t.Fatalf("incrementCount() did not increment the count correctly. Got %d, expected %d", count, expectedCount+1)
	}
}

func TestFileArray_FileRetrieve(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-file")
	defer os.Remove(tmpFile.Name()) // Clean up

	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	fA, err := NewFileArray[number.Number](tmpFile.Name(), false)

	for i := 0; i <= 24; i++ {
		_, err = fA.Append(number.NewNumber(int64(i)))
	}

	if err != nil {
		return
	}

	err = fA.Close()
	if err != nil {
		t.Fatalf("Failed to close file array: %v", err)
	}

	fA2, err := NewFileArray[number.Number](tmpFile.Name(), false)

	if err != nil {
		t.Fatalf("Failed to create file array: %v", err)
	}

	for i := 0; i <= 24; i++ {
		num, err := fA2.GetItemFromIndex(serialization.Offset(i))
		if err != nil {
			t.Fatalf("Failed to get item from index: %v", err)
		}
		if num.Value != int64(i) {
			t.Fatalf("Expected %d, got %d", i, num.Value)
		}
	}
}
