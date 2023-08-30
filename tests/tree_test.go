package tests

import (
	"awesomeProject/internal/bktree"
	"testing"
)

func TestBkTreeSearch(t *testing.T) {
	// Create a BK tree with a root word and seed
	rootWord := []byte("hello")
	seed := int32(12345)
	bk, _ := bktree.NewBkTree("test")

	bk.Add(rootWord, seed)

	// Add words to the BK tree
	wordsToAdd := [][]byte{
		[]byte("helo"),
		[]byte("hell"),
		[]byte("help"),
		[]byte("hall"),
		[]byte("hole"),
	}

	for _, word := range wordsToAdd {
		bk.Add(word, seed)
	}

	// Define test cases
	testCases := []struct {
		queryWord []byte
		tolerance int
		expected  int
	}{
		{[]byte("hello"), 0, 1}, // Expected result: 1 word with distance 0
		{[]byte("helo"), 1, 4},  // Expected result: 3 words with distance <= 1
		{[]byte("hello"), 2, 5}, // Expected result: 5 words with distance <= 2
		{[]byte("help"), 1, 3},  // Expected result: 3 words with distance <= 1
		{[]byte("world"), 1, 0}, // Expected result: 0 words with distance <= 1
	}

	// Run test cases
	for _, testCase := range testCases {
		results := bk.Search(testCase.queryWord, testCase.tolerance)
		if len(results) != testCase.expected {
			t.Errorf("For query '%s' and tolerance %d, expected %d results but got %d", testCase.queryWord, testCase.tolerance, testCase.expected, len(results))
		}
	}
}
