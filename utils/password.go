package utils

import (
	"bytes"
	"crypto/rand"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const bcryptMaxPassLen int = 72

var ErrCouldNotHash error = errors.New("could generate hash of password")

func GenerateSalt() ([]byte, error) {
	salt := make([]byte, bcryptMaxPassLen)
	_, err := rand.Read(salt)
	return salt, err
}

func combineSalt(password, salt []byte) []byte {
	combined := bytes.Clone(password)
	combined = append(combined, salt...)
	if len(combined) > bcryptMaxPassLen {
		combined = combined[:bcryptMaxPassLen]
	}
	return combined
}

func SaltyPassword(password, salt []byte) ([]byte, error) {
	fullPassword := combineSalt(password, salt)
	hashed, err := bcrypt.GenerateFromPassword(fullPassword, bcrypt.DefaultCost)
	if err != nil {
		return nil, ErrCouldNotHash
	}
	return hashed, nil
}

func IsPassword(hashed, password, salt []byte) bool {
	combined := combineSalt(password, salt)
	err := bcrypt.CompareHashAndPassword(hashed, combined)
	return err == nil
}
