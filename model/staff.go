package model

import (
	"github.com/google/uuid"
)

type Staff struct {
	UserId    uuid.UUID `json:"userId" db:"id"`
	NIP       int64     `json:"nip" db:"nip"`
	Name      string    `json:"name" db:"name"`
	Role      string    `json:"role" db:"role"`
	Password  string    `json:"password" db:"password"`
	CreatedAt string    `json:"createdAt" db:"createdAt"`
}

type RegisterStaffRequest struct {
	NIP      int64  `json:"nip" validate:"required,nip"`
	Name     string `json:"name" validate:"required,min=5,max=50"`
	Password string `json:"password" validate:"required,min=5,max=33"`
}

type StaffWithToken struct {
	UserId string `json:"userId"`
	//NIP         string `json:"nip"`
	Name        string `json:"name"`
	Password    string `json:"-"`
	AccessToken string `json:"accessToken"`
}

type RegisterStaffResponse struct {
	Message string         `json:"message"`
	Data    StaffWithToken `json:"data"`
}
