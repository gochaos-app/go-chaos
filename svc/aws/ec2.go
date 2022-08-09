package aws

import (
	"context"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type chaosEC2fn func([]string, int, *ec2.Client)

func ec2Fn(sess aws.Config, tag string, chaos string, number int) {
	svc := ec2.NewFromConfig(sess)
	var key, value string
	if number == 0 {
		log.Println("Will not destroy any EC2")
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
		log.Panicln("Got an error retrieving information about your Amazon EC2 instances:", err)
		return
	}
	var EC2instances []string
	for _, r := range result.Reservations {
		for _, i := range r.Instances {
			EC2instances = append(EC2instances, *i.InstanceId)
		}
	}

	if len(EC2instances) >= number {
		log.Println("EC2 Chaos permitted...")
	} else {
		log.Println("Chaos not permitted", len(EC2instances), "instances found with", key, value, "Number of instances to destroy is:", number)
		return
	}

	ec2Map := map[string]chaosEC2fn{
		"terminate": terminateEC2Fn,
		"stop":      stopEC2Fn,
		"reboot":    rebootEC2Fn,
	}

	ec2Map[chaos](EC2instances, number, svc)

}

func rebootEC2Fn(instances []string, num2Kill int, session *ec2.Client) {
	input := &ec2.RebootInstancesInput{
		InstanceIds: []string{},
	}

	instances = instances[:num2Kill]
	for _, id := range instances {
		log.Println("Rebooting instances:", id)
		input.InstanceIds = append(input.InstanceIds, *aws.String(id))
	}

	_, err := session.RebootInstances(context.TODO(), input)
	if err != nil {
		log.Panicln("Error deleting instances:", err)
	}
}

func stopEC2Fn(instances []string, num2Kill int, session *ec2.Client) {

	input := &ec2.StopInstancesInput{
		InstanceIds: []string{},
	}

	instances = instances[:num2Kill]
	for _, id := range instances {
		log.Println("Stopping instances:", id)
		input.InstanceIds = append(input.InstanceIds, *aws.String(id))
	}

	_, err := session.StopInstances(context.TODO(), input)
	if err != nil {
		log.Println("Error deleting instances:", err)
	}
}

func terminateEC2Fn(instances []string, num2Kill int, session *ec2.Client) {

	input := &ec2.TerminateInstancesInput{
		InstanceIds: []string{},
	}

	instances = instances[:num2Kill]
	for _, id := range instances {
		log.Println("Terminating instance:", id)
		input.InstanceIds = append(input.InstanceIds, *aws.String(id))
	}

	_, err := session.TerminateInstances(context.TODO(), input)
	if err != nil {
		log.Println("Error deleting instances:", err)
	}
}
