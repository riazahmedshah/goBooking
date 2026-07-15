package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateIdempotencyKey() (string, error) {
	return uuid.New().String(), nil
}

func ValidateIdempotencyKey(key string) (bool, error) {
	_, err := uuid.Parse(key)
	if err != nil {
		return false, err
	}
	return true, nil
}

type CustomJWTClaims struct {
	UserID string `json:"userId"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWTToken(userId, role string, jwtSecretKey []byte) (string, error) {
	claim := &CustomJWTClaims{
		UserID: userId,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString(jwtSecretKey)
}
