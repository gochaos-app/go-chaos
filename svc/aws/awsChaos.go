package aws

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

type fn func(string, string)

func AmazonChaos(region string, service string, tags []string, chaos string, number int) {
	//AWS session for each region in the config
	_, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		log.Println("error", err)
	}

	awsMap := map[string]fn{
		"lambda": lambda,
		"ec2":    ec2,
		"s3":     s3,
	}

	switch service {
	case "lambda":
		awsMap[service]("hola desde lambda", "service")
	case "ec2":
		awsMap[service]("hola ec2 que show", "hello")
	case "s3":
		awsMap[service]("Hola desde s3", "hello")
	}
}
