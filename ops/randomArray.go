package ops

import (
	"math/rand"
	"time"
)

func Random(array []string) []string {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(array), func(i, j int) {
		array[i], array[j] = array[j], array[i]
	})
	return array
}
