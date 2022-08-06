package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

func s3Fn(sess aws.Config, tag string, chaos string, number int) {
	svc := lambda.NewFromConfig(sess)

	fmt.Println(svc, tag, chaos, number)

}
