package aws

import (
	"context"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/mental12345/chaosctl/ops"
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
		log.Println("Chaos not permitted: Please use PREFIX or SUFFIX as tag key")
		return
	}
	if len(fixList) == 0 {
		log.Println("Chaos not permitted: Couldn't find a buckets with the characteristics specified in the config file")
		return
	}

	s3Map := map[string]chaosS3fn{
		"terminate":      terminateS3Fn,
		"delete_content": deletectnS3Fn,
	}
	s3Map[chaos](ops.Random(fixList), number, svc)

}

func terminateS3Fn(list []string, number int, session *s3.Client) {
	buckets := len(list)
	if buckets < number {
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
		log.Println("Deleting bucket:", S3Bucket)
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
			log.Println("Deleting object:", item)
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
