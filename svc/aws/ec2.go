package aws

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/gochaos-app/go-chaos/ops"
)

type chaosEC2fn func([]string, int, *ec2.Client) error

func ec2Fn(sess aws.Config, tag string, chaos string, number int, dry bool) error {
	svc := ec2.NewFromConfig(sess)
	var key, value string

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
		return err
	}

	var EC2instances []string
	for _, r := range result.Reservations {
		for _, i := range r.Instances {
			EC2instances = append(EC2instances, *i.InstanceId)
		}
	}

	// FIXME: Implement a better way to handle below logic
	if number <= 0 {
		err := errors.New("Will not destroy any EC2")
		return err
	}
	if len(EC2instances) == 0 {
		err := errors.New("Could not find any instance with provided tag")
		return err
	}
	if dry == true {
		log.Println("Dry mode")
		log.Println("Will apply chaos on ", number, "of EC2 list", EC2instances)
		return nil
	}
	if len(EC2instances) >= number {
		log.Println("EC2 Chaos permitted...")
	} else {
		err := errors.New("Chaos not permitted: Number of instances to destroy is greater than instances found with provided tag")
		return err
	}

	ec2Map := map[string]chaosEC2fn{
		"terminate": terminateEC2Fn,
		"stop":      stopEC2Fn,
		"reboot":    rebootEC2Fn,
	}
	if _, servExists := ec2Map[chaos]; servExists {
		ec2Map[chaos](ops.RandomArray(EC2instances), number, svc)
	} else {
		err := errors.New("Chaos not found")
		return err
	}

	return nil
}

func rebootEC2Fn(instances []string, num2Kill int, session *ec2.Client) error {
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
		return err
	}

	return nil
}

func stopEC2Fn(instances []string, num2Kill int, session *ec2.Client) error {
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
		return err
	}

	return nil
}

func terminateEC2Fn(instances []string, num2Kill int, session *ec2.Client) error {
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
		return err
	}

	return nil
}
