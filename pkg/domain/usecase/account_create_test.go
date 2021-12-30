package usecase

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"github.com/higordasneves/e-corp/pkg/repository"
	repomock "github.com/higordasneves/e-corp/pkg/repository/mock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAccountUseCase_ValidateAccountInput(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		accInput    *AccountInput
		expectedErr error
	}{
		{
			name: "with success",
			accInput: &AccountInput{
				Name:   "Elliot",
				CPF:    "33344455566",
				Secret: "password",
			},
			expectedErr: nil,
		},
		{
			name: "with success, remove blank spaces",
			accInput: &AccountInput{
				Name: "  Elliot  ",
				CPF:  "  33344455567",
				Secret: "password   	",
			},
			expectedErr: nil,
		},
		{
			name: "small password",
			accInput: &AccountInput{
				Name:   "Elliot",
				CPF:    "33344455568",
				Secret: "passwor",
			},
			expectedErr: vos.ErrSmallSecret,
		},
		{
			name: "empty field",
			accInput: &AccountInput{
				Name:   "",
				CPF:    "33344455568",
				Secret: "password",
			},
			expectedErr: entities.ErrEmptyInput,
		},
		{
			name: "CPF format",
			accInput: &AccountInput{
				Name:   "Elliot",
				CPF:    "333.444.555-68",
				Secret: "password",
			},
			expectedErr: vos.ErrCPFFormat,
		},
		{
			name: "CPF length",
			accInput: &AccountInput{
				Name:   "Elliot",
				CPF:    "3334445556",
				Secret: "password",
			},
			expectedErr: vos.ErrCPFLen,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.accInput.ValidateAccountInput()
			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}

func TestAccountUseCase_CreateAccount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		accInput    *AccountInput
		want        *entities.AccountOutput
		expectedErr error
	}{
		{
			name: "with success",
			accInput: &AccountInput{
				Name:   "Elliot",
				CPF:    "33344455566",
				Secret: "password",
			},
			want: &entities.AccountOutput{
				Name:      "Elliot",
				CPF:       "333.444.555-66",
				Balance:   1000000,
				CreatedAt: time.Now().Truncate(time.Minute),
			},
			expectedErr: nil,
		},
		{
			name: "database error: account already exists",
			accInput: &AccountInput{
				Name:   "Elliot",
				CPF:    "33344455566",
				Secret: "password",
			},
			want:        nil,
			expectedErr: entities.ErrAccAlreadyExists,
		},
		{
			name: "database generic error",
			accInput: &AccountInput{
				Name:   "Elliot",
				CPF:    "33344455567",
				Secret: "password",
			},
			want:        nil,
			expectedErr: repository.ErrUnexpected,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// setup
			accRepo := repomock.NewAccountRepo([]entities.Account{}, tt.expectedErr)
			accUseCase := NewAccountUseCase(accRepo)

			// execute
			ctx := context.Background()
			acc, err := accUseCase.CreateAccount(ctx, tt.accInput)
			if tt.want != nil {
				tt.want.ID = acc.ID
				tt.want.CreatedAt = acc.CreatedAt
			}

			//assert
			assert.ErrorIs(t, err, tt.expectedErr)
			assert.Equal(t, tt.want, acc)
		})
	}
}
