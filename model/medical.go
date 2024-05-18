package model

type MedicalGeneralResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PostPatientRequest struct {
	IdentityNumber      int    `json:"identityNumber" validate:"required"`
	PhoneNumber         string `json:"phoneNumber" validate:"required,min=10,max=15,phone_number"`
	Name                string `json:"name" validate:"required,min=3,max=30"`
	BirthDate           string `json:"birthDate" validate:"required"`
	Gender              string `json:"gender" validate:"required"`
	IdentityCardScanImg string `json:"identityCardScanImg" validate:"required,custom_url"`
}

type Patient struct {
	IdentityNumber      int    `json:"identityNumber" db:"identityNumber"`
	PhoneNumber         string `json:"phoneNumber" db:"phoneNumber"`
	Name                string `json:"name" db:"name"`
	BirthDate           string `json:"birthDate" db:"birthDate"`
	Gender              string `json:"gender" db:"gender"`
	IdentityCardScanImg string `json:"identityCardScanImg" db:"identityCardScanImg"`
	CreatedAt           string `json:"createdAt" db:"createdAt"`
}

type PostMedicalRecordRequest struct {
	IdentityNumber int    `json:"identityNumber" validate:"required"`
	Symptoms       string `json:"symptoms" validate:"required,min=1,max=2000"`
	Medications    string `json:"medications" validate:"required,min=1,max=2000"`
}

type MedicalRecord struct {
	IdentityNumber int    `json:"identityNumber" db:"identityNumber"`
	Symptoms       string `json:"symptoms" db:"symptoms"`
	Medications    string `json:"medications" db:"medications"`
	CreatedAt      string `json:"createdAt" db:"createdAt"`
	CreatedBy      string `json:"createdBy" db:"createdBy"`
}

type GetPatientParams struct {
	IdentityNumber *int
	Name           string
	PhoneNumber    string
	Limit          int
	Offset         int
	CreatedAt      string
}

type GetMedicalRecordParams struct {
	IdentityNumber  *int
	CreatedByUserId string
	CreatedByNip    string
	Limit           int
	Offset          int
	CreatedAt       string
}

type GetMedicalRecordData struct {
	IdentityDetail Patient `json:"identityDetail"`
	Symptoms       string  `json:"symptoms" db:"symptoms"`
	Medications    string  `json:"medications" db:"medications"`
	CreatedAt      string  `json:"createdAt" db:"createdAt"`
	CreatedBy      struct {
		Nip    string `json:"nip"`
		Name   string `json:"name"`
		UserId string `json:"userId"`
	} `json:"createdBy"`
}
