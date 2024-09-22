package postgres

import (
	"context"
	"errors"
	"github.com/gofrs/uuid/v5"
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
	"time"

	thelp "github.com/higordasneves/e-corp/extensions/testhelpers"
	"github.com/higordasneves/e-corp/pkg/domain"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
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
				ID:        vos.NewUUID(),
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
				ID:        vos.NewUUID(),
				Name:      "Elliot",
				CPF:       "33344455566",
				Secret:    "password",
				Balance:   0,
				CreatedAt: time.Now().Truncate(time.Second),
			},
			wantErr: domain.Error(domain.InvalidParamErrorType, "account already exists", nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accRepo := NewRepository(db)
			err := accRepo.CreateAccount(context.Background(), tt.input)
			if tt.wantErr != nil {
				thelp.AssertDomainError(t, tt.wantErr, err)
			}
		})
	}
}

func TestAccRepo_ListAccounts_Success(t *testing.T) {
	t.Parallel()

	repo := NewRepository(NewDB(t))

	accounts := []entities.Account{
		{
			ID:        vos.NewUUID(),
			Name:      "Elliot",
			CPF:       "33344455567",
			Secret:    "password",
			Balance:   7000,
			CreatedAt: time.Now().Truncate(time.Second),
		},
		{
			ID:        vos.NewUUID(),
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
			ID:        vos.NewUUID(),
			Name:      "Elliot",
			CPF:       "33344455567",
			Secret:    "password",
			Balance:   7000,
			CreatedAt: time.Now().Truncate(time.Second),
		},
		{
			ID:        vos.NewUUID(),
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

	accRepo := NewRepository(NewDB(t))

	tests := []struct {
		name        string
		acc         entities.Account
		insert      bool
		expectedErr bool
		err         error
	}{
		{
			name: "with success",
			acc: entities.Account{
				ID:        vos.NewUUID(),
				Name:      "Elliot",
				CPF:       "33344455567",
				Secret:    "password",
				Balance:   7000,
				CreatedAt: time.Now().Truncate(time.Second),
			},
			insert:      true,
			expectedErr: false,
			err:         nil,
		},
		{
			name: "with success balance 0",
			acc: entities.Account{
				ID:        vos.NewUUID(),
				Name:      "Elliot",
				CPF:       "33344455568",
				Secret:    "password",
				Balance:   0,
				CreatedAt: time.Now().Truncate(time.Second),
			},
			insert:      true,
			expectedErr: false,
			err:         nil,
		},
		{
			name: "Repository not found",
			acc: entities.Account{
				ID: vos.NewUUID(),
			},
			insert:      false,
			expectedErr: true,
			err:         entities.ErrAccNotFound,
		},
	}

	var GotDBError *domain.DBError
	var WantDBError *domain.DBError

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.insert {
				_ = accRepo.CreateAccount(context.Background(), tt.acc)
			}

			// execute
			result, err := accRepo.GetBalance(context.Background(), tt.acc.ID)

			// assert
			switch {
			case errors.As(err, &GotDBError) && errors.As(tt.err, &WantDBError):
				if GotDBError.Query != WantDBError.Query {
					t.Errorf("got sql error in query: %v, want: %v", GotDBError.Query, WantDBError.Query)
				}
			case !errors.Is(err, tt.err):
				t.Errorf("got error: %v, want: %v", err, tt.err)
			case !tt.expectedErr && result != tt.acc.Balance:
				t.Errorf("got: %v, want: %v", result, tt.acc.Balance)
			}
		})
	}
}

func TestAccRepo_UpdateBalance(t *testing.T) {
	tests := []struct {
		name         string
		acc          entities.Account
		updateAmount int
		insert       bool
		expectedErr  bool
		err          error
	}{
		{
			name: "with success outbound",
			acc: entities.Account{
				ID:        vos.NewUUID(),
				Name:      "Elliot",
				CPF:       "33344455567",
				Secret:    "password",
				Balance:   7000,
				CreatedAt: time.Now().Truncate(time.Second),
			},
			updateAmount: -5000,
			insert:       true,
			expectedErr:  false,
			err:          nil,
		},
		{
			name: "with success inbound",
			acc: entities.Account{
				ID:        vos.NewUUID(),
				Name:      "Elliot",
				CPF:       "33344455568",
				Secret:    "password",
				Balance:   0,
				CreatedAt: time.Now().Truncate(time.Second),
			},
			updateAmount: 5000,
			insert:       true,
			expectedErr:  false,
			err:          nil,
		},
	}

	var GotDBError *domain.DBError
	var WantDBError *domain.DBError

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			accRepo := NewRepository(NewDB(t))
			if tt.insert {
				_ = accRepo.CreateAccount(context.Background(), tt.acc)
			}

			// execute
			err := accRepo.UpdateBalance(context.Background(), tt.acc.ID, tt.updateAmount)

			// assert
			switch {
			case errors.As(err, &GotDBError) && errors.As(tt.err, &WantDBError):
				if GotDBError.Query != WantDBError.Query {
					t.Errorf("got sql error in query: %v, want: %v", GotDBError.Query, WantDBError.Query)
				}
			case err != tt.err:
				t.Errorf("got error: %v, want: %v", err, tt.err)
			case !tt.expectedErr:
				gotBalance, errGetBalance := accRepo.GetBalance(context.Background(), tt.acc.ID)
				if errGetBalance != nil {
					t.Error("unexpected error in get balance query")
				} else if gotBalance != tt.acc.Balance+tt.updateAmount {
					t.Errorf("got: %v, want: %v", gotBalance, tt.acc.Balance+tt.updateAmount)
				}
			}
		})
	}
}

func TestAccRepo_GetAccount(t *testing.T) {
	tests := []struct {
		name        string
		acc         entities.Account
		insert      bool
		expectedErr bool
		err         error
	}{
		{
			name: "with success",
			acc: entities.Account{
				ID:        vos.NewUUID(),
				Name:      "Elliot",
				CPF:       "33344455567",
				Secret:    "password",
				Balance:   7000,
				CreatedAt: time.Now().Truncate(time.Second),
			},
			insert:      true,
			expectedErr: false,
			err:         nil,
		},
		{
			name: "Repository not found",
			acc: entities.Account{
				ID: vos.NewUUID(),
			},
			insert:      false,
			expectedErr: true,
			err:         entities.ErrAccNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			accRepo := NewRepository(NewDB(t))
			if tt.insert {
				_ = accRepo.CreateAccount(context.Background(), tt.acc)
			}

			// execute
			result, err := accRepo.GetAccount(context.Background(), tt.acc.CPF)
			if tt.expectedErr && err != tt.err {
				t.Errorf("got: %v, want: %v", err, tt.err)
			}

			// assert
			if !tt.expectedErr && reflect.DeepEqual(&result, tt.acc) {
				t.Errorf("got: %v, want: %v", &result, tt.acc)
			}
		})
	}
}
