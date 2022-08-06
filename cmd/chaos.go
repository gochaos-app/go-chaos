package cmd

import (
	"fmt"

	"github.com/mental12345/chaosCLI/svc/aws"
	"github.com/mental12345/chaosCLI/svc/k8s"
)

func ExecuteChaos(cfg *GenConfig) error {
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
			fmt.Println("this will create chaos in k8s")
			k8s.KubernetesChaos()
		case "script":
			fmt.Println("This will execute a custom script")
		case "":
			fmt.Println("I dont know what to do")
		default:
			fmt.Println("I dont understand the service to execute chaos on")
		}
	}
	return nil
}
