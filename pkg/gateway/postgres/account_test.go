package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/higordasneves/e-corp/pkg/domain"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
)

func TestAccRepo_CreateAccount(t *testing.T) {
	t.Parallel()
	db := NewDB(t)

	tests := []struct {
		name    string
		input   entities.Account
		wantErr error
	}{
		{
			name: "success",
			input: entities.Account{
				ID:        uuid.Must(uuid.NewV7()),
				Name:      "Elliot",
				CPF:       "33344455566",
				Secret:    "password",
				Balance:   0,
				CreatedAt: time.Now().Truncate(time.Second),
			},
			wantErr: nil,
		},
		{
			name: "fail - account already exists",
			input: entities.Account{
				ID:        uuid.Must(uuid.NewV7()),
				Name:      "Elliot",
				CPF:       "33344455566",
				Secret:    "password",
				Balance:   0,
				CreatedAt: time.Now().Truncate(time.Second),
			},
			wantErr: domain.ErrInvalidParameter,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accRepo := NewRepository(db)
			err := accRepo.CreateAccount(context.Background(), tt.input)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestAccRepo_ListAccounts_Success(t *testing.T) {
	t.Parallel()

	repo := NewRepository(NewDB(t))

	accounts := []entities.Account{
		{
			ID:        uuid.Must(uuid.NewV7()),
			Name:      "Elliot",
			CPF:       "33344455567",
			Secret:    "password",
			Balance:   7000,
			CreatedAt: time.Now().Truncate(time.Second),
		},
		{
			ID:        uuid.Must(uuid.NewV7()),
			Name:      "Mr.Robot",
			CPF:       "33344455568",
			Secret:    "password",
			Balance:   3000,
			CreatedAt: time.Now().Truncate(time.Second),
		},
	}
	for _, acc := range accounts {
		err := repo.CreateAccount(context.Background(), acc)
		if err != nil {
			t.Error("error inserting accounts")
		}
	}

	// execute
	result, err := repo.ListAccounts(context.Background(), usecase.ListAccountsInput{
		IDs:           []uuid.UUID{uuid.FromStringOrNil(accounts[0].ID.String()), uuid.FromStringOrNil(accounts[1].ID.String())},
		LastFetchedID: uuid.UUID{},
		PageSize:      2,
	})
	require.NoError(t, err)

	//assert
	assert.ElementsMatch(t, accounts, result.Accounts)
}

// The purpose of this test is to verify the pagination functionality.
func TestAccRepo_ListAccounts_Success_Pagination(t *testing.T) {
	t.Parallel()

	repo := NewRepository(NewDB(t))

	accounts := []entities.Account{
		{
			ID:        uuid.Must(uuid.NewV7()),
			Name:      "Elliot",
			CPF:       "33344455567",
			Secret:    "password",
			Balance:   7000,
			CreatedAt: time.Now().Truncate(time.Second),
		},
		{
			ID:        uuid.Must(uuid.NewV7()),
			Name:      "Mr.Robot",
			CPF:       "33344455568",
			Secret:    "password",
			Balance:   3000,
			CreatedAt: time.Now().Truncate(time.Second),
		},
	}
	for _, acc := range accounts {
		err := repo.CreateAccount(context.Background(), acc)
		if err != nil {
			t.Error("error inserting accounts")
		}
	}

	// execute: listing page 1
	result, err := repo.ListAccounts(context.Background(), usecase.ListAccountsInput{
		IDs:           []uuid.UUID{uuid.FromStringOrNil(accounts[0].ID.String()), uuid.FromStringOrNil(accounts[1].ID.String())},
		LastFetchedID: uuid.UUID{},
		PageSize:      1,
	})
	require.NoError(t, err)

	// asserting page 1
	assert.Len(t, result.Accounts, 1)
	assert.Equal(t, accounts[1], result.Accounts[0])
	assert.Equal(t, usecase.ListAccountsInput{
		IDs:           []uuid.UUID{uuid.FromStringOrNil(accounts[0].ID.String()), uuid.FromStringOrNil(accounts[1].ID.String())},
		LastFetchedID: uuid.FromStringOrNil(accounts[1].ID.String()),
		PageSize:      1,
	}, *result.NextPage)

	// execute: listing page 2
	result, err = repo.ListAccounts(context.Background(), *result.NextPage)
	require.NoError(t, err)

	// asserting page 2
	assert.Len(t, result.Accounts, 1)
	assert.Equal(t, accounts[0], result.Accounts[0])
	assert.Nil(t, result.NextPage)
}

func TestAccRepo_GetBalance(t *testing.T) {
	t.Parallel()

	r := NewRepository(NewDB(t))
	account := entities.Account{
		ID:        uuid.Must(uuid.NewV7()),
		Name:      "Elliot",
		CPF:       "33344455567",
		Secret:    "password",
		Balance:   7000,
		CreatedAt: time.Now().Truncate(time.Second),
	}
	require.NoError(t, r.CreateAccount(context.Background(), account))

	tests := []struct {
		name    string
		input   uuid.UUID
		want    int
		wantErr error
	}{
		{
			name:    "with success",
			input:   account.ID,
			want:    account.Balance,
			wantErr: nil,
		},
		{
			name:    "account not found",
			want:    0,
			wantErr: domain.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// execute
			result, err := r.GetBalance(context.Background(), tt.input)

			// assert
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, result)

		})
	}
}

func TestAccRepo_UpdateBalance(t *testing.T) {
	t.Parallel()

	// setup
	r := NewRepository(NewDB(t))
	account1 := entities.Account{
		ID:        uuid.Must(uuid.NewV7()),
		Name:      "Elliot",
		CPF:       "33344455567",
		Secret:    "password",
		Balance:   7000,
		CreatedAt: time.Now().Truncate(time.Second),
	}
	require.NoError(t, r.CreateAccount(context.Background(), account1))

	account2 := entities.Account{
		ID:        uuid.Must(uuid.NewV7()),
		Name:      "Elliot",
		CPF:       "33344455568",
		Secret:    "password",
		Balance:   0,
		CreatedAt: time.Now().Truncate(time.Second),
	}
	require.NoError(t, r.CreateAccount(context.Background(), account2))

	tests := []struct {
		name               string
		accountID          uuid.UUID
		amount             int
		wantAccountBalance int
	}{
		{
			name:               "success - positive amount",
			accountID:          account1.ID,
			amount:             -5000,
			wantAccountBalance: 2000,
		},
		{
			name:               "success - negative amount",
			accountID:          account2.ID,
			amount:             5000,
			wantAccountBalance: 5000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// execute
			err := r.UpdateBalance(context.Background(), tt.accountID, tt.amount)

			// assert
			require.NoError(t, err)
			got, err := r.GetBalance(context.Background(), tt.accountID)
			require.NoError(t, err)
			assert.Equal(t, tt.wantAccountBalance, got)
		})
	}
}

func TestAccRepo_GetAccountByDocument(t *testing.T) {
	t.Parallel()

	// setup
	r := NewRepository(NewDB(t))
	account1 := entities.Account{
		ID:        uuid.Must(uuid.NewV7()),
		Name:      "Elliot",
		CPF:       "33344455567",
		Secret:    "password",
		Balance:   7000,
		CreatedAt: time.Now().Truncate(time.Second),
	}
	require.NoError(t, r.CreateAccount(context.Background(), account1))

	tests := []struct {
		name     string
		document vos.CPF
		want     entities.Account
		wantErr  error
	}{
		{
			name:     "with success",
			document: vos.CPF("33344455567"),
			want:     account1,
			wantErr:  nil,
		},
		{
			name:     "Repository not found",
			document: vos.CPF("1"),
			wantErr:  domain.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// execute
			result, err := r.GetAccountByDocument(context.Background(), tt.document)
			// assert
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, result)
		})
	}
}
