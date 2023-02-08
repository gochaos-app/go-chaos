package aws

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/gochaos-app/go-chaos/ops"
)

type chaosLambdafn func([]string, int, *lambda.Client) error

func lambdaFn(sess aws.Config, tag string, chaos string, number int, dry bool) error {
	svc := lambda.NewFromConfig(sess)
	if number <= 0 {
		err := errors.New("Will not destroy any Lambda")
		return err
	}
	parts := strings.Split(tag, ":")
	key := parts[0]
	value := parts[1]
	result, err := svc.ListFunctions(context.TODO(), &lambda.ListFunctionsInput{})
	if err != nil {
		log.Panicln(err)
	}
	var arnFunctions []string
	for _, f := range result.Functions {
		list, err := svc.ListTags(context.TODO(), &lambda.ListTagsInput{
			Resource: f.FunctionArn,
		})
		if err != nil {
			log.Panicln("An error has occurred: ", err)
			return err
		}
		if v, found := list.Tags[key]; found {
			if v == value {
				arnFunctions = append(arnFunctions, *f.FunctionArn)
			} else {
				err := errors.New("Chaos not permitted: Couldn't find lambda functions with the characteristics specified in the config file")
				return err
			}
		}
	}

	// FIXME: Implement a better way to handle this logic
	if len(arnFunctions) == 0 {
		err := errors.New("Chaos not permitted: Couldn't find lambda functions with the characteristics specified in the config file")
		return err
	}
	if dry == true {
		log.Println("Dry mode")
		log.Println("Will apply chaos on ", number, "of lambda list", arnFunctions)
		return nil
	}
	if number > len(arnFunctions) {
		err := errors.New("Chaos not permitted: Out of dimension array, trying to delete more than the functions found")
		return err
	}

	lambdaMap := map[string]chaosLambdafn{
		"terminate": terminateLambdaFn,
		"stop":      stopLambdaFn,
	}
	if _, servExists := lambdaMap[chaos]; servExists {
		lambdaMap[chaos](ops.RandomArray(arnFunctions), number, svc)
	} else {
		err := errors.New("Chaos not found")
		return err
	}

	return nil
}

func terminateLambdaFn(list []string, number int, session *lambda.Client) error {
	list = list[:number]
	for _, lambdaARN := range list {
		input := lambda.DeleteFunctionInput{
			FunctionName: aws.String(lambdaARN),
		}
		log.Println("Terminating Lambda function:", lambdaARN)
		_, err := session.DeleteFunction(context.TODO(), &input)
		if err != nil {
			return err
		}
	}

	return nil
}

func stopLambdaFn(list []string, number int, session *lambda.Client) error {
	list = list[:number]
	for _, lambdaARN := range list {
		input := lambda.PutFunctionConcurrencyInput{
			FunctionName:                 aws.String(lambdaARN),
			ReservedConcurrentExecutions: aws.Int32(0),
		}
		log.Println("Stopping Lambda function:", lambdaARN)
		_, err := session.PutFunctionConcurrency(context.TODO(), &input)
		if err != nil {
			return err
		}
	}

	return nil
}
