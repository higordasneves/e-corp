package usecase_test

import (
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/higordasneves/e-corp/pkg/domain"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/gateway/postgres"
	thelp "github.com/higordasneves/e-corp/utils/testhelpers"
)

func TestTransferUC_Transfer_Success(t *testing.T) {
	t.Parallel()

	// setup
	r := postgres.NewRepository(NewDB(t))
	uc := usecase.TransferUC{R: r}

	accOriginID := uuid.Must(uuid.NewV7())
	accDestinationID := uuid.Must(uuid.NewV7())

	accounts := []entities.Account{
		{
			ID:        accOriginID,
			Name:      "Elliot",
			Document:  "33344455567",
			Secret:    "password",
			Balance:   15,
			CreatedAt: time.Now().Truncate(time.Second),
		},
		{
			ID:        accDestinationID,
			Name:      "Mr.Robot",
			Document:  "33344455568",
			Secret:    "password",
			Balance:   2,
			CreatedAt: time.Now().Truncate(time.Second),
		},
	}

	ctx := thelp.NewCtx(t)
	for _, acc := range accounts {
		err := r.CreateAccount(ctx, acc)
		require.NoError(t, err)
	}

	// execute
	got, err := uc.Transfer(ctx, usecase.TransferInput{
		AccountOriginID:      accOriginID,
		AccountDestinationID: accDestinationID,
		Amount:               5,
	})
	require.NoError(t, err)

	// assert
	assert.Equal(t, accOriginID, got.Transfer.AccountOriginID)
	assert.Equal(t, accDestinationID, got.Transfer.AccountDestinationID)
	assert.Equal(t, 5, got.Transfer.Amount)

	// asserting accounts balance
	accOriginAfterBalance, err := r.GetBalance(ctx, accOriginID)
	require.NoError(t, err)
	assert.Equal(t, 10, accOriginAfterBalance)

	accDestAfterBalance, err := r.GetBalance(ctx, accDestinationID)
	require.NoError(t, err)
	assert.Equal(t, 7, accDestAfterBalance)

	// asserting accounts transfers
	accOriginTransfers, err := r.ListAccountTransfers(ctx, accOriginID)
	require.NoError(t, err)
	assert.Len(t, accOriginTransfers, 1)
	assert.Equal(t, got.Transfer, accOriginTransfers[0])

	accDestTransfers, err := r.ListAccountTransfers(ctx, accDestinationID)
	require.NoError(t, err)
	assert.Len(t, accDestTransfers, 1)
	assert.Equal(t, got.Transfer, accDestTransfers[0])
}

func TestTransferUC_Transfer(t *testing.T) {
	t.Parallel()

	// setup
	r := postgres.NewRepository(NewDB(t))
	uc := usecase.TransferUC{R: r}

	accOriginID := uuid.Must(uuid.NewV7())
	accDestinationID := uuid.Must(uuid.NewV7())

	accounts := []entities.Account{
		{
			ID:        accOriginID,
			Name:      "Elliot",
			Document:  "33344455567",
			Secret:    "password",
			Balance:   15,
			CreatedAt: time.Now().Truncate(time.Second),
		},
		{
			ID:        accDestinationID,
			Name:      "Mr.Robot",
			Document:  "33344455568",
			Secret:    "password",
			Balance:   2,
			CreatedAt: time.Now().Truncate(time.Second),
		},
	}

	ctx := thelp.NewCtx(t)
	for _, acc := range accounts {
		err := r.CreateAccount(ctx, acc)
		require.NoError(t, err)
	}

	tests := []struct {
		name       string
		input      usecase.TransferInput
		wantErr    error
		wantErrMsg string
	}{
		{
			name: "origin and destination account id is the same",
			input: usecase.TransferInput{
				AccountOriginID:      accOriginID,
				AccountDestinationID: accOriginID,
				Amount:               5,
			},
			wantErr:    domain.ErrInvalidParameter,
			wantErrMsg: "the destination account must be different from the origin account",
		},
		{
			name: "invalid transfer amount",
			input: usecase.TransferInput{
				AccountOriginID:      accOriginID,
				AccountDestinationID: accDestinationID,
				Amount:               0,
			},
			wantErr:    domain.ErrInvalidParameter,
			wantErrMsg: "invalid transfer amount, the amount must be greater than 0",
		},
		{
			name: "origin account not fount",
			input: usecase.TransferInput{
				AccountOriginID:      uuid.Must(uuid.NewV7()),
				AccountDestinationID: accDestinationID,
				Amount:               5,
			},
			wantErr:    domain.ErrNotFound,
			wantErrMsg: "getting origin account balance",
		},
		{
			name: "destination account not fount",
			input: usecase.TransferInput{
				AccountOriginID:      accOriginID,
				AccountDestinationID: uuid.Must(uuid.NewV7()),
				Amount:               5,
			},
			wantErr:    domain.ErrNotFound,
			wantErrMsg: "getting destination account balance",
		},
		{
			name: "not enough funds",
			input: usecase.TransferInput{
				AccountOriginID:      accOriginID,
				AccountDestinationID: accDestinationID,
				Amount:               16,
			},
			wantErr:    domain.ErrInvalidParameter,
			wantErrMsg: "insufficient funds",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := uc.Transfer(thelp.NewCtx(t), tt.input)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.ErrorContains(t, err, tt.wantErrMsg)

			// asserting that accounts balance doesn't change.
			accOriginAfterBalance, err := r.GetBalance(ctx, accOriginID)
			require.NoError(t, err)
			assert.Equal(t, 15, accOriginAfterBalance)

			accDestAfterBalance, err := r.GetBalance(ctx, accDestinationID)
			require.NoError(t, err)
			assert.Equal(t, 2, accDestAfterBalance)
		})
	}
}
