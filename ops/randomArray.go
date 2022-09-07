package ops

import (
	"math/rand"
	"time"
)

// Random take an array of string or int arguments, using generics :), and returns the same array but randomized
func Random[T string | int](array []T) []T {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(array), func(i, j int) {
		array[i], array[j] = array[j], array[i]
	})
	return array
}
