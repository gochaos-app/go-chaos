package cmd

import (
	"fmt"

	"github.com/mental12345/go-chaos/notifications"
	"github.com/mental12345/go-chaos/svc/aws"
	"github.com/mental12345/go-chaos/svc/do"
	"github.com/mental12345/go-chaos/svc/gcp"
	"github.com/mental12345/go-chaos/svc/k8s"
	"github.com/mental12345/go-chaos/svc/scripts"
)

func switchService(job JobConfig) {
	switch job.Cloud {
	case "aws":
		aws.AmazonChaos(
			job.Region,
			job.Service,
			job.Chaos.Tag,
			job.Chaos.Chaos,
			job.Chaos.Count)
	case "do":
		do.DigitalOceanChaos(
			job.Region,
			job.Service,
			job.Chaos.Tag,
			job.Chaos.Chaos,
			job.Chaos.Count)
	case "kubernetes":
		k8s.KubernetesChaos(
			job.Namespace,
			job.Service,
			job.Chaos.Tag,
			job.Chaos.Chaos,
			job.Chaos.Count)
	case "gcp":
		gcp.GoogleChaos(
			job.Region,
			job.Project,
			job.Service,
			job.Chaos.Tag,
			job.Chaos.Chaos,
			job.Chaos.Count)
	case "":
		fmt.Println("I dont know what to do")
	default:
		fmt.Println("I dont understand the service to execute chaos on")
	}
}

func ExecuteChaos(cfg *GenConfig) error {
	for i := 0; i < len(cfg.Job); i++ {
		switchService(cfg.Job[i])
	}
	// After executing chaos, if there's a script will be executed
	if cfg.Script != nil {
		scripts.ExecuteScript(cfg.Script.Source, cfg.Script.Executor)
	}
	for i := 0; i < len(cfg.Notifications); i++ {
		switch cfg.Notifications[i].Type {
		case "gmail":
			notifications.GmailNotification(cfg.Notifications[i].ToEmail, cfg.Notifications[i].Body, cfg.Notifications[i].FromEmail)
		case "":
			fmt.Println("I dont know what to do")
		default:
			fmt.Println("I dont understand the notification")
		}
	}
	return nil
}
