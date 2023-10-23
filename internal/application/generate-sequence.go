package application

import (
	"github.com/lukasgolson/FileArray/serialization"
	"math/rand"
)

func GenerateRandomSequence(seed int64, length int64, high int, randSource *rand.Rand) []byte {
	randSource.Seed(seed)

	byteArray := make([]byte, length)

	for i := serialization.Length(0); i < serialization.Length(length); i++ {
		val := randSource.Intn(high)
		byteArray[i] = byte(val)
	}

	return byteArray
}
