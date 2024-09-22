package domain

import "errors"

type ErrorType int

const (
	UndefinedErrorType ErrorType = iota
	NotFoundErrorType
	InvalidParamErrorType
)

type domainError struct {
	errType ErrorType
	msg     string
	fields  map[string]string
}

func (e *domainError) Error() string {
	return e.msg
}

func Error(errType ErrorType, msg string, fields map[string]string) error {
	return &domainError{errType, msg, fields}
}

// GetErrorType returns the ErrorType from a domainError.
// It returns UndefinedErrorType in case the result of errors.As(err, &domainErr) is false.
func GetErrorType(err error) ErrorType {
	domainErr := new(domainError)
	errors.As(err, &domainErr)

	// UndefinedErrorCategory will be returned in case errors.As(err, &domainErr) returned false.
	return domainErr.errType
}
