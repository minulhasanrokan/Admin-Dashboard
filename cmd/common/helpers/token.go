package helpers

import (
	"crypto/rand"
)

func GenerateVerificationToken(length int) (string, error) {

	bytes := make([]byte, length)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	token := make([]byte, length)
	for i, b := range bytes {
		token[i] = '0' + (b % 10)
	}

	return string(token), nil
}

func GeneratePasswordResetToken(length int) (string, error) {

	bytes := make([]byte, length)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	token := make([]byte, length)
	for i, b := range bytes {
		token[i] = '0' + (b % 10)
	}

	return string(token), nil
}
