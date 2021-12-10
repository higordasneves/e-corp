package usecase

import (
	"context"
	"errors"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"github.com/higordasneves/e-corp/pkg/repository"
	repomock "github.com/higordasneves/e-corp/pkg/repository/mock"
	"reflect"
	"testing"
	"time"
)

func TestAccountUseCase_ValidateAccountInput(t *testing.T) {
	tests := []struct {
		name        string
		accInput    *AccountInput
		want        *AccountInput
		expectedErr error
	}{
		{
			name: "with success",
			accInput: &AccountInput{
				Name:   "Elliot",
				CPF:    "33344455566",
				Secret: "password",
			},
			want: &AccountInput{
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
			want: &AccountInput{
				Name:   "Elliot",
				CPF:    "33344455567",
				Secret: "password",
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
			want:        nil,
			expectedErr: vos.ErrSmallSecret,
		},
		{
			name: "empty field",
			accInput: &AccountInput{
				Name:   "",
				CPF:    "33344455568",
				Secret: "password",
			},
			want:        nil,
			expectedErr: entities.ErrEmptyInput,
		},
		{
			name: "CPF format",
			accInput: &AccountInput{
				Name:   "Elliot",
				CPF:    "333.444.555-68",
				Secret: "password",
			},
			want:        nil,
			expectedErr: vos.ErrCPFFormat,
		},
		{
			name: "CPF length",
			accInput: &AccountInput{
				Name:   "Elliot",
				CPF:    "3334445556",
				Secret: "password",
			},
			want:        nil,
			expectedErr: vos.ErrCPFLen,
		},
	}
	for _, test := range tests {
		err := test.accInput.ValidateAccountInput()
		t.Run(test.name, func(t *testing.T) {
			t.Helper()
			if test.expectedErr != err {
				switch {
				case test.expectedErr == nil:
					t.Errorf("didn't want an error, but got the error: %v", err)
				case err == nil:
					t.Error("wanted an error but didn't get one")
				default:
					t.Errorf("got error: %v, want error: %v", err, test.expectedErr)
				}
			}
			if test.want != nil && !reflect.DeepEqual(test.accInput, test.want) {
				t.Errorf("got %v, want %v", test.accInput, test.want)
			}
		})
	}
}

func TestAccountUseCase_CreateAccount(t *testing.T) {
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
				CreatedAt: time.Now().Truncate(time.Hour),
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
			expectedErr: repository.NewDBError(repository.QueryRefCreateAcc, errors.New("any db error")),
		},
	}

	for _, test := range tests {
		accRepo := repomock.NewAccountRepo([]entities.Account{}, test.expectedErr)
		accUseCase := NewAccountUseCase(accRepo)
		acc, err := accUseCase.CreateAccount(context.Background(), test.accInput)

		t.Run(test.name, func(t *testing.T) {
			t.Helper()
			if test.expectedErr != err {
				switch {
				case test.expectedErr == nil:
					t.Errorf("didn't want an error, but got the error: %v", err)
				case err == nil:
					t.Error("wanted an error but didn't get one")
				default:
					t.Errorf("got error: %v, want error: %v", err, test.expectedErr)
				}
			}

			if err == nil {
				acc.CreatedAt = acc.CreatedAt.Truncate(time.Hour)
				if vos.IsValidUUID(acc.ID) != nil {
					t.Error("account was created with invalid id")
				}
				test.want.ID = acc.ID
			}

			if !reflect.DeepEqual(acc, test.want) {
				t.Errorf("got %v, want %v", acc, test.want)
			}
		})
	}
}
