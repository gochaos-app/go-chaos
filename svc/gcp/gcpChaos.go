package gcp

import (
	"context"
	"log"

	resourcemanager "cloud.google.com/go/resourcemanager/apiv3"
	resourcemanagerpb "google.golang.org/genproto/googleapis/cloud/resourcemanager/v3"

	"google.golang.org/api/iterator"
)

type gcpfn func(string, string, string, string, int)

func GoogleChaos(region string, project string, service string, tag string, chaos string, number int) {
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

		if resp.DisplayName == project {
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

	if _, servExists := gcpMap[service]; servExists {
		gcpMap[service](project, region, tag, chaos, number)
	} else {
		log.Println("Service not found")
		return
	}

}
