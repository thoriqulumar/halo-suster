package repo

import (
	"context"
	"helo-suster/model"

	"github.com/jmoiron/sqlx"
)

type MedicalRepo interface {
	NewTx() (*sqlx.Tx, error)
	CreatePatient(ctx context.Context, request model.PostPatientRequest) (patient model.Patient, err error)
	GetPatientByIdentityNumber(ctx context.Context, identityNumber int) (patient model.Patient, err error)
	CreateMedicalRecord(ctx context.Context, requestData model.PostMedicalRecordRequest, createdBy string) (medicalRecord model.MedicalRecord, err error)
}

type medicalRepo struct {
	db *sqlx.DB
}

func NewMedicalRepo(db *sqlx.DB) MedicalRepo {
	return &medicalRepo{
		db: db,
	}
}

func (r *medicalRepo) NewTx() (*sqlx.Tx, error) {
	return r.db.Beginx()
}

var (
	createPatientQuery = `INSERT INTO "patient" ("identityNumber", "phoneNumber", "name", "birthDate", "gender", "identityCardScanImg", "createdAt") 
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
		RETURNING *;`
)

func (r *medicalRepo) CreatePatient(ctx context.Context, requestData model.PostPatientRequest) (patient model.Patient, err error) {
 
	err = r.db.QueryRowxContext(ctx, createPatientQuery, requestData.IdentityNumber, 
		requestData.PhoneNumber, requestData.Name, requestData.BirthDate, requestData.Gender, requestData.IdentityCardScanImg ).StructScan(&patient)
	
	if err != nil{
		return
	}

	return patient, nil
}

var (
	getPatientByIdentityNumberQuery = `SELECT * FROM "patient" WHERE "identityNumber" = $1 LIMIT 1;`
)

func (r *medicalRepo) GetPatientByIdentityNumber(ctx context.Context, identityNumber int) (patient model.Patient, err error) {
	err = r.db.QueryRowxContext(ctx, getPatientByIdentityNumberQuery, identityNumber).StructScan(&patient)

	if err != nil{
		return
	}

	return patient, nil
}


var (
	createMedicalRecordQuery = `INSERT INTO "medicalRecord" ("identityNumber", "symptoms", "medications",  "createdAt", "createdBy") 
		VALUES ($1, $2, $3, NOW(), $4)
		RETURNING *;`
)

func (r *medicalRepo) CreateMedicalRecord(ctx context.Context, requestData model.PostMedicalRecordRequest, createdBy string) (medicalRecord model.MedicalRecord, err error) {
 
	err = r.db.QueryRowxContext(ctx, createMedicalRecordQuery, requestData.IdentityNumber, 
		requestData.Symptoms, requestData.Medications, createdBy).StructScan(&medicalRecord)
	
	if err != nil{
		return
	}

	return medicalRecord, nil
}