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
	accRepo := NewAccountRepo(dbTest)
	ctxDB := context.Background()
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
			err: repository.NewDBError(repository.QueryRefCreateAcc, errors.New("any sql error")),
		},
	}
	defer ClearDB()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var GotDBError *repository.DBError
			var WantDBError *repository.DBError

			resultErr := accRepo.CreateAccount(ctxDB, test.acc)

			switch {
			case errors.As(resultErr, &GotDBError) && !errors.As(test.err, &WantDBError):
				t.Errorf("didn't want sql error, but got the error: %v", resultErr)
			case !errors.As(resultErr, &GotDBError) && errors.As(test.err, &WantDBError):
				t.Error("wanted sql error but didn't get one")
			case errors.As(resultErr, &GotDBError) && errors.As(test.err, &WantDBError):
				if GotDBError.Query != WantDBError.Query {
					t.Errorf("got sql error in query: %v, want: %v", GotDBError.Query, WantDBError.Query)
				}
			case resultErr != test.err:
				t.Errorf("got error: %v, want error: %v", resultErr, test.err)
			}
		})
	}
}

func TestAccRepo_FetchAccounts(t *testing.T) {
	accRepo := NewAccountRepo(dbTest)
	ctxDB := context.Background()

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
		err := accRepo.CreateAccount(ctxDB, &acc)
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

	result, err := accRepo.FetchAccounts(ctxDB)
	if err != nil {
		t.Errorf("didn't want sql error, but got the error: %v", err)
	}

	if !reflect.DeepEqual(want, result) {
		t.Errorf("got: %v, want: %v", result, want)
	}
}

func TestAccRepo_GetBalance(t *testing.T) {
	accRepo := NewAccountRepo(dbTest)
	ctxDB := context.Background()
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
				ID:        vos.NewUUID(),
				Name:      "Elliot",
				CPF:       "33344455569",
				Secret:    "password",
				Balance:   0,
				CreatedAt: time.Now().Truncate(time.Second),
			},
			insert:      false,
			expectedErr: true,
			err:         entities.ErrAccNotFound,
		},
		{
			name: "invalid id",
			acc: &entities.Account{
				ID:        "invalid",
				Name:      "Elliot",
				CPF:       "33344455570",
				Secret:    "password",
				Balance:   0,
				CreatedAt: time.Now().Truncate(time.Second),
			},
			insert:      false,
			expectedErr: true,
			err:         repository.NewDBError(repository.QueryRefGetBalance, errors.New("any sql error")),
		},
	}

	defer ClearDB()

	var GotDBError *repository.DBError
	var WantDBError *repository.DBError

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.insert {
				_ = accRepo.CreateAccount(ctxDB, test.acc)
			}

			result, err := accRepo.GetBalance(context.Background(), test.acc.ID)
			switch {
			case errors.As(err, &GotDBError) && errors.As(test.err, &WantDBError):
				if GotDBError.Query != WantDBError.Query {
					t.Errorf("got sql error in query: %v, want: %v", GotDBError.Query, WantDBError.Query)
				}
			case err != test.err:
				t.Errorf("got error: %v, want: %v", err, test.err)
			case !test.expectedErr && result != test.acc.Balance:
				t.Errorf("got: %v, want: %v", result, test.acc.Balance)
			}
		})
	}
}

func TestAccRepo_GetAccount(t *testing.T) {
	accRepo := NewAccountRepo(dbTest)
	ctxDB := context.Background()
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
				ID:        vos.NewUUID(),
				Name:      "Elliot",
				CPF:       "33344455568",
				Secret:    "password",
				Balance:   0,
				CreatedAt: time.Now().Truncate(time.Second),
			},
			insert:      false,
			expectedErr: true,
			err:         entities.ErrAccNotFound,
		},
	}

	defer ClearDB()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.insert {
				_ = accRepo.CreateAccount(ctxDB, test.acc)
			}

			result, err := accRepo.GetAccount(context.Background(), test.acc.CPF)
			if test.expectedErr && err != test.err {
				t.Errorf("got: %v, want: %v", err, test.err)
			}

			if !test.expectedErr && reflect.DeepEqual(&result, test.acc) {
				t.Errorf("got: %v, want: %v", &result, test.acc)
			}
		})
	}
}
