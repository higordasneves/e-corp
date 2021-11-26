package errors

import "errors"

// account errors
var (
	ErrUnexpected  = errors.New("an unexpected error has occurred trying to process your request")
	ErrAccNotFound = errors.New("account not found")
	ErrEmptyInput  = errors.New("the name, document and password fields are required")
	ErrSmallSecret = errors.New("the password must be at least 8 characters long")
	ErrCPFLen      = errors.New("the CPF must be 11 characters long")
	ErrCPFFormat   = errors.New("the CPF must contain only numbers")
)
