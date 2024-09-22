package postgres

import (
	"context"
	"errors"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/higordasneves/e-corp/pkg/domain"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
)

func TestTransferRepo_CreateTransfer(t *testing.T) {
	// setup
	accOriginID := vos.NewUUID()
	accDestinationID := vos.NewUUID()

	accounts := []entities.Account{
		{
			ID:        accOriginID,
			Name:      "Elliot",
			CPF:       "33344455567",
			Secret:    "password",
			Balance:   7000,
			CreatedAt: time.Now().Truncate(time.Second),
		},
		{
			ID:        accDestinationID,
			Name:      "Mr.Robot",
			CPF:       "33344455568",
			Secret:    "password",
			Balance:   3000,
			CreatedAt: time.Now().Truncate(time.Second),
		},
	}

	repo := NewRepository(NewDB(t))
	ctxDB := context.Background()

	for _, acc := range accounts {
		err := repo.CreateAccount(ctxDB, &acc)
		if err != nil {
			t.Error("error inserting accounts")
		}
	}

	tests := []struct {
		name        string
		transfer    entities.Transfer
		expectedErr error
	}{
		{
			name: "with success",
			transfer: entities.Transfer{
				ID:                   vos.NewUUID(),
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
				ID:                   vos.NewUUID(),
				AccountOriginID:      vos.NewUUID(),
				AccountDestinationID: vos.NewUUID(),
				Amount:               rand.Int(),
				CreatedAt:            time.Now().Truncate(time.Second),
			},
			expectedErr: domain.NewDBError(domain.QueryRefCreateTransfer, errors.New("any sql error"), errors.New("unexpected error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// execute
			resultErr := repo.CreateTransfer(context.Background(), &tt.transfer)

			// assert
			var GotDBError *domain.DBError
			var WantDBError *domain.DBError

			switch {
			case errors.As(resultErr, &GotDBError) && errors.As(tt.expectedErr, &WantDBError):
				if GotDBError.Query != WantDBError.Query {
					t.Errorf("got sql error in query: %v, want: %v", GotDBError.Query, WantDBError.Query)
				}
			case !errors.Is(resultErr, tt.expectedErr):
				t.Errorf("got error: %v want: %v", resultErr, tt.expectedErr)
			}
		})
	}
}

func TestTransferRepo_FetchTransfers(t *testing.T) {
	// setup
	accOriginID := vos.NewUUID()
	accDestinationID := vos.NewUUID()

	accounts := []entities.Account{
		{
			ID:        accOriginID,
			Name:      "Elliot",
			CPF:       "33344455567",
			Secret:    "password",
			Balance:   999999999999,
			CreatedAt: time.Now().Truncate(time.Second),
		},
		{
			ID:        accDestinationID,
			Name:      "Mr.Robot",
			CPF:       "33344455568",
			Secret:    "password",
			Balance:   999999999999,
			CreatedAt: time.Now().Truncate(time.Second),
		},
	}

	r := NewRepository(NewDB(t))
	ctxDB := context.Background()
	for _, acc := range accounts {
		err := r.CreateAccount(ctxDB, &acc)
		if err != nil {
			t.Fatal("error inserting accounts")
		}
	}

	var want []entities.Transfer
	for i := 0; i < 1000; i++ {
		transfer := &entities.Transfer{
			ID:                   vos.NewUUID(),
			AccountOriginID:      accOriginID,
			AccountDestinationID: accDestinationID,
			Amount:               rand.Intn(100),
			CreatedAt:            time.Now().Truncate(time.Second),
		}
		err := r.CreateTransfer(ctxDB, transfer)
		if err != nil {
			t.Fatal("error inserting transfers")
		}
		want = append(want, *transfer)
	}

	//execute
	result, err := r.FetchTransfers(context.Background(), accOriginID)

	// assert
	if err != nil {
		t.Errorf("didn't want sql error, but got the error: %v", err)
	}

	if !reflect.DeepEqual(want, result) {
		t.Errorf("got: %v, want: %v", result, want)
	}
}
