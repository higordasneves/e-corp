package usecase

import (
	"context"
	domainerr "github.com/higordasneves/e-corp/pkg/domain/errors"
	"github.com/higordasneves/e-corp/pkg/domain/models"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"github.com/higordasneves/e-corp/pkg/repository"
	repomock "github.com/higordasneves/e-corp/pkg/repository/mock"
	"github.com/sirupsen/logrus"

	"testing"
)

func TestAccountUseCase_GetBalance(t *testing.T) {
	log := logrus.New()

	accInfo := make(map[vos.Currency]vos.UUID, 3)
	accInfo[162000] = vos.NewUUID()
	accInfo[561300] = vos.NewUUID()

	accounts := make([]models.Account, 0, 3)
	for i, v := range accInfo {
		accounts = append(accounts, models.Account{ID: v, Balance: i})
	}

	tests := []struct {
		name        string
		id          vos.UUID
		want        vos.Currency
		expectedErr error
		dbErr       error
	}{
		{
			name:        "with success 1",
			id:          accInfo[162000],
			want:        vos.Currency(1620),
			expectedErr: nil,
			dbErr:       nil,
		},
		{
			name:        "with success 2",
			id:          accInfo[561300],
			want:        vos.Currency(5613),
			expectedErr: nil,
			dbErr:       nil,
		},
		{
			name:        "err account not found",
			id:          vos.NewUUID(),
			want:        0,
			expectedErr: domainerr.ErrAccNotFound,
			dbErr:       nil,
		},
		{
			name:        "database generic error",
			id:          accInfo[561300],
			want:        vos.Currency(5613),
			expectedErr: repository.ErrGetBalance,
			dbErr:       repository.ErrGetBalance,
		},
	}

	for _, test := range tests {
		accRepo := repomock.NewAccountRepo(accounts, test.dbErr)
		accUseCase := NewAccountUseCase(accRepo, log)
		balance, err := accUseCase.GetBalance(context.Background(), test.id)

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

			if err == nil && *balance != test.want {
				t.Errorf("got %v, want %v", *balance, test.want)
			}
		})
	}
}
