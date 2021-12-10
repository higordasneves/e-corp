package repository

import (
	"errors"
	"fmt"
)

const (
	QueryRefCreateAcc     = "CreateAccount"
	QueryRefFetchAcc      = "FetchAccounts"
	QueryRefGetAcc        = "GetAccount"
	QueryRefGetBalance    = "GetBalance"
	QueryRefUpdateBalance = "UpdateBalance"
)

var (
	ErrTruncDB = errors.New("an unexpected error occurred while deleting tables")
)

type DBError struct {
	Query string
	Err   error
}

func (dbError *DBError) Error() string {
	return fmt.Sprintf("unexpected sql error occurred on query %s: %s", dbError.Query, dbError.Err)
}

func NewDBError(query string, err error) *DBError {
	return &DBError{Query: query, Err: err}
}
