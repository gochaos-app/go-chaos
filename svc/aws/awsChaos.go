package aws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type awsfn func(aws.Config, string, string, int, bool)

func AmazonChaos(region string, service string, tag string, chaos string, number int, dry bool) {
	//AWS session for each region in the config
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))

	if err != nil {
		log.Fatalln("Unable to load config", err)
	}

	awsMap := map[string]awsfn{
		"ec2":             ec2Fn,
		"s3":              s3Fn,
		"lambda":          lambdaFn,
		"ec2_autoscaling": autoscalerFn,
	}

	if _, servExists := awsMap[service]; servExists {
		awsMap[service](cfg, tag, chaos, number, dry)
	} else {
		log.Println("Service not found")
		return
	}

}
