package model

type Staff struct {
	ID                  string `json:"id" db:"id"`
	NIP                 int64  `json:"nip" db:"nip"`
	Name                string `json:"name" db:"name"`
	Role                string `json:"-" db:"role,omitempty"`
	IdentityCardScanImg string `json:"-" db:"identityCardScanImg,omitempty"`
	CreatedAt           string `json:"createdAt" db:"createdAt"`
}

type GetStaffRequest struct {
	UserId    *string `schema:"userId"`
	Limit     *int32  `schema:"limit"`
	Offset    *int32  `schema:"offset"`
	Name      *string `schema:"name"`
	Nip       *int32  `schema:"nip"`
	Role      *string `schema:"role"`
	CreatedAt *string `schema:"createdAt"`
}

type GetStaffResponse struct {
	Message string  `json:"message"`
	Data    []Staff `json:"data"`
}
