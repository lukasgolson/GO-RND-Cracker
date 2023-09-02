package algorithms

import "testing"

func TestMeyersDifferenceAlgorithmEmptySlices(t *testing.T) {
	score := MeyersDifferenceAlgorithm([]byte{}, []byte{})
	if score != 0 {
		t.Errorf("Expected score 0, but got %d", score)
	}
}

func TestMeyersDifferenceAlgorithmOneEmptySlice(t *testing.T) {
	score := MeyersDifferenceAlgorithm([]byte{}, []byte{1, 2, 3})
	if score != 3 {
		t.Errorf("Expected score 3, but got %d", score)
	}
}

func TestMeyersDifferenceAlgorithmEqualSlices(t *testing.T) {
	score := MeyersDifferenceAlgorithm([]byte{1, 2, 3}, []byte{1, 2, 3})
	if score != 0 {
		t.Errorf("Expected score 0, but got %d", score)
	}
}

func TestMeyersDifferenceAlgorithmOnePrefixOfOther(t *testing.T) {
	score := MeyersDifferenceAlgorithm([]byte{1, 2, 3}, []byte{1, 2, 3, 4, 5})
	if score != 2 {
		t.Errorf("Expected score 2, but got %d", score)
	}
}

func TestMeyersDifferenceAlgorithmOneSuffixOfOther(t *testing.T) {
	score := MeyersDifferenceAlgorithm([]byte{1, 2, 3, 4, 5}, []byte{3, 4, 5})
	if score != 2 {
		t.Errorf("Expected score 2, but got %d", score)
	}
}

func TestMeyersDifferenceAlgorithmDifferentSlices(t *testing.T) {
	score := MeyersDifferenceAlgorithm([]byte{1, 2, 3}, []byte{4, 5, 6})
	if score != 3 {
		t.Errorf("Expected score 3, but got %d", score)
	}
}
