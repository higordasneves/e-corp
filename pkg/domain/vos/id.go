package vos

import (
	"errors"
	"github.com/google/uuid"
)

type (
	AccountID string
)

var (
	ErrInvalidID = errors.New("error, invalid id")
)

func (accID AccountID) String() string {
	return string(accID)
}

// newAccID gets uuid using google lib
func NewAccID() AccountID {
	return AccountID(uuid.NewString())
}

// IsValidUUID validates uuid
func IsValidUUID(accID string) error {
	_, err := uuid.Parse(accID)
	if err != nil {
		return ErrInvalidID
	}
	return nil
}
