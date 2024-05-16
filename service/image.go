package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

type ImageService interface {
	UploadImage(file *multipart.FileHeader) <-chan string
}

type imageService struct{}

func NewImageService() ImageService {
	return &imageService{}
}

func (s *imageService) UploadImage(file *multipart.FileHeader) <-chan string {
	fileURLChan := make(chan string, 1)

	go func() {
		defer close(fileURLChan)

		src, err := file.Open()
		if err != nil {
			fileURLChan <- ""
			return
		}
		defer src.Close()

		fileBytes, err := io.ReadAll(src)
		if err != nil {
			fileURLChan <- ""
			return
		}

		uuid := uuid.New().String()
		fileName := uuid + ".jpeg"

		url, err := uploadToS3(fileBytes, fileName)
		if err != nil {
			fileURLChan <- ""
			return
		}

		fileURLChan <- url
	}()

	return fileURLChan
}

func uploadToS3(fileBytes []byte, filename string) (string, error) {
	// Initialize AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("your-region"), // Specify your AWS region
	})
	if err != nil {
		return "", fmt.Errorf("failed to create AWS session: %w", err)
	}

	// Create S3 service client
	svc := s3.New(sess)

	// Specify bucket name and object key
	bucketName := "your-bucket-name"
	objectKey := filename

	// Upload file to S3
	_, err = svc.PutObjectWithContext(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   aws.ReadSeekCloser(bytes.NewReader(fileBytes)),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	// Generate S3 object URL
	objectURL := fmt.Sprintf("https://awss3/%s/%s", bucketName, objectKey)

	return objectURL, nil
}
