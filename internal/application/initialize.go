package application

import (
	"awesomeProject/internal/serialization"
	"awesomeProject/internal/tree"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"
)

func generateRandomSequence(seed int64, length int64, randSource *rand.Rand) []byte {
	randSource.Seed(seed)

	byteArray := make([]byte, length)

	for i := serialization.Length(0); i < serialization.Length(length); i++ {
		val := randSource.Intn(101) // Generate a value between 0 and 100
		byteArray[i] = byte(val)
	}

	return byteArray
}

func processPartition(lo, hi int64, randSource *rand.Rand) error {

	// Create a new tree, name it after the partition. This will be the file name.
	// Place it in a subdirectory of the current directory.

	subdir := fmt.Sprintf("partition-%d", lo)
	os.MkdirAll(subdir, 0755)

	var tree, _ = tree.New(subdir + "/graph")

	var counter = 0

	for i := lo + 1; i < hi; i++ {
		seed := i
		counter++

		// Generate the sequence based on the random source
		var sequence = generateRandomSequence(seed, 32, randSource)

		tree.Add([32]byte(sequence), int32(seed))

	}

	return nil
}

func Initialize() {

	numPartitions := runtime.NumCPU() - 1

	seedCount := int64(1<<31 - 1)

	partitionSize := seedCount / int64(numPartitions)

	var wg sync.WaitGroup

	startTime := time.Now()

	for p := int64(0); p < int64(numPartitions); p++ {
		lo := partitionSize * p
		hi := partitionSize * (p + 1)
		fmt.Printf("Partition %d (%d, %d)\n", p, lo, hi)

		wg.Add(1)

		go func(lo, hi int64) {
			defer wg.Done()

			randSource := rand.New(rand.NewSource(0))
			if err := processPartition(lo, hi, randSource); err != nil {
				log.Printf("Error processing partition: %v\n", err)
			}
		}(lo, hi)
	}

	wg.Wait()

	endTime := time.Now()
	fmt.Printf("Finished seed generation in %s.\n", endTime.Sub(startTime))

}
