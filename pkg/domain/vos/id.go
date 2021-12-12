package vos

import (
	"errors"
	"github.com/google/uuid"
)

type (
	UUID string
)

var (
	ErrInvalidID = errors.New("invalid id")
)

func (id UUID) String() string {
	return string(id)
}

//NewUUID gets uuid using google lib
func NewUUID() UUID {
	return UUID(uuid.NewString())
}

// IsValidUUID validates uuid
func IsValidUUID(id string) error {
	_, err := uuid.Parse(id)
	if err != nil {
		return ErrInvalidID
	}
	return nil
}
