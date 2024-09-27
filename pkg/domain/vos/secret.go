package vos

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Secret string

const minSecretLength = 8

var (
	// ErrSmallSecret occurs when the secret have invalid length.
	ErrSmallSecret = fmt.Errorf("the password must be at least %v characters long", minSecretLength)
	// ErrInvalidPass occurs when secret is invalid.
	ErrInvalidPass = errors.New("invalid password")
)

func (hashSecret Secret) String() string {
	return string(hashSecret)
}

// NewSecret returns hash of password
func NewSecret(s string) (Secret, error) {
	if len(s) < minSecretLength {
		return "", ErrSmallSecret
	}

	hashSecret, err := bcrypt.GenerateFromPassword([]byte(s), 10)
	if err != nil {
		return "", err
	}
	return Secret(hashSecret), nil
}

// CompareHashSecret compares password sent by user and stored password hash.
func (hashSecret Secret) CompareHashSecret(secret string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashSecret), []byte(secret))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidPass
		}
		return err
	}

	return nil
}
