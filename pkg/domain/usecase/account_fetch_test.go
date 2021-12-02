package usecase

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/domain/models"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"github.com/higordasneves/e-corp/pkg/repository"
	repomock "github.com/higordasneves/e-corp/pkg/repository/mock"
	"github.com/sirupsen/logrus"
	"reflect"
	"testing"
	"time"
)

func TestAccountUseCase_FetchAccounts(t *testing.T) {

	log := logrus.New()
	accountsDate := time.Date(2015, time.June, 24, 23, 59, 0, 0, time.UTC)
	IDs := make([]vos.UUID, 0, 4)

	for i := 0; i < cap(IDs); i++ {
		IDs = append(IDs, vos.NewUUID())
	}

	accounts := []models.Account{
		{
			ID:        IDs[0],
			Name:      "Elliot",
			CPF:       "55566677780",
			Secret:    "password1",
			Balance:   9700000,
			CreatedAt: accountsDate,
		},
		{
			ID:        IDs[1],
			Name:      "Elliot",
			CPF:       "55566677781",
			Secret:    "password2",
			Balance:   5596400,
			CreatedAt: accountsDate,
		},
		{
			ID:        IDs[2],
			Name:      "Elliot",
			CPF:       "55566677782",
			Secret:    "password3",
			Balance:   5534513,
			CreatedAt: accountsDate,
		},
		{
			ID:        IDs[3],
			Name:      "Elliot",
			CPF:       "55566677783",
			Secret:    "password4",
			Balance:   12350,
			CreatedAt: accountsDate,
		},
	}

	want := []models.AccountOutput{
		{
			ID:        IDs[0],
			Name:      "Elliot",
			CPF:       "555.666.777-80",
			Balance:   97000.00,
			CreatedAt: accountsDate,
		},
		{
			ID:        IDs[1],
			Name:      "Elliot",
			CPF:       "555.666.777-81",
			Balance:   55964.00,
			CreatedAt: accountsDate,
		},
		{
			ID:        IDs[2],
			Name:      "Elliot",
			CPF:       "555.666.777-82",
			Balance:   55345.13,
			CreatedAt: accountsDate,
		},
		{
			ID:        IDs[3],
			Name:      "Elliot",
			CPF:       "555.666.777-83",
			Balance:   123.50,
			CreatedAt: accountsDate,
		},
	}

	t.Run("with success", func(t *testing.T) {
		accRepo := repomock.NewAccountRepo(accounts, nil)
		accUseCase := NewAccountUseCase(accRepo, log)
		result, err := accUseCase.FetchAccounts(context.Background())
		if err != nil {
			t.Errorf("didn't want an error, but got the error: %v", err)
		}
		if !reflect.DeepEqual(result, want) {
			t.Errorf("got %v, \n want %v", result, want)
		}
	})

	t.Run("expect database error", func(t *testing.T) {
		accRepo := repomock.NewAccountRepo(accounts, repository.ErrFetchAcc)
		accUseCase := NewAccountUseCase(accRepo, log)
		_, err := accUseCase.FetchAccounts(context.Background())
		if err == nil {
			t.Error("wanted an error but didn't get one")
		}
	})
}
