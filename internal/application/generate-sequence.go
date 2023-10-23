package application

import (
	"github.com/lukasgolson/FileArray/serialization"
	"math/rand"
)

func GenerateRandomSequence(seed int64, length int64, randSource *rand.Rand) []byte {
	randSource.Seed(seed)

	byteArray := make([]byte, length)

	for i := serialization.Length(0); i < serialization.Length(length); i++ {
		val := randSource.Intn(101) // Generate a value between 0 and 100
		byteArray[i] = byte(val)
	}

	return byteArray
}
