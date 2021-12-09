package vos

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type Secret string

var (
	//ErrSmallSecret occurs when the secret have invalid length
	ErrSmallSecret = errors.New("the password must be at least 8 characters long")
	//ErrInvalidPass occurs when secret is invalid
	ErrInvalidPass = errors.New("invalid password")
)

func (hashSecret Secret) String() string {
	return string(hashSecret)
}

//GetHashSecret returns hash of password
func GetHashSecret(secret string) (Secret, error) {
	hashSecret, err := bcrypt.GenerateFromPassword([]byte(secret), 10)
	if err != nil {
		return "", err
	}
	return Secret(hashSecret), nil
}

//CompareHashSecret compares password sent by user and stored password hash
func (hashSecret Secret) CompareHashSecret(secret string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashSecret), []byte(secret))
	if err != nil {
		return err
	}
	return nil
}
