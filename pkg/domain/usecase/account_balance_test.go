package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/gateway/postgres"
	thelp "github.com/higordasneves/e-corp/utils/testhelpers"
)

func TestAccountUseCase_GetBalance(t *testing.T) {
	t.Parallel()

	r := postgres.NewRepository(NewDB(t))
	err := r.CreateAccount(context.Background(), entities.Account{
		ID:        uuid.FromStringOrNil("5f2d4920-89c3-4ed5-af8e-1d411588746d"),
		Name:      "Elliot",
		Document:  "12345678900",
		Secret:    "secret",
		Balance:   10,
		CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local),
	})
	require.NoError(t, err)

	uc := usecase.GetAccountBalanceUC{R: r}
	got, err := uc.GetBalance(thelp.NewCtx(t), uuid.FromStringOrNil("5f2d4920-89c3-4ed5-af8e-1d411588746d"))
	require.NoError(t, err)
	assert.Equal(t, 10, got)
}
