package postgres

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/domain/models"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"github.com/higordasneves/e-corp/pkg/repository"
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
			name: "with error",
			acc: &models.Account{
				ID:        vos.NewAccID(),
				Name:      "Elliot",
				CPF:       "33344455566",
				Secret:    "password",
				Balance:   0,
				CreatedAt: time.Now().Truncate(time.Second),
			},
			err: repository.ErrCreateAcc,
		},
	}

	for _, test := range tests {
		result := accRepo.CreateAccount(ctxDB, test.acc)
		if result != test.err {
			t.Errorf("got: %v, want: %v", result, test.err)
		}
	}
}

func TestAccRepo_FetchAccounts(t *testing.T) {
	accRepo := NewAccountRepo(dbTest, logTest)
	ctxDB := context.Background()
	accounts := []*models.Account{
		{
			ID:        vos.NewAccID(),
			Name:      "Elliot",
			CPF:       "33344455566",
			Secret:    "password",
			Balance:   5000,
			CreatedAt: time.Now().Truncate(time.Second),
		},

		{
			ID:        vos.NewAccID(),
			Name:      "Mr.Robot",
			CPF:       "33344455567",
			Secret:    "password",
			Balance:   3000,
			CreatedAt: time.Now().Truncate(time.Second),
		},
	}

	for _, acc := range accounts {
		err := accRepo.CreateAccount(ctxDB, acc)
		if err != nil {
			t.Errorf("error inserting accounts")
		}
	}

}
