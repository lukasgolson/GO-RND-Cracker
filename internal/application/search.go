package application

import (
	"awesomeProject/internal/tree"
	"bufio"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type SeedDistance struct {
	Seed     int32
	Distance uint32
}

func Search(inputFile string, delimiter string, dataDirectories []string, concurrentTrees int, stride int, prefetch bool) error {
	parsedValues, err := readFileAndParse(inputFile, delimiter, 0, 100)
	if err != nil {
		return err
	}

	treeFiles := make([]string, 0)

	for _, directory := range dataDirectories {
		files, err := findNodesFiles(directory)
		if err != nil {
			return err
		}

		treeFiles = append(treeFiles, files...)
	}

	resultsChan := make(chan SeedDistance)

	var wg sync.WaitGroup

	// Process trees in separate goroutine
	go func() {

		counter := 0
		for _, treePath := range treeFiles {
			bkTree, err := tree.NewOrLoad(treePath, false)
			if err != nil {
				fmt.Printf("Error loading tree for path %s: %v\n", treePath, err)
				return
			} else {
				fmt.Printf("Initialized tree for path %s\n", treePath)
			}

			// Prefetch the tree if needed
			if prefetch {
				fmt.Printf("Prefetching tree for path %s\n", treePath)
				bkTree.Prefetch()
				fmt.Printf("Prefetching tree for path %s done\n", treePath)
			}

			for i := len(parsedValues) - 32; i >= 0; i -= stride {
				wg.Add(1)
				sequence := parsedValues[i : i+32]

				go func(seq []byte) {
					defer wg.Done() // Ensure we always decrement the wait group
					
					found, result := searchInTree(seq, bkTree)
					if found {
						resultsChan <- result
					}
				}(sequence)
			}

			counter++

			if counter >= concurrentTrees {
				// Wait for all goroutines to finish for these trees
				wg.Wait()
				counter = 0
			}

		}

		// Close the results channel after all trees are processed
		close(resultsChan)
	}()

	// Process results from the channel
	for result := range resultsChan {
		// Process the result as needed
		fmt.Println("Result:", result.Seed, " ", result.Distance)
	}

	return nil
}

func searchInTree(sequence []byte, bkTree *tree.Tree) (bool, SeedDistance) {

	found, seed, distance := FindClosestInTree(bkTree, sequence)
	if found {
		return found, SeedDistance{Seed: seed, Distance: distance}
	} else {
		return false, SeedDistance{Seed: -1, Distance: math.MaxUint32}
	}
}

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

func FindClosestInTree(bkTree *tree.Tree, sequence []byte) (found bool, seed int32, distance uint32) {
	result := bkTree.FindClosestElement([32]byte(sequence), 16)

	if result.Distance != math.MaxUint32 {
		return true, result.Seed, result.Distance
	}
	return false, -1, math.MaxUint32
}

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

func splitByDelimiter(s string, delimiter string) []string {
	return strings.Split(s, delimiter)
}
