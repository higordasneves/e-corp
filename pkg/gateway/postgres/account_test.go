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
