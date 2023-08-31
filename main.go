package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

const (
	SeqLen = 4 * 8
)

func generateRandomSequence(seed int64, randSource *rand.Rand) []byte {
	randSource.Seed(seed)

	const bitsPerInt = 64

	const bitsPerNumber = 8
	const elementsPerInt = bitsPerInt / bitsPerNumber

	const bytesPerInt = bitsPerInt / 8

	const intCount = SeqLen / elementsPerInt

	var seq [intCount]int64

	for i := 0; i < intCount; i++ {
		var currentInt64 int64 = 0

		for j := 0; j < (elementsPerInt) && i+j < SeqLen; j++ {
			val := randSource.Intn(101) // Generate a value between 0 and 100

			currentInt64 |= int64(val&0x7F) << uint(j*bitsPerNumber)
		}

		seq[i] = currentInt64
	}

	byteArray := make([]byte, len(seq)*bytesPerInt) // 8 bytes per int64 element

	for i, val := range seq {
		binary.BigEndian.PutUint64(byteArray[i*bytesPerInt:], uint64(val))
	}

	return byteArray
}

func processPartition(lo, hi int64, randSource *rand.Rand) error {
	batchSize := 1000

	//var tree, _ = tree.NewTree("main")

	var counter = 0

	for i := lo + 1; i < hi; i++ {
		seed := i
		counter++

		fmt.Printf("", seed)

		// Generate the sequence based on the random source
		//	var sequence = generateRandomSequence(seed, randSource)

		if counter <= batchSize {
			//tree.Add(sequence, int32(seed))
		}
	}

	return nil
}

func main() {

	numPartitions := runtime.NumCPU() - 4 // Adjust the number of partitions as needed

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
