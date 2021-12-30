package usecase

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"github.com/higordasneves/e-corp/pkg/repository"
	repomock "github.com/higordasneves/e-corp/pkg/repository/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAccountUseCase_GetBalance(t *testing.T) {
	t.Parallel()

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
			want:        0,
			expectedErr: repository.ErrUnexpected,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.insert {
				accounts = append(accounts, entities.Account{ID: tt.id, Balance: tt.want})
			}

			// setup
			accRepo := repomock.NewAccountRepo(accounts, tt.expectedErr)
			accUseCase := NewAccountUseCase(accRepo)

			// execute
			ctx := context.Background()
			balance, err := accUseCase.GetBalance(ctx, tt.id)

			// assert
			assert.Equal(t, tt.want, balance)
			assert.ErrorIs(t, tt.expectedErr, err)
		})
	}
}
