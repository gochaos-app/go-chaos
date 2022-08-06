package aws

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type chaosEC2fn func([]string, int, *ec2.Client)

func ec2Fn(sess aws.Config, tag string, chaos string, number int) {
	svc := ec2.NewFromConfig(sess)
	var key, value string
	if number == 0 {
		log.Println("Not going to destroy any EC2")
		return
	}
	parts := strings.Split(tag, ":")
	key = "tag:" + parts[0]
	value = parts[1]
	result, err := svc.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("instance-state-name"),
				Values: []string{*aws.String("running")},
			},
			{
				Name:   aws.String(key),
				Values: []string{*aws.String(value)},
			},
		},
	})
	if err != nil {
		fmt.Println("Got an error retrieving information about your Amazon EC2 instances:")
		fmt.Println(err)
		return
	}
	var EC2instances []string
	for _, r := range result.Reservations {
		for _, i := range r.Instances {
			EC2instances = append(EC2instances, *i.InstanceId)
		}
	}

	if len(EC2instances) >= number {
		log.Println("EC2 Chaos initiated")
	} else {
		log.Println("chaos not permitted", len(EC2instances), "instances found", key, ":", value, "review config")
		return
	}

	ec2Map := map[string]chaosEC2fn{
		"terminate": terminateFn,
		"stop":      stopFn,
		"reboot":    rebootFn,
	}
	if chaos == "random" {
		rand.Seed(time.Now().UnixNano())

		randomSelect := rand.Intn(3) // terminate stop reboot 0 1 2
		switch randomSelect {
		case 0:
			log.Println("Terminating instances... ")
			ec2Map["terminate"](EC2instances, number, svc)
		case 1:
			log.Println("Stopping instances... ")
			ec2Map["stop"](EC2instances, number, svc)
		case 2:
			log.Println("Rebooting instances... ")
			ec2Map["reboot"](EC2instances, number, svc)
		}
	} else {
		ec2Map[chaos](EC2instances, number, svc)
	}

}

func rebootFn(instances []string, num2Kill int, session *ec2.Client) {

	input := &ec2.RebootInstancesInput{
		InstanceIds: []string{},
	}

	if num2Kill != 0 {
		instances = instances[:num2Kill]
	}
	for _, id := range instances {
		log.Println(id)
		input.InstanceIds = append(input.InstanceIds, *aws.String(id))
	}

	_, err := session.RebootInstances(context.TODO(), input)
	if err != nil {
		log.Println("Error deleting instances:", err)
	}
}

func stopFn(instances []string, num2Kill int, session *ec2.Client) {

	input := &ec2.StopInstancesInput{
		InstanceIds: []string{},
	}

	if num2Kill != 0 {
		instances = instances[:num2Kill]
	}
	for _, id := range instances {
		log.Println(id)
		input.InstanceIds = append(input.InstanceIds, *aws.String(id))
	}

	_, err := session.StopInstances(context.TODO(), input)
	if err != nil {
		log.Println("Error deleting instances:", err)
	}
}

func terminateFn(instances []string, num2Kill int, session *ec2.Client) {

	input := &ec2.TerminateInstancesInput{
		InstanceIds: []string{},
	}

	if num2Kill != 0 {
		instances = instances[:num2Kill]
	}
	for _, id := range instances {
		log.Println(id)
		input.InstanceIds = append(input.InstanceIds, *aws.String(id))
	}

	_, err := session.TerminateInstances(context.TODO(), input)
	if err != nil {
		log.Println("Error deleting instances:", err)
	}
}
