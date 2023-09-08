package tree

import (
	"os"
	"testing"
)

func TestNewTree(t *testing.T) {
	tree, err := NewTree("test")
	if err != nil {
		t.Fatalf("Error creating tree: %v", err)
	}

	if tree == nil {
		t.Fatalf("Tree is nil")
	}

	file1, file2, file3 := tree.getFileNames()

	if err := tree.Close(); err != nil {
		t.Fatalf("Error closing tree: %v", err)
	}

	os.Remove(file1)
	os.Remove(file2)
	os.Remove(file3)
}

func TestTreeIsEmpty(t *testing.T) {
	tree, err := NewTree("test")

	if err != nil {
		t.Fatalf("Error creating tree: %v", err)
	}

	// Use defer to schedule the file removal functions
	defer func() {
		file1, file2, file3 := tree.getFileNames()

		if err := tree.Close(); err != nil {
			t.Fatalf("Error closing tree: %v", err)
		}

		if err := os.Remove(file1); err != nil {
			t.Fatalf("Error deleting file1: %v", err)
		}

		if err := os.Remove(file2); err != nil {
			t.Fatalf("Error deleting file2: %v", err)
		}

		if err := os.Remove(file3); err != nil {
			t.Fatalf("Error deleting file2: %v", err)
		}
	}()

	if !tree.isEmpty() {
		t.Fatalf("Tree is not empty")
	}

}

func TestTreeAddNode(t *testing.T) {
	tree, err := NewTree("test")

	if err != nil {
		t.Fatalf("Error creating tree: %v", err)
	}

	// Use defer to schedule the file removal functions
	defer func() {
		file1, file2, file3 := tree.getFileNames()

		if err := tree.Close(); err != nil {
			t.Fatalf("Error closing tree: %v", err)
		}

		if err := os.Remove(file1); err != nil {
			t.Fatalf("Error deleting file1: %v", err)
		}

		if err := os.Remove(file2); err != nil {
			t.Fatalf("Error deleting file2: %v", err)
		}

		if err := os.Remove(file3); err != nil {
			t.Fatalf("Error deleting file2: %v", err)
		}
	}()

	tree.addNode([NodeWordSize]byte{}, 0)
}

func TestTreeAddEdge(t *testing.T) {
	tree, err := NewTree("test")

	if err != nil {
		t.Fatalf("Error creating tree: %v", err)
	}

	// Use defer to schedule the file removal functions
	defer func() {
		file1, file2, file3 := tree.getFileNames()

		if err := tree.Close(); err != nil {
			t.Fatalf("Error closing tree: %v", err)
		}

		if err := os.Remove(file1); err != nil {
			t.Fatalf("Error deleting file1: %v", err)
		}

		if err := os.Remove(file2); err != nil {
			t.Fatalf("Error deleting file2: %v", err)
		}

		if err := os.Remove(file3); err != nil {
			t.Fatalf("Error deleting file2: %v", err)
		}
	}()

	err = tree.AddEdge(0, 1, 0)
	if err != nil {
		t.Fatalf("Error adding edge: %v", err)
	}

}
