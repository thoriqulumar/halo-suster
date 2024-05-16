package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"helo-suster/config"
	"io"
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ImageService interface {
	UploadImage(file *multipart.FileHeader) <-chan string
}

type imageService struct {
	cfg *config.Config
	log *zap.Logger
}

func NewImageService(cfg *config.Config, log *zap.Logger) ImageService {
	return &imageService{
		cfg: cfg,
		log: log,
	}
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

		url, err := uploadToS3(fileBytes, fileName, s.cfg)
		if err != nil {
			fileURLChan <- ""
			return
		}

		fileURLChan <- url
	}()

	return fileURLChan
}

func uploadToS3(fileBytes []byte, filename string, cfg *config.Config) (string, error) {
	// Initialize AWS session
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(cfg.S3Region),
		Credentials: credentials.NewStaticCredentials(cfg.S3Id, cfg.S3Secret, ""),
	})
	if err != nil {
		return "", errors.New("failed to create AWS session")
	}

	// Create S3 service client
	svc := s3.New(sess)

	// Specify bucket name and object key
	bucketName := cfg.S3Bucket
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
