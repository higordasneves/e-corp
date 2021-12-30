package usecase

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"github.com/higordasneves/e-corp/pkg/repository"
	repomock "github.com/higordasneves/e-corp/pkg/repository/mock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTransferUseCase_FetchTransfers(t *testing.T) {
	t.Parallel()

	firstAccountID := vos.NewUUID()
	secondAccountID := vos.NewUUID()
	thirdAccountID := vos.NewUUID()

	accounts := []entities.Account{
		{
			ID:        firstAccountID,
			Name:      "ecorp",
			CPF:       "55566677781",
			Secret:    "password2",
			Balance:   10000,
			CreatedAt: time.Now().Truncate(time.Second),
		},
		{
			ID:        secondAccountID,
			Name:      "Elliot",
			CPF:       "55566677782",
			Secret:    "password123",
			Balance:   2000,
			CreatedAt: time.Now().Truncate(time.Second),
		},
		{
			ID:        thirdAccountID,
			Name:      "penny pincher",
			CPF:       "55566677783",
			Secret:    "password123",
			Balance:   80000000000,
			CreatedAt: time.Now().Truncate(time.Second),
		},
	}

	transfersFirstAcc := []entities.Transfer{
		{
			ID:                   vos.NewUUID(),
			AccountOriginID:      firstAccountID,
			AccountDestinationID: secondAccountID,
			Amount:               10,
			CreatedAt:            time.Now(),
		},
		{
			ID:                   vos.NewUUID(),
			AccountOriginID:      firstAccountID,
			AccountDestinationID: secondAccountID,
			Amount:               3,
			CreatedAt:            time.Now(),
		},
	}

	transfersSecondAcc := []entities.Transfer{
		{
			ID:                   vos.NewUUID(),
			AccountOriginID:      secondAccountID,
			AccountDestinationID: firstAccountID,
			Amount:               43,
			CreatedAt:            time.Now(),
		},
		{
			ID:                   vos.NewUUID(),
			AccountOriginID:      secondAccountID,
			AccountDestinationID: firstAccountID,
			Amount:               86,
			CreatedAt:            time.Now(),
		},
	}

	transfers := append(transfersFirstAcc, transfersSecondAcc...)

	tests := []struct {
		name        string
		id          string
		want        []entities.Transfer
		dbErr       error
		expectedErr error
	}{
		{
			name:        "with success first account",
			id:          firstAccountID.String(),
			want:        transfersFirstAcc,
			dbErr:       nil,
			expectedErr: nil,
		},
		{
			name:        "with success second account",
			id:          secondAccountID.String(),
			want:        transfersSecondAcc,
			dbErr:       nil,
			expectedErr: nil,
		},
		{
			name:        "invalid ID",
			id:          "invalid",
			want:        nil,
			dbErr:       nil,
			expectedErr: vos.ErrInvalidID,
		},
		{
			name:        "account not found",
			id:          vos.NewUUID().String(),
			want:        nil,
			dbErr:       nil,
			expectedErr: entities.ErrAccNotFound,
		},
		{
			name:        "zero transfers",
			id:          thirdAccountID.String(),
			want:        nil,
			dbErr:       nil,
			expectedErr: nil,
		},
		{
			name:        "database error",
			id:          firstAccountID.String(),
			want:        nil,
			dbErr:       repository.ErrUnexpected,
			expectedErr: repository.ErrUnexpected,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// setup
			accRepo := repomock.NewAccountRepo(accounts, nil)
			tRepo := repomock.NewTransferRepo(transfers, tt.dbErr)
			tUseCase := NewTransferUseCase(accRepo, tRepo)

			// execute
			results, err := tUseCase.FetchTransfers(context.Background(), tt.id)

			// assert
			assert.ErrorIs(t, err, tt.expectedErr)
			assert.Equal(t, tt.want, results)
		})
	}

}
