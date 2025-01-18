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
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"github.com/higordasneves/e-corp/pkg/gateway/config"
	"github.com/higordasneves/e-corp/pkg/gateway/postgres"
	thelp "github.com/higordasneves/e-corp/utils/testhelpers"
)

func TestAuthUC_Login_Success(t *testing.T) {
	t.Parallel()

	r := postgres.NewRepository(NewDB(t))

	secret, err := vos.NewSecret("password123")
	require.NoError(t, err)

	err = r.CreateAccount(thelp.NewCtx(t), entities.Account{
		ID:        uuid.Must(uuid.NewV7()),
		Name:      "Elliot",
		Document:  "43663412309",
		Secret:    secret,
		Balance:   0,
		CreatedAt: time.Now(),
	})
	require.NoError(t, err)

	uc := usecase.NewAuthUC(r, &config.AuthConfig{
		Duration:  time.Minute,
		SecretKey: "secret_key_test",
	})

	output, err := uc.Login(thelp.NewCtx(t), usecase.LoginInput{
		Document: "43663412309",
		Secret:   "password123",
	})
	require.NoError(t, err)
	assert.WithinDuration(t, output.IssuedAt, time.Now(), time.Minute)
	assert.Equal(t, output.ExpiresAt, output.IssuedAt.Add(time.Minute))
}

func TestAuthUC_Login_Failure_InvalidPass(t *testing.T) {
	t.Parallel()

	r := postgres.NewRepository(NewDB(t))

	secret, err := vos.NewSecret("password123")
	require.NoError(t, err)

	err = r.CreateAccount(thelp.NewCtx(t), entities.Account{
		ID:        uuid.Must(uuid.NewV7()),
		Name:      "Elliot",
		Document:  "43663412309",
		Secret:    secret,
		Balance:   0,
		CreatedAt: time.Now(),
	})
	require.NoError(t, err)

	uc := usecase.NewAuthUC(r, &config.AuthConfig{
		Duration:  time.Minute,
		SecretKey: "secret_key_test",
	})

	_, err = uc.Login(thelp.NewCtx(t), usecase.LoginInput{
		Document: "43663412309",
		Secret:   "password124",
	})
	assert.ErrorIs(t, err, vos.ErrInvalidPass)
	assert.ErrorIs(t, err, domain.ErrInvalidParameter)
}
