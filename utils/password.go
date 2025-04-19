package utils

import (
	"bytes"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var ErrCouldNotHash error = errors.New("could generate hash of password")

func combineSalt(password, salt []byte) []byte {
	combined := bytes.Clone(password)
	combined = append(combined, salt...)
	if len(combined) > 72 {
		combined = combined[:72]
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
