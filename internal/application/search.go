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
)

type SeedDistance struct {
	Seed     int32
	Distance uint32
}

func Search(inputFile string, delimiter string) error {

	trees, err := findNodesFiles("data")
	if err != nil {
		return err
	}

	parsedValues, err := readFileAndParse(inputFile, delimiter, 0, 100)
	if err != nil {
		return err
	}

	var seedDistances []SeedDistance

	for _, s := range trees {
		// Define the window size and stride
		stride := 1

		// Loop through the byte slice
		for i := 0; i <= len(parsedValues)-32; i += stride {
			sequence := parsedValues[i : i+32]

			found, seed, distance := FindClosestInTree(s, [32]byte(sequence))

			if found {
				seedDistances = append(seedDistances, SeedDistance{Seed: seed, Distance: distance})
			}

			//println(FormatByteArrayAsNumbers(sequence, 32))

		}
	}

	// Print the results
	for _, sd := range seedDistances {
		fmt.Printf("Seed: %d, Match Distance: %d\n", sd.Seed, sd.Distance)
	}

	return nil

}

func findNodesFiles(dataPath string) ([]string, error) {
	var nodesFiles []string
	stack := []string{dataPath}

	for len(stack) > 0 {
		dir := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// List the contents of the current directory
		entries, err := os.ReadDir(dir)
		if err != nil {
			return nil, err
		}

		// Process entries
		for _, entry := range entries {
			entryPath := filepath.Join(dir, entry.Name())

			if entry.IsDir() {
				// If it's a directory, add it to the stack for further processing
				stack = append(stack, entryPath)
			} else if strings.HasSuffix(entry.Name(), ".nodes.bin") {
				// If it's a file ending with ".nodes.bin", add its path without extension to the list
				fileWithoutExt := strings.TrimSuffix(entryPath, ".nodes.bin")
				nodesFiles = append(nodesFiles, fileWithoutExt)
			}
		}
	}

	return nodesFiles, nil
}

func FindClosestInTree(treePath string, sequence [32]byte) (found bool, seed int32, distance uint32) {
	bkTree, err := tree.New(treePath)

	if err != nil {
		panic(err)
	}

	result := bkTree.FindClosestElement(sequence, 16)

	if result.Distance != math.MaxUint32 {
		return true, result.Seed, result.Distance
	} else {
		return false, -1, math.MaxUint32
	}
}

// readFileAndParse reads a text file, parses numbers using the specified delimiter, and returns them as a byte slice.
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

// splitByDelimiter splits a string using the specified delimiter and returns a slice of substrings.
func splitByDelimiter(s string, delimiter string) []string {
	return strings.Split(s, delimiter)
}
