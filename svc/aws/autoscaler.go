package aws

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
)

type chaosAutoscalerfn func([]string, int, string, *autoscaling.Client) error

func autoscalerFn(sess aws.Config, tag string, chaos string, number int, dry bool) error {
	svc := autoscaling.NewFromConfig(sess)

	var key, value string

	parts := strings.Split(tag, ":")
	key = "tag:" + parts[0]
	value = parts[1]

	result, err := svc.DescribeAutoScalingGroups(context.TODO(), &autoscaling.DescribeAutoScalingGroupsInput{
		Filters: []types.Filter{
			{
				Name:   aws.String(key),
				Values: []string{*aws.String(value)},
			},
		},
	})
	if err != nil {
		log.Println("Got an error retrieving information about EC2 autoscaling:", err)
		return err
	}

	var autoscalingList []string
	for _, r := range result.AutoScalingGroups {
		autoscalingList = append(autoscalingList, *r.AutoScalingGroupName)
	}

	if len(autoscalingList) == 0 {
		err := errors.New("Chaos not permitted: autoscaling groups not found with associated tag")
		return err
	}
	if dry == true {
		log.Println("Dry mode")
		log.Println("Will apply chaos on ", number, "of Autoscaling", autoscalingList)
		return nil
	}

	autoscalingMap := map[string]chaosAutoscalerfn{
		"terminate": terminateAutoScalingFn,
		"update":    updateAutoscalingFn,
		"desired":   addtoAutoscalingFn,
	}
	if _, servExists := autoscalingMap[chaos]; servExists {
		autoscalingMap[chaos](autoscalingList, number, tag, svc)
	} else {
		err := errors.New("chaos not found")
		return err
	}

	return nil
}

func updateAutoscalingFn(list []string, num int, tag string, session *autoscaling.Client) error {
	num32 := int32(num)
	if len(list) > 1 {
		err := errors.New("Found more than one autoscaling groups with provided tags")
		return err
	}
	autoscalingName := list[0]
	input := &autoscaling.UpdateAutoScalingGroupInput{
		AutoScalingGroupName: aws.String(autoscalingName),
		MaxSize:              aws.Int32(num32),
		DesiredCapacity:      aws.Int32(num32),
		MinSize:              aws.Int32(0),
	}

	log.Println("Updating autoscaling group:", autoscalingName)
	_, err := session.UpdateAutoScalingGroup(context.TODO(), input)
	if err != nil {
		log.Println("Error updating autoscaling group:", err)
	}

	return nil
}

func terminateAutoScalingFn(list []string, num int, tag string, session *autoscaling.Client) error {
	if num <= 0 {
		err := errors.New("Error, when terminate AWS autoscaler, count parameter should be a positive integer")
		return err
	}
	if num > len(list) {
		err := errors.New("Chaos not permitted: autoscaling groups found with associated tag is smaller than the count")
		return err
	}

	list = list[:num]
	for _, name := range list {
		_, err := session.DeleteAutoScalingGroup(context.TODO(), &autoscaling.DeleteAutoScalingGroupInput{
			AutoScalingGroupName: aws.String(name),
			ForceDelete:          aws.Bool(true),
		})
		log.Println("Terminating autoscaling group:", name)
		if err != nil {
			log.Println("Error terminating autoscaling group:", err)
		}
	}

	return nil
}

func addtoAutoscalingFn(list []string, num int, tag string, session *autoscaling.Client) error {
	if len(list) > 1 {
		err := errors.New("Found more than one autoscaling groups with the associated tags")
		return err
	}
	num32 := int32(num)
	autoscalingName := list[0]
	input := &autoscaling.SetDesiredCapacityInput{
		AutoScalingGroupName: aws.String(autoscalingName),
		DesiredCapacity:      aws.Int32(num32),
	}

	_, err := session.SetDesiredCapacity(context.TODO(), input)
	if err != nil {
		return err
	}

	return nil
}
