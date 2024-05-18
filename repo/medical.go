package repo

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"halo-suster/model"
)

type MedicalRepo interface {
	NewTx() (*sqlx.Tx, error)
	CreatePatient(ctx context.Context, request model.PostPatientRequest) (patient model.Patient, err error)
	GetPatientByIdentityNumber(ctx context.Context, identityNumber int) (patient model.Patient, err error)
	CreateMedicalRecord(ctx context.Context, requestData model.PostMedicalRecordRequest, createdBy string) (medicalRecord model.MedicalRecord, err error)
	GetPatient(ctx context.Context, params model.GetPatientParams) (patients []model.Patient, err error)
	GetMedicalRecord(ctx context.Context, params model.GetMedicalRecordParams) (records []model.GetMedicalRecordData, err error)
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
		requestData.PhoneNumber, requestData.Name, requestData.BirthDate, requestData.Gender, requestData.IdentityCardScanImg).StructScan(&patient)

	if err != nil {
		return
	}

	return patient, nil
}

var (
	getPatientByIdentityNumberQuery = `SELECT * FROM "patient" WHERE "identityNumber" = $1 LIMIT 1;`
)

func (r *medicalRepo) GetPatientByIdentityNumber(ctx context.Context, identityNumber int) (patient model.Patient, err error) {
	err = r.db.QueryRowxContext(ctx, getPatientByIdentityNumberQuery, identityNumber).StructScan(&patient)

	if err != nil {
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

	if err != nil {
		return
	}

	return medicalRecord, nil
}

func (r *medicalRepo) GetPatient(ctx context.Context, params model.GetPatientParams) (patients []model.Patient, err error) {
	var listPatient []model.Patient
	var getPatientQuery = `SELECT * FROM "patient" WHERE true`

	if params.Name != "" {
		name := "%" + params.Name + "%"
		getPatientQuery += fmt.Sprintf(` AND "name" ILIKE '%s'`, name)
	}

	if params.IdentityNumber != nil {
		getPatientQuery += fmt.Sprintf(` AND "identityNumber" = %d`, *params.IdentityNumber)
	}

	if params.PhoneNumber != "" {
		if len(params.PhoneNumber) < 10 {
			params.PhoneNumber = "%" + params.PhoneNumber + "%"
		}
		getPatientQuery += fmt.Sprintf(` AND "phoneNumber" ILIKE '%s'`, params.PhoneNumber)
	}

	if params.CreatedAt != "" {
		if params.CreatedAt != "desc" && params.CreatedAt != "asc" {
			params.CreatedAt = "desc"
		}
		getPatientQuery += fmt.Sprintf(` ORDER BY "createdAt" %s`, params.CreatedAt)
	} else {
		getPatientQuery += ` ORDER BY "createdAt" DESC`
	}

	if params.Limit == 0 {
		params.Limit = 5 // default limit
	}

	getPatientQuery += fmt.Sprintf(` LIMIT %d OFFSET %d`, params.Limit, params.Offset)

	rows, err := r.db.QueryContext(ctx, getPatientQuery)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// Iterate over the rows and scan each row into a struct
	for rows.Next() {
		var patient model.Patient
		if err := rows.Scan(&patient.IdentityNumber, &patient.PhoneNumber, &patient.Name, &patient.BirthDate, &patient.Gender, &patient.IdentityCardScanImg, &patient.CreatedAt); err != nil {
			return nil, err
		}

		listPatient = append(listPatient, patient)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return listPatient, nil
}

var (
	getMedicalRecordQuery = `SELECT 
							p."identityNumber", p."phoneNumber", p.name, p."birthDate", p.gender, p."identityCardScanImg",
							mr.symptoms, mr.medications, mr."createdAt",
							s.nip, s.name, s.id as userId
						FROM "medicalRecord" mr
						JOIN patient p ON mr."identityNumber" = p."identityNumber"
						JOIN staff s ON mr."createdBy" = s.id
					    WHERE true`
)

func (r *medicalRepo) GetMedicalRecord(ctx context.Context, params model.GetMedicalRecordParams) (records []model.GetMedicalRecordData, err error) {
	var listRecord []model.GetMedicalRecordData
	var getMedicalRecordQuery = getMedicalRecordQuery

	if params.IdentityNumber != nil {
		getMedicalRecordQuery += fmt.Sprintf(` AND "p.identityNumber" = %d`, params.IdentityNumber)
	}

	if params.CreatedByUserId != "" {
		getMedicalRecordQuery += fmt.Sprintf(` AND "s.id" = %s`, params.CreatedByUserId)
	}
	if params.CreatedByNip != "" {
		getMedicalRecordQuery += fmt.Sprintf(` AND "s.nip" = %s`, params.CreatedByNip)
	}

	if params.CreatedAt != "" {
		if params.CreatedAt != "desc" && params.CreatedAt != "asc" {
			params.CreatedAt = "desc"
		}
		getMedicalRecordQuery += fmt.Sprintf(` ORDER BY "createdAt" %s`, params.CreatedAt)
	} else {
		getMedicalRecordQuery += ` ORDER BY "createdAt" DESC`
	}

	if params.Limit == 0 {
		params.Limit = 5 // default limit
	}

	getMedicalRecordQuery += fmt.Sprintf(` LIMIT %d OFFSET %d`, params.Limit, params.Offset)

	rows, err := r.db.QueryContext(ctx, getMedicalRecordQuery)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// Iterate over the rows and scan each row into a struct
	for rows.Next() {
		var medicalRecord model.GetMedicalRecordData
		if err := rows.Scan(
			&medicalRecord.IdentityDetail.IdentityNumber,
			&medicalRecord.IdentityDetail.PhoneNumber,
			&medicalRecord.IdentityDetail.Name,
			&medicalRecord.IdentityDetail.BirthDate,
			&medicalRecord.IdentityDetail.Gender,
			&medicalRecord.IdentityDetail.IdentityCardScanImg,
			&medicalRecord.Symptoms,
			&medicalRecord.Medications,
			&medicalRecord.CreatedAt,
			&medicalRecord.CreatedBy.Nip,
			&medicalRecord.CreatedBy.Name,
			&medicalRecord.CreatedBy.UserId,
		); err != nil {
			return nil, err
		}

		listRecord = append(listRecord, medicalRecord)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return listRecord, nil
}
