package main

import (
	"context"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Creates a new folder in bucket for user
func createUserHolder(keyName string) (err error) {
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
		//Delimiter: aws.String("/"),
	})
	if err != nil {
		return nil, err
	}
	keyNameStruct := make([]string, 0)
	for _, item := range resp.Contents {
		keyName := DerefString(item.Key)
		if keyName != folderName {
			nwKeyName := strings.TrimPrefix(keyName, folderName)
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

func createFolder(folderName, userFolder string) error {
	_, err := svc.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(userFolder + folderName),
	})
	if err != nil {
		return err
	}
	return nil
}

func downloadFile(keyName string) (io.Reader, error) {
	object, err := svc.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(keyName),
	})
	if err != nil {
		return nil, err
	}
	return object.Body, nil
}

//us-east-1
//AWS access portal URL
//https://d-9067b9f60f.awsapps.com/start
