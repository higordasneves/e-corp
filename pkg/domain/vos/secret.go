package vos

import (
	"golang.org/x/crypto/bcrypt"
)

type Secret string

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

func (hashSecret Secret) CompareHashSecret(secret string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashSecret), []byte(secret))
	if err != nil {
		return err
	}
	return nil
}
