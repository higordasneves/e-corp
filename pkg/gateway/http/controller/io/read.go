package io

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrReadRequest = errors.New("invalid request body")
	ErrTokenFormat = errors.New("invalid token format")
)

func ReadRequestBody(r *http.Request, obj interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(&obj); err != nil {
		return fmt.Errorf("%w: %s", ErrReadRequest, err)
	}
	return nil
}
