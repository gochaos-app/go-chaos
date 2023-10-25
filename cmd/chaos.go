package cmd

import (
	"log"
	"strconv"
	"sync"

	"github.com/gochaos-app/go-chaos/notifications"
	"github.com/gochaos-app/go-chaos/svc/aws"
	"github.com/gochaos-app/go-chaos/svc/do"
	"github.com/gochaos-app/go-chaos/svc/gcp"
	"github.com/gochaos-app/go-chaos/svc/k8s"
	"github.com/gochaos-app/go-chaos/svc/scripts"
)

func switchService(job JobConfig, dry bool) {
	switch job.Cloud {
	case "aws":
		aws.AmazonChaos(
			job.Region,
			job.Service,
			job.Chaos.Tag,
			job.Chaos.Chaos,
			job.Chaos.Count,
			dry)
	case "do":
		do.DigitalOceanChaos(
			job.Region,
			job.Service,
			job.Chaos.Tag,
			job.Chaos.Chaos,
			job.Chaos.Count,
			dry)
	case "kubernetes":
		k8s.KubernetesChaos(
			job.Namespace,
			job.Service,
			job.Chaos.Tag,
			job.Chaos.Chaos,
			job.Chaos.Count,
			dry)
	case "gcp":
		gcp.GoogleChaos(
			job.Region,
			job.Project,
			job.Service,
			job.Chaos.Tag,
			job.Chaos.Chaos,
			job.Chaos.Count,
			dry)
	case "script":
		scripts.ExecuteScript(
			job.Region,
			job.Project,
			job.Service,
			job.Namespace,
			job.Chaos.Tag,
			job.Chaos.Chaos,
			job.Chaos.Count,
			dry)

	case "":
		log.Println("I don't know what to do")
	default:
		log.Println("I don't understand the service to execute chaos on")
	}
}

func selectFunction(cfg *GenConfig) []JobConfig {
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

func ExecuteChaos(cfg *GenConfig, dryFlag bool) error {
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
		/*
			if cfg.Script != nil {
				scripts.ExecuteScript(cfg.Script.Source, cfg.Script.Executor)
			}
		*/
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
