package usecase

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"github.com/higordasneves/e-corp/pkg/repository"
	repomock "github.com/higordasneves/e-corp/pkg/repository/mock"
	"github.com/sirupsen/logrus"

	"testing"
)

func TestAccountUseCase_GetBalance(t *testing.T) {
	log := logrus.New()

	accInfo := make(map[int]vos.UUID, 3)
	accInfo[162000] = vos.NewUUID()
	accInfo[561300] = vos.NewUUID()

	accounts := make([]entities.Account, 0, 3)
	for i, v := range accInfo {
		accounts = append(accounts, entities.Account{ID: v, Balance: i})
	}

	tests := []struct {
		name        string
		id          vos.UUID
		want        int
		expectedErr error
		dbErr       error
	}{
		{
			name:        "with success 1",
			id:          accInfo[162000],
			want:        162000,
			expectedErr: nil,
			dbErr:       nil,
		},
		{
			name:        "with success 2",
			id:          accInfo[561300],
			want:        561300,
			expectedErr: nil,
			dbErr:       nil,
		},
		{
			name:        "err account not found",
			id:          vos.NewUUID(),
			want:        0,
			expectedErr: entities.ErrAccNotFound,
			dbErr:       nil,
		},
		{
			name:        "database generic error",
			id:          accInfo[561300],
			want:        561300,
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

			if err == nil && balance != test.want {
				t.Errorf("got %v, want %v", balance, test.want)
			}
		})
	}
}
