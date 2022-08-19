package cmd

import (
	"fmt"

	"github.com/mental12345/chaosCLI/svc/aws"
	"github.com/mental12345/chaosCLI/svc/k8s"
	"github.com/mental12345/chaosCLI/svc/scripts"
)

func ExecuteChaos(cfg *GenConfig) error {
	// Before executing chaos, if there's a script will be executed

	if cfg.Script != nil {
		scripts.ExecuteScript(cfg.Script.Source, cfg.Script.Executor)
	}

	for i := 0; i < len(cfg.Job); i++ {
		switch cfg.Job[i].Cloud {
		case "aws":
			aws.AmazonChaos(
				cfg.Job[i].Region,
				cfg.Job[i].Service,
				cfg.Job[i].Chaos.Tag,
				cfg.Job[i].Chaos.Chaos,
				cfg.Job[i].Chaos.Count)
		case "gcp":
			fmt.Println("his will impact on GCP")
		case "kubernetes":
			k8s.KubernetesChaos(
				cfg.Job[i].Namespace,
				cfg.Job[i].Service,
				cfg.Job[i].Chaos.Tag,
				cfg.Job[i].Chaos.Chaos,
				cfg.Job[i].Chaos.Count)
		case "":
			fmt.Println("I dont know what to do")
		default:
			fmt.Println("I dont understand the service to execute chaos on")
		}
	}
	return nil
}
