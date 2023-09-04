package serialization

import "testing"

type SampleStruct1 struct {
	Name  string
	Age   int
	Score float64
}

type SampleStruct2 struct {
	ID    int
	Title string
}

func TestGenerateStructStructureHash_ValidStruct1(t *testing.T) {
	data := SampleStruct1{}
	hash := GenerateStructStructureHash(data)
	expectedHash := uint32(3586140521)
	if hash != expectedHash {
		t.Errorf("TestGenerateStructStructureHash_ValidStruct1 failed: Expected %d, got %d", expectedHash, hash)
	}
}

func TestGenerateStructStructureHash_ValidStruct2(t *testing.T) {
	data := SampleStruct2{}
	hash := GenerateStructStructureHash(data)
	expectedHash := uint32(2542233370)
	if hash != expectedHash {
		t.Errorf("TestGenerateStructStructureHash_ValidStruct2 failed: Expected %d, got %d", expectedHash, hash)
	}
}

func TestGenerateStructStructureHash_InvalidInput(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("TestGenerateStructStructureHash_InvalidInput failed: Expected panic, but function did not panic")
		}
	}()
	data := 42 // Passing an integer, which is not a struct
	GenerateStructStructureHash(data)
}

func TestGenerateStructStructureHash_NoCollision(t *testing.T) {
	// Create instances of the two structs
	data1 := SampleStruct1{}
	data2 := SampleStruct2{}

	// Generate hash values for both structs
	hash1 := GenerateStructStructureHash(data1)
	hash2 := GenerateStructStructureHash(data2)

	// Ensure that the hash values are not equal
	if hash1 == hash2 {
		t.Errorf("TestGenerateStructStructureHash_NoCollision failed: Hash values for different structs are equal")
	}
}
