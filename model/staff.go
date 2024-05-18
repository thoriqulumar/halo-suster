package model

import (
	"github.com/google/uuid"
	"time"
)

type Role string

const (
	RoleUnknown Role = "unknown"
	RoleIt      Role = "it"
	RoleNurse   Role = "nurse"
)

type Status string

const (
	StatusUnknown Status = "unknown"
	StatusActive  Status = "active"
	StatusDeleted Status = "deleted"
)

type Staff struct {
	UserId              uuid.UUID `json:"userId" db:"id"`
	NIP                 int64     `json:"nip" db:"nip"`
	Name                string    `json:"name" db:"name"`
	Role                Role      `json:"role" db:"role"`
	IdentityCardScanImg string    `json:"identityCardScanImg" db:"identityCardScanImg"`
	Status              Status    `json:"status" db:"status"`
	Password            string    `json:"password" db:"password"`
	CreatedAt           time.Time `json:"createdAt" db:"createdAt"`
}

type RegisterStaffRequest struct {
	NIP      int64  `json:"nip" validate:"required,nip"`
	Name     string `json:"name" validate:"required,min=5,max=50"`
	Password string `json:"password" validate:"required,min=5,max=33"`
}

type RegisterNurseRequest struct {
	NIP                 int64  `json:"nip" validate:"required,nip"`
	Name                string `json:"name" validate:"required,min=5,max=50"`
	IdentityCardScanImg string `json:"identityCardScanImg" validate:"required,customURL"`
}

type UpdateNurseRequest struct {
	ID   uuid.UUID `json:"-"`
	NIP  int64     `json:"nip" validate:"required,nip"`
	Name string    `json:"name" validate:"required,min=5,max=50"`
}

type GrantAccessNurseRequest struct {
	Password string `json:"password" validate:"required"`
}

type LoginStaffRequest struct {
	NIP      int64  `json:"nip" validate:"required,nip"`
	Password string `json:"password" validate:"required"`
}

type StaffWithToken struct {
	UserId      string `json:"userId"`
	NIP         int64  `json:"nip"`
	Name        string `json:"name"`
	Password    string `json:"-"`
	AccessToken string `json:"accessToken"`
}

type RegisterStaffResponse struct {
	Message string         `json:"message"`
	Data    StaffWithToken `json:"data"`
}

type GetListUserParams struct {
	ID     *uuid.UUID `json:"userId"`
	Limit  *int       `json:"limit"`
	Offset *int       `json:"offset"`
	Name   *string    `json:"name"`
	NIP    *int64     `json:"nip"`
	Role   *Role      `json:"role"`
	Status *Status    `json:"status"`
	Sort   UserSorting
}

type UserSorting struct {
	CreatedAt *string
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
