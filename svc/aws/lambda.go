package aws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/gochaos-app/go-chaos/ops"
)

type chaosLambdafn func([]string, int, *lambda.Client)

func lambdaFn(sess aws.Config, key string, value string, chaos string, number int, dry bool) {
	svc := lambda.NewFromConfig(sess)
	if number <= 0 {
		log.Println("Will not destroy any Lambda")
		return
	}

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
			return
		}
		if v, found := list.Tags[key]; found {
			if v == value {
				arnFunctions = append(arnFunctions, *f.FunctionArn)
			} else {
				log.Println("Chaos not permitted: Couldn't find lambda functions with the characteristics specified in the config file")
				return
			}
		}
	}

	if len(arnFunctions) == 0 {
		log.Println("Chaos not permitted: Couldn't find lambda functions with the characteristics specified in the config file")
		return
	}
	if dry == true {
		log.Println("Dry mode")
		log.Println("Will apply chaos on ", number, "of lambda list", arnFunctions)
		return
	}
	if number > len(arnFunctions) {
		log.Println("Chaos not permitted: Out of dimension array, trying to delete", number, "functions.", len(arnFunctions), "functions found")
		return
	}

	lambdaMap := map[string]chaosLambdafn{
		"terminate": terminateLambdaFn,
		"stop":      stopLambdaFn,
	}
	if _, servExists := lambdaMap[chaos]; servExists {
		lambdaMap[chaos](ops.RandomArray(arnFunctions), number, svc)
	} else {
		log.Println("Chaos not found")
		return
	}

}

func terminateLambdaFn(list []string, number int, session *lambda.Client) {
	list = list[:number]
	for _, lambdaARN := range list {
		input := lambda.DeleteFunctionInput{
			FunctionName: aws.String(lambdaARN),
		}
		log.Println("Terminating Lambda function:", lambdaARN)
		_, err := session.DeleteFunction(context.TODO(), &input)
		if err != nil {
			log.Panicln("Error:", err)
		}
	}
}

func stopLambdaFn(list []string, number int, session *lambda.Client) {
	list = list[:number]
	for _, lambdaARN := range list {
		input := lambda.PutFunctionConcurrencyInput{
			FunctionName:                 aws.String(lambdaARN),
			ReservedConcurrentExecutions: aws.Int32(0),
		}
		log.Println("Stopping Lambda function:", lambdaARN)
		_, err := session.PutFunctionConcurrency(context.TODO(), &input)
		if err != nil {
			log.Panicln("Error:", err)
		}
	}
}
