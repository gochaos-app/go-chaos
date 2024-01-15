package cmd

import (
	"log"
	"strconv"
	"sync"

	"github.com/gochaos-app/go-chaos/config"
	"github.com/gochaos-app/go-chaos/notifications"
	"github.com/gochaos-app/go-chaos/svc/aws"
	"github.com/gochaos-app/go-chaos/svc/do"
	"github.com/gochaos-app/go-chaos/svc/gcp"
	"github.com/gochaos-app/go-chaos/svc/k8s"
	"github.com/gochaos-app/go-chaos/svc/scripts"
)

type cloudfn func(config.JobConfig, bool)

func switchService(job config.JobConfig, dry bool) {
	cloudMap := map[string]cloudfn{
		"aws":        aws.AmazonChaos,
		"do":         do.DigitalOceanChaos,
		"kubernetes": k8s.KubernetesChaos,
		"gcp":        gcp.GoogleChaos,
		"script":     scripts.ExecuteScript,
	}
	if _, servExists := cloudMap[job.Cloud]; servExists {
		cloudMap[job.Cloud](job, dry)
	} else {
		log.Println("Service not found")
		return
	}
}

func selectFunction(cfg *config.GenConfig) []config.JobConfig {
	switch cfg.Function {
	case "random":
		return randomJobs(cfg.Job)
	case "sequential":
		return cfg.Job
	case "":
		return cfg.Job
	}
	return nil
}

func ExecuteChaos(cfg *config.GenConfig, dryFlag bool) error {
	var wg sync.WaitGroup
	done := make(chan struct{})
	if cfg.Hypothesis != nil {

		workers, _ := strconv.Atoi(cfg.Hypothesis.Pings)
		wg.Add(workers)
		for i := 0; i < workers; i++ {
			go Ping(cfg.Hypothesis.Url, cfg.Hypothesis.Report, &wg)
		}

	}

	selectFunction(cfg)

	for i := 0; i < len(cfg.Job); i++ {
		switchService(cfg.Job[i], dryFlag)
	}
	// After executing chaos, if there's a script will be executed

	if !dryFlag {
		for i := 0; i < len(cfg.Notifications); i++ {
			switch cfg.Notifications[i].Type {
			case "gmail":
				notifications.GmailNotification(cfg.Notifications[i].To, cfg.Notifications[i].Body, cfg.Notifications[i].From)
			case "slack":
				notifications.SlackNotification(cfg.Notifications[i].To, cfg.Notifications[i].Body)
			default:
				log.Println("I don't understand the notification")
			}
		}
	} else {
		log.Println("Dry run, no chaos executed, ignoring notifications")
	}

	close(done)
	wg.Wait()
	return nil
}
