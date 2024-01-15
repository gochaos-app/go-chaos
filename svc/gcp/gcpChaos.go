package gcp

import (
	"context"
	"log"

	resourcemanager "cloud.google.com/go/resourcemanager/apiv3"
	"github.com/gochaos-app/go-chaos/config"
	resourcemanagerpb "google.golang.org/genproto/googleapis/cloud/resourcemanager/v3"

	"google.golang.org/api/iterator"
)

type gcpfn func(string, string, string, string, int, bool)

func GoogleChaos(cfg config.JobConfig, dry bool) {
	//search for project, if it doesn't exists return and print and error
	ctx := context.Background()
	c, err := resourcemanager.NewProjectsClient(ctx)
	if err != nil {
		log.Println(err)
	}
	defer c.Close()
	rqst := &resourcemanagerpb.SearchProjectsRequest{}
	it := c.SearchProjects(ctx, rqst)

	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Println(err)
		}

		if resp.DisplayName == cfg.Project {
			log.Println("GCP Project exists")
			break
		} else {
			log.Println("GCP Project does not exists")
			return
		}
	}

	gcpMap := map[string]gcpfn{
		"vm": instanceFn,
	}

	if _, servExists := gcpMap[cfg.Service]; servExists {
		gcpMap[cfg.Service](cfg.Project, cfg.Region, cfg.Chaos.Tag, cfg.Chaos.Chaos, cfg.Chaos.Count, dry)
	} else {
		log.Println("Service not found")
		return
	}

}
