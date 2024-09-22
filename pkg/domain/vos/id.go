package vos

import (
	"errors"
	"github.com/gofrs/uuid/v5"
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

// NewUUID gets uuid using google lib
func NewUUID() UUID {
	return UUID(uuid.Must(uuid.NewV7()).String())
}

// IsValidUUID validates uuid
func IsValidUUID(id string) error {
	u := uuid.FromStringOrNil(id)
	if u.IsNil() {
		return ErrInvalidID
	}

	return nil
}
