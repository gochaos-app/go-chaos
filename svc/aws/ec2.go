package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

func ec2Fn(sess *session.Session, tags []string, chaos string, number int) {
	svc := lambda.New(sess)
	fmt.Println(svc, tags, chaos, number)
}
