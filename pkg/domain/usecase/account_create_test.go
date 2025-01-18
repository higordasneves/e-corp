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
	"github.com/higordasneves/e-corp/pkg/domain/usecase/mocks"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"github.com/higordasneves/e-corp/pkg/gateway/postgres"
	thelp "github.com/higordasneves/e-corp/utils/testhelpers"
)

func TestAccountUseCase_CreateAccount_Success(t *testing.T) {
	t.Parallel()

	// setup
	ctx := thelp.NewCtx(t)
	r := postgres.NewRepository(NewDB(t))
	uc := usecase.CreateAccountUC{R: r, B: &mocks.CreateAccountUCBrokerMock{}}

	// execute
	got, err := uc.CreateAccount(ctx, usecase.CreateAccountInput{
		Name:     "Elliot",
		Document: "43663312487",
		Secret:   "password123@",
	})
	require.NoError(t, err)

	// assert
	assert.Equal(t, "Elliot", got.Account.Name)
	assert.Equal(t, vos.Document("43663312487"), got.Account.Document)
	assert.Equal(t, 0, got.Account.Balance)
	assert.WithinDuration(t, time.Now(), got.Account.CreatedAt, time.Hour)
}

func TestAccountUseCase_CreateAccount_Failure(t *testing.T) {
	t.Parallel()

	// setup
	ctx := thelp.NewCtx(t)
	r := postgres.NewRepository(NewDB(t))
	uc := usecase.CreateAccountUC{R: r, B: &mocks.CreateAccountUCBrokerMock{}}

	tests := []struct {
		name        string
		setup       func(t *testing.T, r postgres.Repository)
		input       usecase.CreateAccountInput
		wantErr     error
		errContains string
	}{
		{
			name:  "empty account name",
			setup: func(t *testing.T, r postgres.Repository) {},
			input: usecase.CreateAccountInput{
				Name:     "",
				Document: "1234567899",
				Secret:   "secret123@",
			},
			wantErr:     domain.ErrInvalidParameter,
			errContains: "(name): required field",
		},
		{
			name:  "document length",
			setup: func(t *testing.T, r postgres.Repository) {},
			input: usecase.CreateAccountInput{
				Name:     "Elliot",
				Document: "123",
				Secret:   "secret123@",
			},
			wantErr:     domain.ErrInvalidParameter,
			errContains: vos.ErrDocumentLen.Error(),
		},
		{
			name:  "document format",
			setup: func(t *testing.T, r postgres.Repository) {},
			input: usecase.CreateAccountInput{
				Name:     "Elliot",
				Document: "1234567890A",
				Secret:   "secret123@",
			},
			wantErr:     domain.ErrInvalidParameter,
			errContains: vos.ErrDocumentFormat.Error(),
		},
		{
			name:  "secret length",
			setup: func(t *testing.T, r postgres.Repository) {},
			input: usecase.CreateAccountInput{
				Name:     "Elliot",
				Document: "12345678901",
				Secret:   "123",
			},
			wantErr:     domain.ErrInvalidParameter,
			errContains: vos.ErrSmallSecret.Error(),
		},
		{
			name: "account already exists",
			setup: func(t *testing.T, r postgres.Repository) {
				acc := entities.Account{
					ID:        uuid.Must(uuid.NewV7()),
					Name:      "Elliot",
					Document:  "43663312344",
					Secret:    "validsecret123",
					Balance:   0,
					CreatedAt: time.Now(),
				}
				require.NoError(t, r.CreateAccount(ctx, acc))
			},
			input: usecase.CreateAccountInput{
				Name:     "Elliot",
				Document: "43663312344",
				Secret:   "validsecret123",
			},
			wantErr:     domain.ErrInvalidParameter,
			errContains: "account with document 43663312344 already exists",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.setup(t, r)

			_, err := uc.CreateAccount(ctx, tt.input)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.ErrorContains(t, err, tt.errContains)
		})
	}
}
