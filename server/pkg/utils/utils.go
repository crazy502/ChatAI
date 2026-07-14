package utils

import (
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/google/uuid"
)

func GetRandomNumbers(length int) (string, error) {
	if length <= 0 {
		return "", errors.New("random number length must be positive")
	}

	digits := make([]byte, length)
	for index := range digits {
		value, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		digits[index] = byte(value.Int64()) + '0'
	}
	return string(digits), nil
}

func GenerateUUID() string {
	return uuid.New().String()
}
