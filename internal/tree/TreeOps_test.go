package tree

import (
	"os"
	"testing"
)

func TestAddToBKTree_EmptyTree(t *testing.T) {
	// Create a test instance of the BKTree
	tree, err := NewTree("Test") // You need to implement a function to create a new BKTree instance

	if err != nil {
		t.Fatalf("Error creating tree: %v", err)
	}

	// Use defer to schedule the file removal functions
	defer cleanup(tree)

	word1 := [NodeWordSize]byte{'c', 'a', 't'}
	seed1 := int32(42)
	err = tree.AddToBKTree(0, word1, seed1)
	if err != nil {
		t.Errorf("Expected no error when adding to an empty tree, but got: %v", err)
	}
}

func TestAddToBKTree_AddingDuplicate(t *testing.T) {
	tree, err := NewTree("Test") // You need to implement a function to create a new BKTree instance
	if err != nil {
		t.Fatalf("Error creating tree: %v", err)
	}

	defer cleanup(tree)

	word1 := [NodeWordSize]byte{'c', 'a', 't'}
	seed1 := int32(42)
	err = tree.AddToBKTree(0, word1, seed1)
	if err != nil {
		t.Errorf("Expected no error when adding to an empty tree, but got: %v", err)
	}

	err = tree.AddToBKTree(0, word1, seed1)
	if err != nil {
		t.Errorf("Expected no error when adding to an empty tree, but got: %v", err)
	}

}

func TestAddToBKTree_AddingToNonEmptyTree(t *testing.T) {
	tree, err := NewTree("Test") // Implement a function to create a new BKTree instance
	if err != nil {
		t.Fatalf("Error creating tree: %v", err)
	}

	defer cleanup(tree)

	wordStrings := []string{
		"cat",
		"car",
		"cart",
		"carts",
		"dog",
		"dogs",
		"hello",
		"world",
		"fox",
		"bird",
		"fish",
		"pen",
		"pencil",
		"book",
		"books",
		"apple",
		"orange",
		"banana",
		"grape",
		"mango",
	}

	for i, wordString := range wordStrings {
		word := make([]byte, NodeWordSize)

		copy(word[:], wordString)

		seed := i
		err = tree.AddToBKTree(0, [32]byte(word), int32(seed))
		if err != nil {
			t.Errorf("Expected no error when adding to a non-empty tree, but got: %v", err)
		}
	}
}

func cleanup(tree *Tree) {
	file1, file2 := tree.getFileNames()

	tree.Close()

	os.Remove(file1)

	os.Remove(file2)
}
