package postgres

import (
	"context"
	"errors"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"github.com/higordasneves/e-corp/pkg/repository"
	"reflect"
	"testing"
	"time"
)

func TestAccRepo_CreateAccount(t *testing.T) {

	tests := []struct {
		name string
		acc  *entities.Account
		err  error
	}{
		{
			name: "with success",
			acc: &entities.Account{
				ID:        vos.NewUUID(),
				Name:      "Elliot",
				CPF:       "33344455566",
				Secret:    "password",
				Balance:   0,
				CreatedAt: time.Now().Truncate(time.Second),
			},
			err: nil,
		},
		{
			name: "check error account already exists",
			acc: &entities.Account{
				ID:        vos.NewUUID(),
				Name:      "Elliot",
				CPF:       "33344455566",
				Secret:    "password",
				Balance:   0,
				CreatedAt: time.Now().Truncate(time.Second),
			},
			err: entities.ErrAccAlreadyExists,
		},
		{
			name: "invalid id",
			acc: &entities.Account{
				ID:        "invalid",
				Name:      "Elliot",
				CPF:       "33344455567",
				Secret:    "password",
				Balance:   0,
				CreatedAt: time.Now().Truncate(time.Second),
			},
			err: repository.NewDBError(repository.QueryRefCreateAcc, errors.New("any sql error"), errors.New("unexpected error")),
		},
	}

	defer ClearDB()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var GotDBError *repository.DBError
			var WantDBError *repository.DBError

			accRepo := NewAccountRepo(dbTest)
			ctxDB := context.Background()
			resultErr := accRepo.CreateAccount(ctxDB, tt.acc)

			switch {
			case errors.As(resultErr, &GotDBError) && errors.As(tt.err, &WantDBError):
				if GotDBError.Query != WantDBError.Query {
					t.Errorf("got sql error in query: %v, want: %v", GotDBError.Query, WantDBError.Query)
				}
			case resultErr != tt.err:
				t.Errorf("got error: %v, want error: %v", resultErr, tt.err)
			}
		})
	}
}

func TestAccRepo_FetchAccounts(t *testing.T) {
	// setup
	accRepo := NewAccountRepo(dbTest)

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
	var want []entities.Account
	for _, acc := range accounts {
		err := accRepo.CreateAccount(context.Background(), &acc)
		if err != nil {
			t.Error("error inserting accounts")
		}
		want = append(want, entities.Account{
			ID:        acc.ID,
			Name:      acc.Name,
			CPF:       acc.CPF,
			Balance:   acc.Balance,
			CreatedAt: acc.CreatedAt,
		})
	}

	defer ClearDB()

	// execute
	result, err := accRepo.FetchAccounts(context.Background())
	if err != nil {
		t.Errorf("didn't want sql error, but got the error: %v", err)
	}

	//assert
	if !reflect.DeepEqual(want, result) {
		t.Errorf("got: %v, want: %v", result, want)
	}
}

func TestAccRepo_GetBalance(t *testing.T) {
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
			name: "account not found",
			acc: &entities.Account{
				ID: vos.NewUUID(),
			},
			insert:      false,
			expectedErr: true,
			err:         entities.ErrAccNotFound,
		},
		{
			name: "invalid id",
			acc: &entities.Account{
				ID: "invalid",
			},
			insert:      false,
			expectedErr: true,
			err:         repository.NewDBError(repository.QueryRefGetBalance, errors.New("any sql error"), errors.New("unexpected error")),
		},
	}

	defer ClearDB()

	var GotDBError *repository.DBError
	var WantDBError *repository.DBError

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// setup
			accRepo := NewAccountRepo(dbTest)
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
			case err != tt.err:
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
		{
			name: "account not found",
			acc: &entities.Account{
				ID: vos.NewUUID(),
			},
			updateAmount: 1000,
			insert:       false,
			expectedErr:  true,
			err:          entities.ErrZeroRowsAffectedUpdateBalance,
		},
		{
			name: "invalid id",
			acc: &entities.Account{
				ID: "invalid",
			},
			updateAmount: 1000,
			insert:       false,
			expectedErr:  true,
			err:          repository.NewDBError(repository.QueryRefUpdateBalance, errors.New("any sql error"), errors.New("unexpected error")),
		},
	}

	defer ClearDB()

	var GotDBError *repository.DBError
	var WantDBError *repository.DBError

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			accRepo := NewAccountRepo(dbTest)
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
			name: "account not found",
			acc: &entities.Account{
				ID: vos.NewUUID(),
			},
			insert:      false,
			expectedErr: true,
			err:         entities.ErrAccNotFound,
		},
	}

	defer ClearDB()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			accRepo := NewAccountRepo(dbTest)
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
