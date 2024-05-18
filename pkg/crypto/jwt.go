package crypto

import (
	"fmt"
	"halo-suster/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(staff model.Staff, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, model.JWTClaims{
		UserId: staff.UserId.String(),
		NIP:    fmt.Sprint("%d", staff.NIP),
		Name:   staff.Name,
		Role:   string(staff.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Minute)),
		},
	})

	tokenString, err := token.SignedString([]byte(secret))
	return tokenString, err
}

func VerifyToken(token, secretKey string) (*model.JWTPayload, error) {
	claims := &model.JWTClaims{}

	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims.RegisteredClaims.ExpiresAt.Before(time.Now()) {
		return nil, err
	}

	payload := &model.JWTPayload{
		Id:   claims.UserId,
		NIP:  claims.NIP,
		Name: claims.Name,
		Role: model.Role(claims.Role),
	}

	return payload, nil
}
