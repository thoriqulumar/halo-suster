package service

import (
	"context"
	"database/sql"

	"errors"
	cerr "helo-suster/pkg/customError"

	"helo-suster/model"
	"helo-suster/repo"
	"net/http"

	"go.uber.org/zap"
)

type MedicalService interface {
	CreateNewPatient(ctx context.Context, request model.PostPatientRequest) (patient model.Patient, err error)
}

type medicalService struct {
	repo   repo.MedicalRepo
	logger *zap.Logger
}

func NewMedicalService(r repo.MedicalRepo, logger *zap.Logger) MedicalService {
	return &medicalService{
		repo:   r,
		logger: logger,
	}
}

func (s *medicalService) CreateNewPatient(ctx context.Context, request model.PostPatientRequest) (patient model.Patient, err error) {
	_, err = s.repo.GetPatientByIdentityNumber(ctx, request.IdentityNumber)
	if !errors.Is(err, sql.ErrNoRows) {
		return model.Patient{}, cerr.New(http.StatusConflict, "phoneNumber already exists")
	}

	patient, err = s.repo.CreatePatient(ctx, request)
	if err != nil {
		return model.Patient{}, cerr.New(http.StatusInternalServerError, "Internal Server Error")
	}

	return patient, nil
}
