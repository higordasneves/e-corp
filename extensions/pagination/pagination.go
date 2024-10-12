package pagination

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
)

const (
	maxSize     = 200
	defaultSize = 10
)

var ErrInvalidToken = errors.New("invalid page token")

func ValidatePageSize(size uint32) int {
	switch {
	case size == 0:
		return defaultSize
	case size > maxSize:
		return maxSize
	default:
		return int(size)
	}
}

func Extract(token string, page interface{}) error {
	data, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidToken, err)
	}

	if err = json.Unmarshal(data, page); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidToken, err)
	}

	return nil
}

func NewToken(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrInvalidToken, err)
	}

	token := base64.StdEncoding.EncodeToString(b)
	return token, nil
}
