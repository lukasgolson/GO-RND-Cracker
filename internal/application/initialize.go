package application

import (
	"awesomeProject/internal/serialization"
	"awesomeProject/internal/tree"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
)

func processPartition(lo, hi, fileCount int64, graphPath string, randSource *rand.Rand) error {
	numberOfSeeds := hi - lo

	if fileCount > numberOfSeeds {
		fileCount = numberOfSeeds
	}

	seedsPerFile := numberOfSeeds / fileCount

	for fileIndex := int64(0); fileIndex < fileCount; fileIndex++ {
		var bkTree, err = tree.NewOrLoad(graphPath+fmt.Sprintf("-%d", fileIndex), false)
		if err != nil {
			return err
		}
		startSeed := lo + (fileIndex * seedsPerFile)
		endSeed := startSeed + seedsPerFile
		loadedSeedPosition := serialization.Length(lo) + bkTree.Length()

		if bkTree.Length() > 0 {
			fmt.Println("Loading existing tree... Start seed", startSeed, "end seed:", endSeed, "previous end seed:", loadedSeedPosition)

			if loadedSeedPosition > serialization.Length(startSeed) && loadedSeedPosition < serialization.Length(endSeed) {

				fmt.Println("Previous tree is in range. Everything is fine.")

				const seedOverlap = 5

				newStartSeed := int64(loadedSeedPosition) - seedOverlap

				if newStartSeed < startSeed {
					newStartSeed = startSeed
				}

				startSeed = newStartSeed

			} else if loadedSeedPosition < serialization.Length(startSeed) {
				return fmt.Errorf("previous tree ends before our start seed")
			} else if loadedSeedPosition > serialization.Length(endSeed) {
				return fmt.Errorf("previous tree ends after our end seed")
			}
		} else {
			fmt.Println("No existing tree found. Creating new tree with name", graphPath, "... Start seed", startSeed, "end seed:", endSeed)
			err := bkTree.PreExpand(serialization.Length(seedsPerFile))
			if err != nil {
				return err
			}
		}

		for seed := startSeed; seed < endSeed; seed++ {
			sequence := GenerateRandomSequence(seed, 32, randSource)

			err := bkTree.Add([32]byte(sequence), int32(seed))
			if err != nil {
				return err
			}
		}

		err = bkTree.ShrinkWrap()
		if err != nil {
			return err
		}
	}

	return nil
}

func Initialize(coreCount int, fileCount int, seedCount int64, dataDirectories []string) error {

	if fileCount < 1 {
		return fmt.Errorf("file count must be at least 1")
	}

	if fileCount < coreCount {
		return fmt.Errorf("file count must be greater than or equal to the core count")
	}

	if fileCount%coreCount != 0 {
		return fmt.Errorf("file count must be divisible by the core count")
	}

	if seedCount < 1 {
		return fmt.Errorf("seed count must be at least 1")
	}

	if int64(fileCount) > seedCount {
		return fmt.Errorf("file count must be less than or equal to the seed count")
	}

	for _, dir := range dataDirectories {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	var wg sync.WaitGroup

	partitionSize := seedCount / int64(coreCount)

	filesPerPartition := int64(fileCount) / int64(coreCount)

	for p := int64(0); p < int64(coreCount); p++ {
		lo := partitionSize * p
		hi := partitionSize * (p + 1)

		dirIndex := int(p) % len(dataDirectories)
		dataDirectory := dataDirectories[dirIndex]

		wg.Add(1)

		go func(lo, hi int64, partitionID int64) {
			defer wg.Done()
			randSource := rand.New(rand.NewSource(0))
			dir := fmt.Sprintf("%s/graph-%d", dataDirectory, partitionID)

			if err := processPartition(lo, hi, filesPerPartition, dir, randSource); err != nil {
				log.Printf("Error processing partition: %v\n", err)
			}
		}(lo, hi, p)
	}

	wg.Wait()

	return nil
}
