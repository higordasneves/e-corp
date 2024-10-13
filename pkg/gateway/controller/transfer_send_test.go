package controller_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/higordasneves/e-corp/pkg/domain"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/gateway/config"
	"github.com/higordasneves/e-corp/pkg/gateway/controller"
	"github.com/higordasneves/e-corp/pkg/gateway/controller/mocks"
	"github.com/higordasneves/e-corp/pkg/gateway/controller/server"
)

func TestTransferController_Transfer(t *testing.T) {
	t.Parallel()

	type fields struct {
		tUseCase controller.TransferUseCase
	}

	type args struct {
		ctxWithValue context.Context
		requestBody  *bytes.Reader
	}

	tests := []struct {
		name         string
		fields       fields
		args         args
		want         string
		expectedCode int
	}{
		{
			name: "with success",
			fields: fields{
				tUseCase: &mocks.TransferUseCaseMock{
					TransferFunc: func(ctx context.Context, input usecase.TransferInput) (usecase.TransferOutput, error) {
						return usecase.TransferOutput{
							Transfer: entities.Transfer{
								ID:                   uuid.FromStringOrNil("9ee14852-1011-422e-b9f3-abd905d5103c"),
								AccountOriginID:      input.AccountOriginID,
								AccountDestinationID: input.AccountDestinationID,
								Amount:               input.Amount,
								CreatedAt:            time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
							},
						}, nil
					},
				},
			},
			args: args{
				ctxWithValue: context.WithValue(context.Background(), "subject", "b59c5660-d62f-4f3e-91b4-5f8e236e5d3d"),
				requestBody:  bytes.NewReader([]byte(`{"destination_id": "5f2d4920-89c3-4ed5-af8e-1d411588746d", "amount": 10827}`)),
			},
			want:         `{"id":"9ee14852-1011-422e-b9f3-abd905d5103c","account_origin_id":"b59c5660-d62f-4f3e-91b4-5f8e236e5d3d","account_destination_id":"5f2d4920-89c3-4ed5-af8e-1d411588746d","amount":10827,"created_at":"2024-01-01T00:00:00Z"}`,
			expectedCode: http.StatusCreated,
		},
		{
			name: "same account id in origin and destination should return an error and status code 400",
			fields: fields{
				tUseCase: &mocks.TransferUseCaseMock{
					TransferFunc: func(ctx context.Context, input usecase.TransferInput) (usecase.TransferOutput, error) {
						return usecase.TransferOutput{}, domain.ErrInvalidParameter
					},
				},
			},
			args: args{
				ctxWithValue: context.WithValue(context.Background(), "subject", "b59c5660-d62f-4f3e-91b4-5f8e236e5d3d"),
				requestBody:  bytes.NewReader([]byte(`{"destinationID": "b59c5660-d62f-4f3e-91b4-5f8e236e5d3d", "amount": 10}`)),
			},
			want:         fmt.Sprintf(`{"error":"%s"}`, domain.ErrInvalidParameter),
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid origin id should return an error and status code 400",
			fields: fields{
				tUseCase: &mocks.TransferUseCaseMock{
					TransferFunc: func(ctx context.Context, input usecase.TransferInput) (usecase.TransferOutput, error) {
						return usecase.TransferOutput{}, domain.ErrInvalidParameter
					},
				},
			},
			args: args{
				ctxWithValue: context.WithValue(context.Background(), "subject", "invalid"),
				requestBody:  bytes.NewReader([]byte(`{"destinationID": "b59c5660-d62f-4f3e-91b4-5f8e236e5d3d", "amount": 10}`)),
			},
			want:         fmt.Sprintf(`{"error":"%s"}`, domain.ErrInvalidParameter),
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "when transfer amount < 0 should return an error and status code 400",
			fields: fields{
				tUseCase: &mocks.TransferUseCaseMock{
					TransferFunc: func(ctx context.Context, input usecase.TransferInput) (usecase.TransferOutput, error) {
						return usecase.TransferOutput{}, domain.ErrInvalidParameter
					},
				},
			},
			args: args{
				ctxWithValue: context.WithValue(context.Background(), "subject", "b59c5660-d62f-4f3e-91b4-5f8e236e5d3d"),
				requestBody:  bytes.NewReader([]byte(`{"destinationID": "5f2d4920-89c3-4ed5-af8e-1d411588746d", "amount": -10}`)),
			},
			want:         fmt.Sprintf(`{"error":"%s"}`, domain.ErrInvalidParameter),
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "when origin account balance < transfer amount should return an error and status code 400",
			fields: fields{
				tUseCase: &mocks.TransferUseCaseMock{
					TransferFunc: func(ctx context.Context, input usecase.TransferInput) (usecase.TransferOutput, error) {
						return usecase.TransferOutput{}, domain.ErrInvalidParameter
					},
				},
			},
			args: args{
				ctxWithValue: context.WithValue(context.Background(), "subject", "b59c5660-d62f-4f3e-91b4-5f8e236e5d3d"),
				requestBody:  bytes.NewReader([]byte(`{"destinationID": "5f2d4920-89c3-4ed5-af8e-1d411588746d", "amount": 1000000000000000000}`)),
			},
			want:         fmt.Sprintf(`{"error":"%s"}`, domain.ErrInvalidParameter),
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "when destination account doesn't exists should return an error and status code 400",
			fields: fields{
				tUseCase: &mocks.TransferUseCaseMock{
					TransferFunc: func(ctx context.Context, input usecase.TransferInput) (usecase.TransferOutput, error) {
						return usecase.TransferOutput{}, domain.ErrNotFound
					},
				},
			},
			args: args{
				ctxWithValue: context.WithValue(context.Background(), "subject", "b59c5660-d62f-4f3e-91b4-5f8e236e5d3d"),
				requestBody:  bytes.NewReader([]byte(`{"destinationID": "5f2d4920-89c3-4ed5-af8e-1d411588746d", "amount": 1000000000000000000}`)),
			},
			want:         fmt.Sprintf(`{"error":"%s"}`, domain.ErrNotFound),
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// setup
			tUseCase := tt.fields.tUseCase
			tCtrl := controller.NewTransferController(tUseCase)
			api := controller.API{
				TransferController: tCtrl,
			}

			now := time.Now()
			claims := &jwt.StandardClaims{
				Issuer:    "login",
				Subject:   "b59c5660-d62f-4f3e-91b4-5f8e236e5d3d",
				IssuedAt:  now.UTC().Unix(),
				ExpiresAt: now.UTC().Add(time.Hour).Unix(),
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := token.SignedString([]byte("test_secret_key"))
			require.NoError(t, err)

			handler := server.HTTPHandler(zaptest.NewLogger(t), api, config.Config{Auth: config.AuthConfig{SecretKey: "test_secret_key"}})
			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/transfers"), tt.args.requestBody)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))
			response := httptest.NewRecorder()

			// execute
			handler.ServeHTTP(response, req)

			// assert
			assert.Equal(t, strings.TrimSpace(tt.want), strings.TrimSpace(response.Body.String()))
			assert.Equal(t, tt.expectedCode, response.Code)
		})
	}
}
