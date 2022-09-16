package main

import (
	"context"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	bucketName = "localhost-fileserver"
	region     = "eu-west-2"
)

func DerefString(s *string) string {
	if s != nil {
		return *s
	}

	return ""
}

// Creates a new aws bucket for user
func createUserBucket(keyName string) (err error) {
	_, err = svc.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(keyName + "/"),
	})
	return
}

func listFiles(folderName string) ([]string, error) {
	resp, err := svc.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(folderName),
	})
	if err != nil {
		return nil, err
	}
	keyNameStruct := make([]string, 0)
	for _, item := range resp.Contents {
		keyName := DerefString(item.Key)
		if keyName != folderName+"/" {
			nwKeyName := strings.TrimPrefix(keyName, folderName+"/")
			keyNameStruct = append(keyNameStruct, nwKeyName)
		}

	}
	return keyNameStruct, nil
}

func uploadFile(keyName string, file io.Reader) error {
	_, err := svc.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(keyName),
		Body:   file,
	})
	if err != nil {
		return err
	}

	return nil
}

func downloadFile(keyName string) io.Reader {
	object, err := svc.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(keyName),
	})
	if err != nil {
		panic("Couldn't download object")
	}
	return object.Body
}

//us-east-1
//s$ export AWS_ACCESS_KEY_ID=YOUR_AKID
//$ export AWS_SECRET_ACCESS_KEY=FNRLlaA08iqN9nJh33C1AjXTOmH+ajZ8kYBvnXi4
//AWS access portal URL
//https://d-9067b9f60f.awsapps.com/start
