package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

type awsfn func(*session.Session, []string, string, int)

func AmazonChaos(region string, service string, tags []string, chaos string, number int) {
	//AWS session for each region in the config
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	fmt.Printf("%T", sess)
	if err != nil {
		log.Println("error", err)
	}

	awsMap := map[string]awsfn{
		"lambda": lambdaFn,
		"ec2":    ec2Fn,
		"s3":     s3Fn,
	}
	awsMap[service](sess, tags, chaos, number)
	/*switch service {
	case "lambda":
		//string, string, *session.Session, []string, string, int)
		awsMap[service](sess, tags, chaos, number)
	case "ec2":
		awsMap[service](sess, tags, chaos, number)
	case "s3":
		awsMap[service](sess, tags, chaos, number)
	}*/
}
