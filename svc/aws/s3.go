package aws

import (
	"context"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type chaosS3fn func([]string, int, *s3.Client)

func s3Fn(sess aws.Config, tag string, chaos string, number int) {
	svc := s3.NewFromConfig(sess)

	if number == 0 {
		log.Println("Will not destroy any S3")
		return
	}

	s3List, err := svc.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		log.Println("An error has occurred: ", err)
		return
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
		log.Println("I don't understand")
	}
	if len(fixList) == 0 {
		log.Println("Couldn't find a buckets with the characteristics specified in the config file")
		return
	} else {
		log.Println("S3 Chaos permitted...")
	}

	s3Map := map[string]chaosS3fn{
		"terminate":      terminateS3Fn,
		"delete_content": deletectnS3Fn,
	}
	if chaos == "random" {
		rand.Seed(time.Now().UnixNano())
		randomSelect := rand.Intn(2)
		switch randomSelect {
		case 0:
			log.Println("Terminating S3 buckets")
			s3Map["stop"](fixList, number, svc)
		case 1:
			log.Println("Deleting content on S3 buckets")
			s3Map["delete_content"](fixList, number, svc)
		}
	} else {
		s3Map[chaos](fixList, number, svc)
	}
}

func terminateS3Fn(list []string, number int, session *s3.Client) {
	buckets := len(list)
	if buckets == 0 {
		log.Println("Can't find any s3 buckets with the specified config")
		return
	} else if buckets < number {
		log.Println("Out of dimension array, trying to delete", number, "buckets.", buckets, "S3 buckets found")
		return
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
				log.Panicln("Error:", err)
				return
			}
		}
		input := &s3.DeleteBucketInput{
			Bucket: aws.String(S3Bucket),
		}
		_, s3error := session.DeleteBucket(context.TODO(), input)
		if s3error != nil {
			log.Panicln(err)
			return
		}
	}

}

func deletectnS3Fn(list []string, number int, session *s3.Client) {
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
			log.Println("No objects in bucket", S3Bucket)
			return
		} else if objects < number {
			log.Println("Out of dimension array, trying to delete", number, "objects.", objects, "objects found in bucket:", S3Bucket)
			return
		}
		objectsArray := results.Contents[0:number]

		for _, item := range objectsArray {
			input_delete := &s3.DeleteObjectInput{
				Bucket: aws.String(S3Bucket),
				Key:    aws.String(*item.Key),
			}
			_, err := session.DeleteObject(context.TODO(), input_delete)
			if err != nil {
				log.Panicln("Error:", err)
				return
			}
		}
	}
}
