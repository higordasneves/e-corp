package repository

import (
	"errors"
	"fmt"
)

const (
	QueryRefCreateAcc          = "CreateAccount"
	QueryRefFetchAcc           = "FetchAccounts"
	QueryRefGetAcc             = "GetAccount"
	QueryRefGetBalance         = "GetBalance"
	QueryRefUpdateBalance      = "UpdateBalance"
	QueryRefPerformTransaction = "PerformTransaction"
	QueryRefCreateTransfer     = "CreateTransfer"
)

var (
	ErrTruncDB    = errors.New("an unexpected error occurred while deleting tables")
	ErrUnexpected = errors.New("unexpected error occurred")
)

type DBError struct {
	Query      string
	DBErr      error
	GenericErr error
}

func (dbError *DBError) Error() string {
	return fmt.Sprintf("unexpected sql error occurred on query %s: %s", dbError.Query, dbError.DBErr)
}

func NewDBError(query string, sqlError error, genericError error) *DBError {
	return &DBError{Query: query, DBErr: sqlError, GenericErr: genericError}
}
