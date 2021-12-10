package vos

import (
	"errors"
	"regexp"
	"unicode"
)

type (
	CPF string
)

var (
	//ErrCPFLen occurs when the cpf received have invalid length
	ErrCPFLen = errors.New("the CPF must be 11 characters long")
	//ErrCPFFormat occurs when the cpf contains invalid characters
	ErrCPFFormat = errors.New("the CPF must contain only numbers")
)

var cpfModel = regexp.MustCompile(`^([\d]{3})([\d]{3})([\d]{3})([\d]{2})$`)

func (cpf CPF) String() string {
	return string(cpf)
}

//ValidateInput validates a CPF
func (cpf CPF) ValidateInput() error {
	if err := cpf.validateInputFormat(); err != nil {
		return err
	}
	return cpf.validateInputLen()
}

//validateInputLen validates the CPF length
func (cpf CPF) validateInputLen() error {
	if len(cpf) != 11 {
		return ErrCPFLen
	}
	return nil

}

//validateInputFormat validates if the CPF has only numbers
func (cpf CPF) validateInputFormat() error {
	for _, v := range cpf {
		if !unicode.IsDigit(v) {
			return ErrCPFFormat
		}
	}
	return nil
}

//FormatOutput formats CPF of account owner to pattern "xxx-xxx-xxx-xx"
func (cpf *CPF) FormatOutput() string {
	return cpfModel.ReplaceAllString(cpf.String(), "$1.$2.$3-$4")
}
