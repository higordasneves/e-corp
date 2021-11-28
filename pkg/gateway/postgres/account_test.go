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
				ID:        vos.NewAccID(),
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
				ID:        vos.NewAccID(),
				Name:      "Elliot",
				CPF:       "33344455566",
				Secret:    "password",
				Balance:   0,
				CreatedAt: time.Now().Truncate(time.Second),
			},
			err: domainerr.ErrAccAlreadyExists,
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
			ID:        vos.NewAccID(),
			Name:      "Elliot",
			CPF:       "33344455567",
			Secret:    "password",
			Balance:   7000,
			CreatedAt: time.Now().Truncate(time.Second),
		},

		{
			ID:        vos.NewAccID(),
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
		t.Error(repository.ErrFetchAcc)
	}

	if !reflect.DeepEqual(want, result) {
		t.Errorf("got: %v, want: %v", result, want)
	}
}

func TestAccRepo_GetBalance(t *testing.T) {
	accRepo := NewAccountRepo(dbTest, logTest)
	ctxDB := context.Background()

	accounts := []models.Account{
		{
			ID:        vos.NewAccID(),
			Name:      "Elliot",
			CPF:       "33344455567",
			Secret:    "password",
			Balance:   7000,
			CreatedAt: time.Now().Truncate(time.Second),
		},

		{
			ID:        vos.NewAccID(),
			Name:      "Mr.Robot",
			CPF:       "33344455568",
			Secret:    "password",
			Balance:   3000,
			CreatedAt: time.Now().Truncate(time.Second),
		},
	}

	for _, acc := range accounts {
		err := accRepo.CreateAccount(ctxDB, &acc)
		if err != nil {
			t.Error("error inserting accounts")
		}
	}

	defer ClearDB()
	for _, acc := range accounts {
		result, err := accRepo.GetBalance(context.Background(), acc.ID)
		if err != nil {
			t.Error(repository.ErrGetBalance)
		}
		if *result != acc.Balance {
			t.Errorf("got: %v, want: %v", *result, acc.Balance)
		}
	}
}
