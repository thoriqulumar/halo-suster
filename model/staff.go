package model

type Staff struct {
	ID                  string `json:"id" db:"id"`
	NIP                 int64  `json:"nip" db:"nip"`
	Name                string `json:"name" db:"name"`
	Role                string `json:"role" db:"role,omitempty"`
	IdentityCardScanImg string `json:"identityCardScanImg" db:"identityCardScanImg,omitempty"`
	CreatedAt           string `json:"createdAt" db:"createdAt"`
}
