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

// NewSecret returns hash of password.
// If the number of character is less than minSecretLength returns ErrSmallSecret.
func NewSecret(s string) (Secret, error) {
	if len(s) < minSecretLength {
		return "", ErrSmallSecret
	}

	hashSecret, err := bcrypt.GenerateFromPassword([]byte(s), 10)
	if err != nil {
		return "", fmt.Errorf("unexpeted error generating secret: %w", err)
	}

	return Secret(hashSecret), nil
}

// CompareHashSecret compares password sent by user and stored password hash.
// Returns ErrInvalidPass if the provided secret doesn't match with the hash password.
func (hashSecret Secret) CompareHashSecret(secret string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashSecret), []byte(secret))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidPass
		}
		return fmt.Errorf("unexpeted error comparing password: %w", err)
	}

	return nil
}
