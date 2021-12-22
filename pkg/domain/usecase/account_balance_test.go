package usecase

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"github.com/higordasneves/e-corp/pkg/repository"
	repomock "github.com/higordasneves/e-corp/pkg/repository/mock"
	"testing"
)

func TestAccountUseCase_GetBalance(t *testing.T) {
	ctx := context.Background()

	accounts := make([]entities.Account, 0, 10)

	tests := []struct {
		name        string
		id          vos.UUID
		insert      bool
		want        int
		expectedErr error
	}{
		{
			name:        "with success 1",
			id:          vos.NewUUID(),
			insert:      true,
			want:        162000,
			expectedErr: nil,
		},
		{
			name:        "with success 2",
			id:          vos.NewUUID(),
			insert:      true,
			want:        561300,
			expectedErr: nil,
		},
		{
			name:        "err account not found",
			id:          vos.NewUUID(),
			insert:      false,
			want:        0,
			expectedErr: entities.ErrAccNotFound,
		},
		{
			name:        "database generic error",
			id:          vos.NewUUID(),
			insert:      false,
			want:        561300,
			expectedErr: repository.ErrUnexpected,
		},
	}

	for _, test := range tests {
		if test.insert {
			accounts = append(accounts, entities.Account{ID: test.id, Balance: test.want})
		}
		accRepo := repomock.NewAccountRepo(accounts, test.expectedErr)
		accUseCase := NewAccountUseCase(accRepo)
		balance, err := accUseCase.GetBalance(ctx, test.id)

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

			if test.expectedErr == nil && balance != test.want {
				t.Errorf("got %v, want %v", balance, test.want)
			}
		})
	}
}
