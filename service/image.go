package service

import (
	"helo-suster/repo"
	"io"

	"github.com/google/uuid"
)

type ImageService interface {
	SaveImage(file io.Reader) (string, error)
}

type imageService struct {
	repo repo.ImageRepository
}

func NewImageService(repo repo.ImageRepository) ImageService {
	return &imageService{
		repo: repo,
	}
}

func (s *imageService) SaveImage(file io.Reader) (string, error) {
	uuid := uuid.New().String()
	fileName := uuid + ".jpeg"

	err := s.repo.SaveImage(file, fileName)
	if err != nil {
		return "", err
	}

	return fileName, err
}
