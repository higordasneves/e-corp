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

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/higordasneves/e-corp/pkg/domain"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"github.com/higordasneves/e-corp/pkg/gateway/controller"
	"github.com/higordasneves/e-corp/pkg/gateway/controller/mocks"
)

func TestAccountController_CreateAccount(t *testing.T) {
	t.Parallel()

	type fields struct {
		accUseCase controller.AccountUseCase
	}

	tests := []struct {
		name         string
		requestBody  *bytes.Reader
		fields       fields
		want         string
		expectedCode int
	}{
		{
			name:        "with success",
			requestBody: bytes.NewReader([]byte(`{"name":"Elliot", "document":"44455566678", "secret":"12345678"}`)),
			fields: fields{
				accUseCase: &mocks.AccountUseCaseMock{
					CreateAccountFunc: func(ctx context.Context, input usecase.CreateAccountInput) (usecase.CreateAccountOutput, error) {
						return usecase.CreateAccountOutput{
							Account: entities.Account{
								ID:        uuid.FromStringOrNil("5f2d4920-89c3-4ed5-af8e-1d411588746d"),
								Name:      input.Name,
								Document:  vos.Document(input.Document),
								Balance:   0,
								CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
							},
						}, nil
					},
				},
			},
			want:         `{"id":"5f2d4920-89c3-4ed5-af8e-1d411588746d","name":"Elliot","document":"44455566678","balance":0,"created_at":"2024-01-01T00:00:00Z"}`,
			expectedCode: http.StatusCreated,
		},
		{
			name:        "when account already exists should return error and status code 400",
			requestBody: bytes.NewReader([]byte(`{"name":"Elliot", "document":"44455566678", "secret":"12345678"}`)),
			fields: fields{
				accUseCase: &mocks.AccountUseCaseMock{
					CreateAccountFunc: func(ctx context.Context, input usecase.CreateAccountInput) (usecase.CreateAccountOutput, error) {
						return usecase.CreateAccountOutput{}, domain.ErrInvalidParameter
					},
				},
			},
			want:         fmt.Sprintf(`{"error":"%s"}`, domain.ErrInvalidParameter),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "invalid cpf length should return error and status code 400",
			requestBody: bytes.NewReader([]byte(`{"name":"Elliot", "document":"111", "secret":"12345678"}`)),
			fields: fields{
				accUseCase: &mocks.AccountUseCaseMock{
					CreateAccountFunc: func(ctx context.Context, input usecase.CreateAccountInput) (usecase.CreateAccountOutput, error) {
						return usecase.CreateAccountOutput{}, domain.ErrInvalidParameter
					},
				},
			},
			want:         fmt.Sprintf(`{"error":"%s"}`, domain.ErrInvalidParameter),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "invalid cpf format should return error and status code 400",
			requestBody: bytes.NewReader([]byte(`{"name":"Elliot", "document":"111.233", "secret":"12345678"}`)),
			fields: fields{
				accUseCase: &mocks.AccountUseCaseMock{
					CreateAccountFunc: func(ctx context.Context, input usecase.CreateAccountInput) (usecase.CreateAccountOutput, error) {
						return usecase.CreateAccountOutput{}, domain.ErrInvalidParameter
					},
				},
			},
			want:         fmt.Sprintf(`{"error":"%s"}`, domain.ErrInvalidParameter),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "invalid secret length should return error and status code 400",
			requestBody: bytes.NewReader([]byte(`{"name":"Elliot", "document":"44455566678", "secret":"123456"}`)),
			fields: fields{
				accUseCase: &mocks.AccountUseCaseMock{
					CreateAccountFunc: func(ctx context.Context, input usecase.CreateAccountInput) (usecase.CreateAccountOutput, error) {
						return usecase.CreateAccountOutput{}, domain.ErrInvalidParameter
					},
				},
			},
			want:         fmt.Sprintf(`{"error":"%s"}`, domain.ErrInvalidParameter),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "empty required fields should return error and status code 400",
			requestBody: bytes.NewReader([]byte(`{"name":"", "document":"", "secret":""}`)),
			fields: fields{
				accUseCase: &mocks.AccountUseCaseMock{
					CreateAccountFunc: func(ctx context.Context, input usecase.CreateAccountInput) (usecase.CreateAccountOutput, error) {
						return usecase.CreateAccountOutput{}, domain.ErrInvalidParameter
					},
				},
			},
			want:         fmt.Sprintf(`{"error":"%s"}`, domain.ErrInvalidParameter),
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// setup
			accUseCase := tt.fields.accUseCase
			accController := controller.NewAccountController(accUseCase, logTest)

			router := mux.NewRouter()
			router.HandleFunc("/accounts", accController.CreateAccount).Methods(http.MethodPost)
			req := httptest.NewRequest(http.MethodPost, "/accounts", tt.requestBody)
			response := httptest.NewRecorder()

			//execute
			router.ServeHTTP(response, req)

			//assert
			assert.Equal(t, strings.TrimSpace(tt.want), strings.TrimSpace(response.Body.String()))
			assert.Equal(t, tt.expectedCode, response.Code)
		})
	}
}
