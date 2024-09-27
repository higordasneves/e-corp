package postgres

import (
	"context"
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/higordasneves/e-corp/pkg/domain/entities"
)

func TestTransferRepo_CreateTransfer(t *testing.T) {
	t.Parallel()

	// setup
	accOriginID := uuid.Must(uuid.NewV7())
	accDestinationID := uuid.Must(uuid.NewV7())

	accounts := []entities.Account{
		{
			ID:        accOriginID,
			Name:      "Elliot",
			Document:  "33344455567",
			Secret:    "password",
			Balance:   7000,
			CreatedAt: time.Now().Truncate(time.Second),
		},
		{
			ID:        accDestinationID,
			Name:      "Mr.Robot",
			Document:  "33344455568",
			Secret:    "password",
			Balance:   3000,
			CreatedAt: time.Now().Truncate(time.Second),
		},
	}

	r := NewRepository(NewDB(t))
	ctxDB := context.Background()

	for _, acc := range accounts {
		err := r.CreateAccount(ctxDB, acc)
		require.NoError(t, err)
	}

	tests := []struct {
		name        string
		transfer    entities.Transfer
		expectedErr error
	}{
		{
			name: "success",
			transfer: entities.Transfer{
				ID:                   uuid.Must(uuid.NewV7()),
				AccountOriginID:      accOriginID,
				AccountDestinationID: accDestinationID,
				Amount:               rand.Int(),
				CreatedAt:            time.Now().Truncate(time.Second),
			},
			expectedErr: nil,
		},
		{
			name: "violates foreign key constraint",
			transfer: entities.Transfer{
				ID:                   uuid.FromStringOrNil("5f2d4920-89c3-4ed5-af8e-1d411588746d"),
				AccountOriginID:      uuid.Must(uuid.NewV7()),
				AccountDestinationID: uuid.Must(uuid.NewV7()),
				Amount:               rand.Int(),
				CreatedAt:            time.Now().Truncate(time.Second),
			},
			expectedErr: errors.New("inserting transfer with id 5f2d4920-89c3-4ed5-af8e-1d411588746d"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// execute
			err := r.CreateTransfer(context.Background(), tt.transfer)
			if tt.expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.expectedErr.Error())
			}
		})
	}
}

func TestTransferRepo_ListAccountTransfers(t *testing.T) {
	t.Parallel()

	// setup
	accOriginID := uuid.Must(uuid.NewV7())
	accDestinationID := uuid.Must(uuid.NewV7())

	accounts := []entities.Account{
		{
			ID:        accOriginID,
			Name:      "Elliot",
			Document:  "33344455567",
			Secret:    "password",
			Balance:   999999999999,
			CreatedAt: time.Now().Truncate(time.Second),
		},
		{
			ID:        accDestinationID,
			Name:      "Mr.Robot",
			Document:  "33344455568",
			Secret:    "password",
			Balance:   999999999999,
			CreatedAt: time.Now().Truncate(time.Second),
		},
	}

	r := NewRepository(NewDB(t))
	ctxDB := context.Background()
	for _, acc := range accounts {
		err := r.CreateAccount(ctxDB, acc)
		require.NoError(t, err)
	}

	var want []entities.Transfer
	for i := 0; i < 1000; i++ {
		transfer := entities.Transfer{
			ID:                   uuid.Must(uuid.NewV7()),
			AccountOriginID:      accOriginID,
			AccountDestinationID: accDestinationID,
			Amount:               rand.Intn(100),
			CreatedAt:            time.Now().Truncate(time.Second),
		}
		err := r.CreateTransfer(ctxDB, transfer)
		if err != nil {
			t.Fatal("error inserting transfers")
		}
		want = append(want, transfer)
	}

	t.Run("listing the transfers sent by an account", func(t *testing.T) {
		// execute
		result, err := r.ListAccountTransfers(context.Background(), accOriginID)
		// assert
		require.NoError(t, err)
		assert.ElementsMatch(t, want, result)
	})

	t.Run("listing the transfers received by an account", func(t *testing.T) {
		// execute
		result, err := r.ListAccountTransfers(context.Background(), accDestinationID)
		// assert
		require.NoError(t, err)
		assert.ElementsMatch(t, want, result)
	})
}
