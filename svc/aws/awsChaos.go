package aws

import (
	"context"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	jobCfg "github.com/gochaos-app/go-chaos/config"
)

type awsfn func(aws.Config, string, string, string, int, bool)

func AmazonChaos(job jobCfg.JobConfig, dry bool) {
	//AWS session for each region in the config
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(job.Region))

	if err != nil {
		log.Fatalln("Unable to load config", err)
	}

	var key, value string
	if strings.Count(job.Chaos.Tag, ":") != 1 {
		log.Println("aws tags must be in format key:value ")
		return
	}
	parts := strings.Split(job.Chaos.Tag, ":")
	key = parts[0]
	value = parts[1]

	awsMap := map[string]awsfn{
		"ec2":             ec2Fn,
		"s3":              s3Fn,
		"lambda":          lambdaFn,
		"ec2_autoscaling": autoscalerFn,
	}

	if _, servExists := awsMap[job.Service]; servExists {
		awsMap[job.Service](cfg, key, value, job.Chaos.Chaos, job.Chaos.Count, dry)
	} else {
		log.Println("Service not found")
		return
	}

}
