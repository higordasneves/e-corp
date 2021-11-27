package vos

import (
	"errors"
	"github.com/google/uuid"
)

type (
	UUID string
)

var (
	ErrInvalidID = errors.New("error, invalid id")
)

func (id UUID) String() string {
	return string(id)
}

//NewAccID gets uuid using google lib
func NewAccID() UUID {
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
