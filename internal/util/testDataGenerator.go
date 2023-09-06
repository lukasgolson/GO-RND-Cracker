package util

import (
	"math/rand"
)

func GetWordList() []string {

	shuffled := []string{
		"abandon",
		"ability",
		"able",
		"about",
		"above",
		"absent",
		"absolute",
		"absolutely",
		"absorb",
		"carrot",
		"carry",
		"case",
		"cash",
		"cast",
		"cat",
		"catch",
		"category",
		"cattle",
		"cause",
		"ceiling",
		"celebrate",
		"flag",
		"flame",
		"flash",
		"flat",
		"flavor",
		"flee",
	}

	// Shuffle the copy
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	return shuffled
}
