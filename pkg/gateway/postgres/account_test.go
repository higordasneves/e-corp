package postgres

import (
	"context"
	domainerr "github.com/higordasneves/e-corp/pkg/domain/errors"
	"github.com/higordasneves/e-corp/pkg/domain/models"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"github.com/higordasneves/e-corp/pkg/repository"
	"reflect"
	"testing"
	"time"
)

func TestAccRepo_CreateAccount(t *testing.T) {
	accRepo := NewAccountRepo(dbTest, logTest)
	ctxDB := context.Background()
	tests := []struct {
		name string
		acc  *models.Account
		err  error
	}{
		{
			name: "with success",
			acc: &models.Account{
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
			acc: &models.Account{
				ID:        vos.NewUUID(),
				Name:      "Elliot",
				CPF:       "33344455566",
				Secret:    "password",
				Balance:   0,
				CreatedAt: time.Now().Truncate(time.Second),
			},
			err: domainerr.ErrAccAlreadyExists,
		},
		{
			name: "invalid id",
			acc: &models.Account{
				ID:        "invalid",
				Name:      "Elliot",
				CPF:       "33344455567",
				Secret:    "password",
				Balance:   0,
				CreatedAt: time.Now().Truncate(time.Second),
			},
			err: repository.ErrCreateAcc,
		},
	}
	defer ClearDB()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := accRepo.CreateAccount(ctxDB, test.acc)
			if result != test.err {
				t.Errorf("got: %v, want: %v", result, test.err)
			}
		})
	}
}

func TestAccRepo_FetchAccounts(t *testing.T) {
	accRepo := NewAccountRepo(dbTest, logTest)
	ctxDB := context.Background()

	accounts := []models.Account{
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
	var want []models.Account
	for _, acc := range accounts {
		err := accRepo.CreateAccount(ctxDB, &acc)
		if err != nil {
			t.Error("error inserting accounts")
		}
		want = append(want, models.Account{
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
		t.Errorf("got: %v, want: %v", err, nil)
	}

	if !reflect.DeepEqual(want, result) {
		t.Errorf("got: %v, want: %v", result, want)
	}
}

func TestAccRepo_GetBalance(t *testing.T) {
	accRepo := NewAccountRepo(dbTest, logTest)
	ctxDB := context.Background()
	tests := []struct {
		name        string
		acc         *models.Account
		insert      bool
		expectedErr bool
		err         error
	}{
		{
			name: "with success",
			acc: &models.Account{
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
			acc: &models.Account{
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
			acc: &models.Account{
				ID:        vos.NewUUID(),
				Name:      "Elliot",
				CPF:       "33344455569",
				Secret:    "password",
				Balance:   0,
				CreatedAt: time.Now().Truncate(time.Second),
			},
			insert:      false,
			expectedErr: true,
			err:         domainerr.ErrAccNotFound,
		},
		{
			name: "invalid id",
			acc: &models.Account{
				ID:        "invalid",
				Name:      "Elliot",
				CPF:       "33344455570",
				Secret:    "password",
				Balance:   0,
				CreatedAt: time.Now().Truncate(time.Second),
			},
			insert:      false,
			expectedErr: true,
			err:         repository.ErrGetBalance,
		},
	}

	defer ClearDB()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.insert {
				_ = accRepo.CreateAccount(ctxDB, test.acc)
			}

			result, err := accRepo.GetBalance(context.Background(), test.acc.ID)
			if test.expectedErr && err != test.err {
				t.Errorf("got: %v, want: %v", err, test.err)
			}

			if !test.expectedErr && *result != test.acc.Balance {
				t.Errorf("got: %v, want: %v", *result, test.acc.Balance)
			}
		})
	}
}

func TestAccRepo_GetAccount(t *testing.T) {
	accRepo := NewAccountRepo(dbTest, logTest)
	ctxDB := context.Background()
	tests := []struct {
		name        string
		acc         *models.Account
		insert      bool
		expectedErr bool
		err         error
	}{
		{
			name: "with success",
			acc: &models.Account{
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
			acc: &models.Account{
				ID:        vos.NewUUID(),
				Name:      "Elliot",
				CPF:       "33344455568",
				Secret:    "password",
				Balance:   0,
				CreatedAt: time.Now().Truncate(time.Second),
			},
			insert:      false,
			expectedErr: true,
			err:         domainerr.ErrAccNotFound,
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
