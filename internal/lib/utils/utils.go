package utils

import "github.com/google/uuid"

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
