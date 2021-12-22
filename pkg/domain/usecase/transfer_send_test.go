package usecase

import (
	"context"
	"errors"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"github.com/higordasneves/e-corp/pkg/repository"
	repomock "github.com/higordasneves/e-corp/pkg/repository/mock"
	"reflect"
	"testing"
	"time"
)

func TestTransferUseCase_Transfer(t *testing.T) {
	accOriginID := vos.NewUUID()
	accDestinationID := vos.NewUUID()

	accounts := []entities.Account{
		{
			ID:        accOriginID,
			Name:      "ecorp",
			CPF:       "55566677781",
			Secret:    "password2",
			Balance:   10000,
			CreatedAt: time.Now().Truncate(time.Second),
		},
		{
			ID:        accDestinationID,
			Name:      "Elliot",
			CPF:       "55566677782",
			Secret:    "password123",
			Balance:   2000,
			CreatedAt: time.Now().Truncate(time.Second),
		},
	}
	tests := []struct {
		name          string
		transferInput TransferInput
		want          *entities.Transfer
		dbErr         error
		expectedErr   error
	}{
		{
			name: "with success",
			transferInput: TransferInput{
				AccountOriginID:      accOriginID.String(),
				AccountDestinationID: accDestinationID.String(),
				Amount:               1000,
			},
			want: &entities.Transfer{
				AccountOriginID:      accOriginID,
				AccountDestinationID: accDestinationID,
				Amount:               1000,
			},
			dbErr:       nil,
			expectedErr: nil,
		},
		{
			name: "insufficient funds",
			transferInput: TransferInput{
				AccountOriginID:      accOriginID.String(),
				AccountDestinationID: accDestinationID.String(),
				Amount:               10001,
			},
			want:        nil,
			dbErr:       nil,
			expectedErr: entities.ErrTransferInsufficientFunds,
		},
		{
			name: "origin account not found",
			transferInput: TransferInput{
				AccountOriginID:      vos.NewUUID().String(),
				AccountDestinationID: accDestinationID.String(),
				Amount:               10001,
			},
			want:        nil,
			dbErr:       nil,
			expectedErr: entities.ErrAccNotFound,
		},
		{
			name: "destination account not found",
			transferInput: TransferInput{
				AccountOriginID:      accOriginID.String(),
				AccountDestinationID: vos.NewUUID().String(),
				Amount:               10001,
			},
			want:        nil,
			dbErr:       nil,
			expectedErr: entities.ErrAccNotFound,
		},
		{
			name: "invalid origin account id",
			transferInput: TransferInput{
				AccountOriginID:      "invalid",
				AccountDestinationID: vos.NewUUID().String(),
				Amount:               10001,
			},
			want:        nil,
			dbErr:       nil,
			expectedErr: entities.ErrOriginAccID,
		},
		{
			name: "invalid destination account id",
			transferInput: TransferInput{
				AccountOriginID:      accOriginID.String(),
				AccountDestinationID: "invalid",
				Amount:               10001,
			},
			want:        nil,
			dbErr:       nil,
			expectedErr: entities.ErrDestAccID,
		},
		{
			name: "invalid amount",
			transferInput: TransferInput{
				AccountOriginID:      accOriginID.String(),
				AccountDestinationID: accDestinationID.String(),
				Amount:               -1,
			},
			want:        nil,
			dbErr:       nil,
			expectedErr: entities.ErrTransferAmount,
		},
		{
			name: "self transfer",
			transferInput: TransferInput{
				AccountOriginID:      accOriginID.String(),
				AccountDestinationID: accOriginID.String(),
				Amount:               10,
			},
			want:        nil,
			dbErr:       nil,
			expectedErr: entities.ErrSelfTransfer,
		},
		{
			name: "db error",
			transferInput: TransferInput{
				AccountOriginID:      accOriginID.String(),
				AccountDestinationID: accDestinationID.String(),
				Amount:               10,
			},
			want:        nil,
			dbErr:       repository.ErrUnexpected,
			expectedErr: repository.ErrUnexpected,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Helper()
			accRepo := repomock.NewAccountRepo(accounts, nil)
			tRepo := repomock.NewTransferRepo([]entities.Transfer{}, test.dbErr)
			tUseCase := NewTransferUseCase(accRepo, tRepo)

			result, err := tUseCase.Transfer(context.Background(), &test.transferInput)

			if err == nil {
				test.want.ID = result.ID
				test.want.CreatedAt = result.CreatedAt
			}

			switch {
			case !errors.Is(err, test.expectedErr):
				t.Errorf("got error: %v, expected error: %v", err, test.expectedErr)
			case test.expectedErr == nil && !reflect.DeepEqual(result, test.want):
				t.Errorf("got: %v, want: %v", result, test.want)
			}
		})
	}

}
