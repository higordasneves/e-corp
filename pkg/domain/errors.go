package domain

import "errors"

var (
	ErrNotFound         = errors.New("not found")
	ErrInvalidParameter = errors.New("invalid parameter")
)
