package vos

import (
	"errors"
	"unicode"
)

type (
	Document string
)

var (
	// ErrDocumentLen occurs when the cpf received have invalid length.
	ErrDocumentLen = errors.New("the document must have 11 or 14 characters")
	// ErrDocumentFormat occurs when the cpf contains invalid characters.
	ErrDocumentFormat = errors.New("the document must contain only numbers")
)

func (cpf Document) String() string {
	return string(cpf)
}

// NewDocument creates a new document from a string.
// returns ErrDocumentLen if the number of the digits is invalid.
// returns ErrDocumentFormat id the format of the document is invalid.
func NewDocument(s string) (Document, error) {
	if err := validateDocumentLen(s); err != nil {
		return "", err
	}

	if err := validateDocumentFormat(s); err != nil {
		return "", err
	}

	return Document(s), nil
}

// validateDocumentLen validates the Document length
func validateDocumentLen(s string) error {
	if n := len(s); n != 11 && n != 14 {
		return ErrDocumentLen
	}

	return nil
}

// validateInputFormat validates if the Document has only numbers
func validateDocumentFormat(s string) error {
	for _, v := range s {
		if !unicode.IsDigit(v) {
			return ErrDocumentFormat
		}
	}

	return nil
}
