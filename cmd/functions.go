package cmd

import (
	"math/rand"
	"time"
)

func randomJobs(jobs []JobConfig) []JobConfig {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(jobs), func(i, j int) {
		jobs[i], jobs[j] = jobs[j], jobs[i]
	})
	return jobs
}
