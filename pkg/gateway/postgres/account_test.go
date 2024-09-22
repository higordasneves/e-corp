package postgres

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/higordasneves/e-corp/pkg/domain"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
)

func TestAccRepo_CreateAccount(t *testing.T) {
	t.Parallel()
	db := NewDB(t)

	tests := []struct {
		name        string
		input       entities.Account
		wantErrType domain.ErrorType
	}{
		{
			name: "success",
			input: entities.Account{
				ID:        vos.NewUUID(),
				Name:      "Elliot",
				CPF:       "33344455566",
				Secret:    "password",
				Balance:   0,
				CreatedAt: time.Now().Truncate(time.Second),
			},
			wantErrType: 0,
		},
		{
			name: "fail - account already exists",
			input: entities.Account{
				ID:        vos.NewUUID(),
				Name:      "Elliot",
				CPF:       "33344455566",
				Secret:    "password",
				Balance:   0,
				CreatedAt: time.Now().Truncate(time.Second),
			},
			wantErrType: domain.InvalidParamErrorType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accRepo := NewRepository(db)
			err := accRepo.CreateAccount(context.Background(), &tt.input)
			assert.Equal(t, tt.wantErrType, domain.GetErrorType(err))
		})
	}
}

func TestAccRepo_FetchAccounts(t *testing.T) {
	t.Parallel()

	repo := NewRepository(NewDB(t))

	accounts := []entities.Account{
		{
			ID:        vos.NewUUID(),
			Name:      "Elliot",
			CPF:       "33344455567",
			Secret:    "password",
			Balance:   7000,
			CreatedAt: time.Now().Truncate(time.Second),
		},
		{
			ID:        vos.NewUUID(),
			Name:      "Mr.Robot",
			CPF:       "33344455568",
			Secret:    "password",
			Balance:   3000,
			CreatedAt: time.Now().Truncate(time.Second),
		},
	}
	for _, acc := range accounts {
		err := repo.CreateAccount(context.Background(), &acc)
		if err != nil {
			t.Error("error inserting accounts")
		}
	}

	// execute
	result, err := repo.FetchAccounts(context.Background())
	if err != nil {
		t.Errorf("didn't want sql error, but got the error: %v", err)
	}

	//assert
	if !reflect.DeepEqual(accounts, result) {
		t.Errorf("got: %v, want: %v", result, accounts)
	}
}

func TestAccRepo_GetBalance(t *testing.T) {
	t.Parallel()

	accRepo := NewRepository(NewDB(t))

	tests := []struct {
		name        string
		acc         *entities.Account
		insert      bool
		expectedErr bool
		err         error
	}{
		{
			name: "with success",
			acc: &entities.Account{
				ID:        vos.NewUUID(),
				Name:      "Elliot",
				CPF:       "33344455567",
				Secret:    "password",
				Balance:   7000,
				CreatedAt: time.Now().Truncate(time.Second),
			},
			insert:      true,
			expectedErr: false,
			err:         nil,
		},
		{
			name: "with success balance 0",
			acc: &entities.Account{
				ID:        vos.NewUUID(),
				Name:      "Elliot",
				CPF:       "33344455568",
				Secret:    "password",
				Balance:   0,
				CreatedAt: time.Now().Truncate(time.Second),
			},
			insert:      true,
			expectedErr: false,
			err:         nil,
		},
		{
			name: "Repository not found",
			acc: &entities.Account{
				ID: vos.NewUUID(),
			},
			insert:      false,
			expectedErr: true,
			err:         entities.ErrAccNotFound,
		},
	}

	var GotDBError *domain.DBError
	var WantDBError *domain.DBError

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.insert {
				_ = accRepo.CreateAccount(context.Background(), tt.acc)
			}

			// execute
			result, err := accRepo.GetBalance(context.Background(), tt.acc.ID)

			// assert
			switch {
			case errors.As(err, &GotDBError) && errors.As(tt.err, &WantDBError):
				if GotDBError.Query != WantDBError.Query {
					t.Errorf("got sql error in query: %v, want: %v", GotDBError.Query, WantDBError.Query)
				}
			case !errors.Is(err, tt.err):
				t.Errorf("got error: %v, want: %v", err, tt.err)
			case !tt.expectedErr && result != tt.acc.Balance:
				t.Errorf("got: %v, want: %v", result, tt.acc.Balance)
			}
		})
	}
}

