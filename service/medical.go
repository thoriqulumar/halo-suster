package service

import (
	"context"
	"database/sql"

	"errors"
	cerr "halo-suster/pkg/customErr"

	"halo-suster/model"
	"halo-suster/repo"
	"net/http"

	"go.uber.org/zap"
)

type MedicalService interface {
	CreateNewPatient(ctx context.Context, request model.PostPatientRequest) (patient model.Patient, err error)
	CreateNewMedicalRecord(ctx context.Context, request model.PostMedicalRecordRequest, createdBy string) (medicalRecord model.MedicalRecord, err error)
	GetAllPatient(ctx context.Context, params model.GetPatientParams) (listPatient []model.Patient, err error)
	GetAllMedicalRecord(ctx context.Context, params model.GetMedicalRecordParams) (listMedicalRecord []model.GetMedicalRecordData, err error)
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
		return model.Patient{}, cerr.New(http.StatusConflict, "identityNumber already exists")
	}

	patient, err = s.repo.CreatePatient(ctx, request)
	if err != nil {
		return model.Patient{}, cerr.New(http.StatusInternalServerError, "Internal Server Error")
	}

	return patient, nil
}

func (s *medicalService) CreateNewMedicalRecord(ctx context.Context, request model.PostMedicalRecordRequest, createdBy string) (medicalRecord model.MedicalRecord, err error) {
	_, err = s.repo.GetPatientByIdentityNumber(ctx, request.IdentityNumber)
	if errors.Is(err, sql.ErrNoRows) {
		return model.MedicalRecord{}, cerr.New(http.StatusNotFound, "identityNumber not found")
	}

	medicalRecord, err = s.repo.CreateMedicalRecord(ctx, request, createdBy)
	if err != nil {
		return model.MedicalRecord{}, cerr.New(http.StatusInternalServerError, "Internal Server Error")
	}

	return medicalRecord, nil
}

func (s *medicalService) GetAllPatient(ctx context.Context, params model.GetPatientParams) (listPatient []model.Patient, err error) {
	listPatient, err = s.repo.GetPatient(ctx, params)
	if err != nil {
		return
	}

	return listPatient, nil
}

func (s *medicalService) GetAllMedicalRecord(ctx context.Context, params model.GetMedicalRecordParams) (listMedicalRecord []model.GetMedicalRecordData, err error) {
	listMedicalRecord, err = s.repo.GetMedicalRecord(ctx, params)
	if err != nil {
		return
	}

	return listMedicalRecord, nil
}
