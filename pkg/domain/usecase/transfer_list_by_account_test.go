package usecase_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/higordasneves/e-corp/pkg/domain"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/gateway/postgres"
)

func TestTransferUseCase_ListAccountTransfers_Success(t *testing.T) {
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

	r := postgres.NewRepository(NewDB(t))
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

	uc := usecase.TransferUseCase{R: r}

	t.Run("listing the transfers sent by an account", func(t *testing.T) {
		// execute
		result, err := uc.ListAccountTransfers(context.Background(), usecase.ListAccountTransfersInput{AccountID: accOriginID})
		// assert
		require.NoError(t, err)
		assert.ElementsMatch(t, want, result.Transfers)
	})

	t.Run("listing the transfers received by an account", func(t *testing.T) {
		// execute
		result, err := uc.ListAccountTransfers(context.Background(), usecase.ListAccountTransfersInput{AccountID: accDestinationID})
		// assert
		require.NoError(t, err)
		assert.ElementsMatch(t, want, result.Transfers)
	})
}

func TestTransferUseCase_ListAccountTransfers_Failure_NotFound(t *testing.T) {
	t.Parallel()

	// setup
	r := postgres.NewRepository(NewDB(t))
	uc := usecase.TransferUseCase{R: r}
	//execute
	_, err := uc.ListAccountTransfers(context.Background(), usecase.ListAccountTransfersInput{AccountID: uuid.Must(uuid.NewV7())})
	// assert
	require.ErrorIs(t, err, domain.ErrNotFound)
}
