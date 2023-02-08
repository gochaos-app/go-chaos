package aws

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gochaos-app/go-chaos/ops"
)

type chaosS3fn func([]string, int, *s3.Client) error

func s3Fn(sess aws.Config, tag string, chaos string, number int, dry bool) error {
	svc := s3.NewFromConfig(sess)

	if number <= 0 {
		err := errors.New("Will not destroy any S3")
		return err
	}

	s3List, err := svc.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		return err
	}
	var bucketList []string
	for _, bucket := range s3List.Buckets {
		bucketList = append(bucketList, *bucket.Name)
	}

	//Separate tag into key, value components
	parts := strings.Split(tag, ":")
	var fixList []string
	switch parts[0] {
	case "SUFFIX":
		for _, name := range bucketList {
			if strings.HasSuffix(name, parts[1]) {
				fixList = append(fixList, name)
			}
		}
	case "PREFIX":
		for _, name := range bucketList {
			if strings.HasPrefix(name, parts[1]) {
				fixList = append(fixList, name)
			}
		}
	default:
		err := errors.New("Chaos not permitted: Please use PREFIX or SUFFIX as tag key")
		return err
	}
	if len(fixList) == 0 {
		err := errors.New("Chaos not permitted: Couldn't find a buckets with the characteristics specified in the config file")
		return err
	}
	if dry == true {
		log.Println("Dry mode")
		log.Println("Will apply chaos on ", number, "of s3 list", bucketList)
		return nil
	}

	s3Map := map[string]chaosS3fn{
		"terminate":      terminateS3Fn,
		"delete_content": deletectnS3Fn,
	}
	if _, servExists := s3Map[chaos]; servExists {
		s3Map[chaos](ops.RandomArray(fixList), number, svc)
	} else {
		err := errors.New("Chaos not found")
		return err
	}

	return nil
}

func terminateS3Fn(list []string, number int, session *s3.Client) error {
	buckets := len(list)
	if buckets < number {
		err := errors.New("Out of dimension array, trying to delete more S3 buckets than found")
		return err
	}
	bucketsArray := list[0:number]
	for _, S3Bucket := range bucketsArray {
		input_list := &s3.ListObjectsV2Input{
			Bucket: aws.String(S3Bucket),
		}
		results, err := session.ListObjectsV2(context.TODO(), input_list)
		for _, item := range results.Contents {
			input_delete := &s3.DeleteObjectInput{
				Bucket: aws.String(S3Bucket),
				Key:    aws.String(*item.Key),
			}
			_, err := session.DeleteObject(context.TODO(), input_delete)
			if err != nil {
				return err
			}
		}
		log.Println("Deleting bucket:", S3Bucket)
		input := &s3.DeleteBucketInput{
			Bucket: aws.String(S3Bucket),
		}
		_, s3error := session.DeleteBucket(context.TODO(), input)
		if s3error != nil {
			return err
		}
	}

	return nil
}

func deletectnS3Fn(list []string, number int, session *s3.Client) error {
	for _, S3Bucket := range list {
		input_list := &s3.ListObjectsV2Input{
			Bucket: aws.String(S3Bucket),
		}
		results, err := session.ListObjectsV2(context.TODO(), input_list)
		if err != nil {
			log.Panicln("Error:", err)
		}
		objects := len(results.Contents)
		if objects == 0 {
			err := errors.New("No objects in to delete bucket")
			return err
		} else if objects < number {
			err := errors.New("Out of dimension array, trying to delete more objects than objects found in the bucket")
			return err
		}
		objectsArray := results.Contents[0:number]

		for _, item := range objectsArray {
			log.Println("Deleting object:", item)
			input_delete := &s3.DeleteObjectInput{
				Bucket: aws.String(S3Bucket),
				Key:    aws.String(*item.Key),
			}
			_, err := session.DeleteObject(context.TODO(), input_delete)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
