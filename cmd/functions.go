package cmd

import (
	"math/rand"
	"time"

	"github.com/gochaos-app/go-chaos/config"
)

func randomJobs(jobs []config.JobConfig) []config.JobConfig {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	rand.Shuffle(len(jobs), func(i, j int) {
		jobs[i], jobs[j] = jobs[j], jobs[i]
	})
	return jobs
}
