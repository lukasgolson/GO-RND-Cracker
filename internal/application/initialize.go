package application

import (
	"awesomeProject/internal/tree"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
)

func processPartition(lo, hi, fileCount int64, directory string, randSource *rand.Rand) error {
	totalSeeds := hi - lo

	if fileCount > totalSeeds {
		fileCount = totalSeeds
	}

	seedsPerFile := totalSeeds / fileCount

	err := os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		return err
	}

	for fileIndex := int64(0); fileIndex < fileCount; fileIndex++ {
		startSeed := lo + (fileIndex * seedsPerFile)
		endSeed := startSeed + seedsPerFile

		treeDir := fmt.Sprintf("%s/graph-%d", directory, fileIndex)
		var bktree, err = tree.New(treeDir)
		if err != nil {
			return err
		}

		for seed := startSeed; seed < endSeed; seed++ {
			sequence := GenerateRandomSequence(seed, 32, randSource)

			err := bktree.Add([32]byte(sequence), int32(seed))
			if err != nil {
				return err
			}
		}

		err = bktree.ShrinkWrap()
		if err != nil {
			return err
		}
	}

	return nil
}

func Initialize(coreCount int, fileCount int, seedCount int64, dataDirectory string) error {

	if fileCount < 1 {
		return fmt.Errorf("file count must be at least 1")
	}

	if fileCount < coreCount {
		return fmt.Errorf("file count must be greater than or equal to the core count")
	}

	err := os.MkdirAll(dataDirectory, os.ModePerm)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	partitionSize := seedCount / int64(coreCount)

	filesPerPartition := int64(fileCount) / int64(coreCount)

	for p := int64(0); p < int64(coreCount); p++ {
		lo := partitionSize * p
		hi := partitionSize * (p + 1)
		fmt.Printf("Processing partition %d (%d, %d)\n", p, lo, hi)

		wg.Add(1)

		go func(lo, hi int64, partitionID int64) {
			defer wg.Done()
			randSource := rand.New(rand.NewSource(0))
			subdir := fmt.Sprintf("%s/partition-%d", dataDirectory, partitionID)

			if err := processPartition(lo, hi, filesPerPartition, subdir, randSource); err != nil {
				log.Printf("Error processing partition: %v\n", err)
			}
		}(lo, hi, p)
	}

	wg.Wait()

	return nil
}
