package application

import (
	"awesomeProject/internal/tree"
	"bufio"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

type SeedDistance struct {
	Seed     int32
	Distance uint32
}

// Search performs a search operation on a set of trees.
func Search(inputFile string, delimiter string, numCPU int) error {
	// Save the current GOMAXPROCS value and set it to numCPU.
	previousMaxProcs := runtime.GOMAXPROCS(numCPU)
	defer runtime.GOMAXPROCS(previousMaxProcs) // Restore the original GOMAXPROCS value when done.

	trees, err := findNodesFiles("data")
	if err != nil {
		return err
	}

	parsedValues, err := readFileAndParse(inputFile, delimiter, 0, 100)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	resultsChan := make(chan []SeedDistance)

	// Perform search in parallel for each tree.
	for _, treePath := range trees {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			results := searchInTree(parsedValues, path)
			resultsChan <- results
		}(treePath)
	}

	// Start a goroutine to collect and print the results.
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Collect and print the results from the channel.
	seedDistances := make([]SeedDistance, 0)
	for result := range resultsChan {
		seedDistances = append(seedDistances, result...)
	}

	// Print the results
	for _, sd := range seedDistances {
		fmt.Printf("Seed: %d, Match Distance: %d\n", sd.Seed, sd.Distance)
	}

	return nil
}

// searchInTree searches for sequences in a tree and returns the results.
func searchInTree(parsedValues []byte, treePath string) []SeedDistance {
	const sequenceLength = 32
	stride := 1
	seedDistances := make([]SeedDistance, 0)

	bkTree, err := tree.New(treePath)
	if err != nil {
		panic(err)
	}

	for i := 0; i <= len(parsedValues)-sequenceLength; i += stride {
		sequence := parsedValues[i : i+sequenceLength]
		found, seed, distance := FindClosestInTree(bkTree, sequence)

		if found {
			seedDistances = append(seedDistances, SeedDistance{Seed: seed, Distance: distance})
		}
	}

	return seedDistances
}

// findNodesFiles finds and returns a list of node files in a directory tree.
func findNodesFiles(dataPath string) ([]string, error) {
	var nodesFiles []string
	stack := []string{dataPath}

	for len(stack) > 0 {
		dir := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		entries, err := os.ReadDir(dir)
		if err != nil {
			return nil, err
		}

		for _, entry := range entries {
			entryPath := filepath.Join(dir, entry.Name())

			if entry.IsDir() {
				stack = append(stack, entryPath)
			} else if strings.HasSuffix(entry.Name(), ".nodes.bin") {
				fileWithoutExt := strings.TrimSuffix(entryPath, ".nodes.bin")
				nodesFiles = append(nodesFiles, fileWithoutExt)
			}
		}
	}

	return nodesFiles, nil
}

// FindClosestInTree finds the closest element in the tree for a given sequence.
func FindClosestInTree(bkTree *tree.Tree, sequence []byte) (found bool, seed int32, distance uint32) {
	result := bkTree.FindClosestElement([32]byte(sequence), 16)

	if result.Distance != math.MaxUint32 {
		return true, result.Seed, result.Distance
	}
	return false, -1, math.MaxUint32
}

// readFileAndParse reads a file, parses it, and returns the parsed values as a byte array.
func readFileAndParse(filename string, delimiter string, minNumber, maxNumber int) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var byteArray []byte
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		nums := splitByDelimiter(line, delimiter)

		for _, numStr := range nums {
			num, err := strconv.Atoi(numStr)
			if err != nil {
				return nil, err
			}
			if num < minNumber || num > maxNumber {
				return nil, fmt.Errorf("number %d is out of range (%d-%d)", num, minNumber, maxNumber)
			}
			byteArray = append(byteArray, byte(num))
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return byteArray, nil
}

// splitByDelimiter splits a string by a delimiter and returns the substrings.
func splitByDelimiter(s string, delimiter string) []string {
	return strings.Split(s, delimiter)
}
