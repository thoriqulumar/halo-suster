package service

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Service interface {
	UploadImage(fileName string, file *os.File) (string, error)
}

type s3Uploader struct {
	sess *session.Session
}

func NewS3Service(region string) (S3Service, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, err
	}
	return &s3Uploader{sess: sess}, nil
}

func (s *s3Uploader) UploadImage(fileName string, file *os.File) (string, error) {
	s3Svc := s3.New(s.sess)

	bucketName := "your-bucket-name"
	objectKey := fileName

	_, err := s3Svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   file,
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		return "", err
	}

	// Construct the image URL
	imageURL := fmt.Sprintf("https://awss3/%s/%s", bucketName, objectKey)
	return imageURL, nil
}
