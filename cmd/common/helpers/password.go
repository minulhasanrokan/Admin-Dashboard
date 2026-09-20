package helpers

import (
	"admin-dashboard/cmd/internal/constants"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), constants.BcryptCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckPasswordHash(password, hash string) (bool, error) {

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	if err != nil {

		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {

			return false, nil
		}
		return false, err
	}
	return true, nil
}
