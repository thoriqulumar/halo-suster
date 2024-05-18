package model

import "github.com/golang-jwt/jwt/v5"

type JWTClaims struct {
	UserId string `json:"id"`
	Name   string `json:"name"`
	//NIP    string `json:"nip"`
	jwt.RegisteredClaims
}

type JWTPayload struct {
	Id   string
	Name string
	NIP  string
}