func TestAccRepo_UpdateBalance(t *testing.T) {
	tests := []struct {
		name         string
		acc          *entities.Account
		updateAmount int
		insert       bool
		expectedErr  bool
		err          error
	}{
		{
			name: "with success outbound",
			acc: &entities.Account{
				ID:        vos.NewUUID(),
				Name:      "Elliot",
				CPF:       "33344455567",
				Secret:    "password",
				Balance:   7000,
				CreatedAt: time.Now().Truncate(time.Second),
			},
			updateAmount: -5000,
			insert:       true,
			expectedErr:  false,
			err:          nil,
		},
		{
			name: "with success inbound",
			acc: &entities.Account{
				ID:        vos.NewUUID(),
				Name:      "Elliot",
				CPF:       "33344455568",
				Secret:    "password",
				Balance:   0,
				CreatedAt: time.Now().Truncate(time.Second),
			},
			updateAmount: 5000,
			insert:       true,
			expectedErr:  false,
			err:          nil,
		},
	}

	var GotDBError *domain.DBError
	var WantDBError *domain.DBError

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			accRepo := NewRepository(NewDB(t))
			if tt.insert {
				_ = accRepo.CreateAccount(context.Background(), tt.acc)
			}

			// execute
			err := accRepo.UpdateBalance(context.Background(), tt.acc.ID, tt.updateAmount)

			// assert
			switch {
			case errors.As(err, &GotDBError) && errors.As(tt.err, &WantDBError):
				if GotDBError.Query != WantDBError.Query {
					t.Errorf("got sql error in query: %v, want: %v", GotDBError.Query, WantDBError.Query)
				}
			case err != tt.err:
				t.Errorf("got error: %v, want: %v", err, tt.err)
			case !tt.expectedErr:
				gotBalance, errGetBalance := accRepo.GetBalance(context.Background(), tt.acc.ID)
				if errGetBalance != nil {
					t.Error("unexpected error in get balance query")
				} else if gotBalance != tt.acc.Balance+tt.updateAmount {
					t.Errorf("got: %v, want: %v", gotBalance, tt.acc.Balance+tt.updateAmount)
				}
			}
		})
	}
}

func TestAccRepo_GetAccount(t *testing.T) {
	tests := []struct {
		name        string
		acc         *entities.Account
		insert      bool
		expectedErr bool
		err         error
	}{
		{
			name: "with success",
			acc: &entities.Account{
				ID:        vos.NewUUID(),
				Name:      "Elliot",
				CPF:       "33344455567",
				Secret:    "password",
				Balance:   7000,
				CreatedAt: time.Now().Truncate(time.Second),
			},
			insert:      true,
			expectedErr: false,
			err:         nil,
		},
		{
			name: "Repository not found",
			acc: &entities.Account{
				ID: vos.NewUUID(),
			},
			insert:      false,
			expectedErr: true,
			err:         entities.ErrAccNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			accRepo := NewRepository(NewDB(t))
			if tt.insert {
				_ = accRepo.CreateAccount(context.Background(), tt.acc)
			}

			// execute
			result, err := accRepo.GetAccount(context.Background(), tt.acc.CPF)
			if tt.expectedErr && err != tt.err {
				t.Errorf("got: %v, want: %v", err, tt.err)
			}

			// assert
			if !tt.expectedErr && reflect.DeepEqual(&result, tt.acc) {
				t.Errorf("got: %v, want: %v", &result, tt.acc)
			}
		})
	}
}
