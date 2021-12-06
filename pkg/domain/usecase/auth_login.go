package usecase

import (
	"context"
	domainerr "github.com/higordasneves/e-corp/pkg/domain/errors"
	"github.com/higordasneves/e-corp/pkg/domain/models"
	"time"
)

//LoginInput represents information necessary to access a bank account
type LoginInput struct {
	CPF    string `json:"cpf"`
	Secret string `json:"secret"`
}

func (authUC authUseCase) Login(ctx context.Context, input *LoginInput) (*models.Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	acc, err := authUC.accountRepo.GetAccount(ctx, input.CPF)
	if err != nil {
		return nil, err
	}

	err = acc.CompareHashSecret(input.Secret)
	if err != nil {
		return nil, domainerr.ErrInvalidPass
	}

	return acc, nil
}
