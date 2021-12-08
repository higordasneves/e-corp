package vos

import (
	domainerr "github.com/higordasneves/e-corp/pkg/domain/errors"
	"regexp"
	"unicode"
)

type (
	CPF string
)

func (cpf CPF) String() string {
	return string(cpf)
}

//ValidateInput validates a CPF
func (cpf *CPF) ValidateInput() error {
	if err := cpf.validateInputLen(); err != nil {
		return err
	}
	return cpf.validateInputFormat()
}

//validateInputLen validates the CPF length
func (cpf *CPF) validateInputLen() error {
	if len(*cpf) != 11 {
		return domainerr.ErrCPFLen
	}
	return nil

}

//validateInputFormat validates if the CPF has only numbers
func (cpf *CPF) validateInputFormat() error {
	for _, v := range *cpf {
		if !unicode.IsDigit(v) {
			return domainerr.ErrCPFFormat
		}
	}
	return nil
}

//FormatOutput formats CPF of account owner to pattern "xxx-xxx-xxx-xx"
func (cpf *CPF) FormatOutput() {
	cpfModel, err := regexp.Compile(`^([\d]{3})([\d]{3})([\d]{3})([\d]{2})$`)
	if err == nil {
		*cpf = CPF(cpfModel.ReplaceAllString(cpf.String(), "$1.$2.$3-$4"))
	}
}
