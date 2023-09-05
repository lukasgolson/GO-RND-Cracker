package tree

import (
	"awesomeProject/internal/util"
	"os"
	"testing"
)

func TestAddToBKTree_EmptyTree(t *testing.T) {
	// Create a test instance of the BKTree
	tmpFile, _ := os.CreateTemp("", "test-file")
	tree, err := NewTree(tmpFile.Name())
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
	tmpFile, _ := os.CreateTemp("", "test-file")
	tree, err := NewTree(tmpFile.Name())

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

	tmpFile, _ := os.CreateTemp("", "test-file")
	tree, err := NewTree(tmpFile.Name())

	if err != nil {
		t.Fatalf("Error creating tree: %v", err)
	}

	defer cleanup(tree)

	wordStrings := util.GetWordList()

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

func TestFuzzyMatch(t *testing.T) {
	tree, err := NewTree("Test")
	if err != nil {
		t.Fatalf("Error creating tree: %v", err)
	}

	defer cleanup(tree)

	for i, wordString := range util.GetWordList() {
		word := make([]byte, NodeWordSize)

		copy(word[:], wordString)

		seed := i
		err = tree.AddToBKTree(0, [32]byte(word), int32(seed))
		if err != nil {
			t.Errorf("Expected no error when adding to a non-empty tree, but got: %v", err)
		}

	}

	word := make([]byte, NodeWordSize)
	copy(word[:], "cat")

	foundNode, distance := tree.FindClosestElement(0, [32]byte(word), 1)

	if distance != 0 {
		t.Errorf("Expected distance to be 0, but got: %v.", distance)
	}

	word = make([]byte, NodeWordSize)
	copy(word[:], "cats")

	_, distance = tree.FindClosestElement(0, [32]byte(word), 10)

	if distance != 1 {
		t.Errorf("Expected distance to be 1, but got: %v. Tested word %v, found word %v", distance, string(word[:]), string(foundNode.Word[:]))
	}
}

func cleanup(tree *Tree) {
	file1, file2 := tree.getFileNames()

	tree.Close()

	os.Remove(file1)

	os.Remove(file2)
}
