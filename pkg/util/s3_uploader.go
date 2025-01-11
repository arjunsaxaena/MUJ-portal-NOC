package util

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func createS3Client() (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-north-1"))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}

	svc := s3.NewFromConfig(cfg)
	return svc, nil
}

func UploadFileToS3(bucketName string, file io.Reader, key string) error {
	svc, err := createS3Client()
	if err != nil {
		return err
	}

	_, err = svc.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String("application/pdf"),
	})

	if err != nil {
		return fmt.Errorf("unable to upload to %v, %v", bucketName, err)
	}

	fmt.Printf("Successfully uploaded to %v\n", bucketName)
	return nil
}
