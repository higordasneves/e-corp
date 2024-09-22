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
	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"

	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"github.com/higordasneves/e-corp/pkg/gateway/http/controller"
	"github.com/higordasneves/e-corp/pkg/gateway/http/controller/mocks"
)

const balanceInit = 1000000

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
					CreateAccountFunc: func(ctx context.Context, input *usecase.CreateAccountInput) (*entities.AccountOutput, error) {
						return &entities.AccountOutput{
							ID:        uuid.FromStringOrNil("5f2d4920-89c3-4ed5-af8e-1d411588746d"),
							Name:      input.Name,
							CPF:       input.Document,
							Balance:   balanceInit,
							CreatedAt: time.Now().Truncate(time.Minute),
						}, nil
					},
				},
			},
			want:         fmt.Sprintf(`{"id": "5f2d4920-89c3-4ed5-af8e-1d411588746d", "name": "Elliot", "cpf": "44455566678", "balance": %v, "created_at": "<<PRESENCE>>"}`, balanceInit),
			expectedCode: http.StatusCreated,
		},
		{
			name:        "when account already exists should return error and status code 400",
			requestBody: bytes.NewReader([]byte(`{"name":"Elliot", "cpf":"44455566678", "secret":"12345678"}`)),
			fields: fields{
				accUseCase: &mocks.AccountUseCaseMock{
					CreateAccountFunc: func(ctx context.Context, input *usecase.CreateAccountInput) (*entities.AccountOutput, error) {
						return nil, entities.ErrAccAlreadyExists
					},
				},
			},
			want:         fmt.Sprintf(`{"error": "%s"}`, entities.ErrAccAlreadyExists),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "invalid cpf length should return error and status code 400",
			requestBody: bytes.NewReader([]byte(`{"name":"Elliot", "cpf":"111", "secret":"12345678"}`)),
			fields: fields{
				accUseCase: &mocks.AccountUseCaseMock{
					CreateAccountFunc: func(ctx context.Context, input *usecase.CreateAccountInput) (*entities.AccountOutput, error) {
						return nil, vos.ErrDocumentLen
					},
				},
			},
			want:         fmt.Sprintf(`{"error": "%s"}`, vos.ErrDocumentLen),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "invalid cpf format should return error and status code 400",
			requestBody: bytes.NewReader([]byte(`{"name":"Elliot", "cpf":"111.233", "secret":"12345678"}`)),
			fields: fields{
				accUseCase: &mocks.AccountUseCaseMock{
					CreateAccountFunc: func(ctx context.Context, input *usecase.CreateAccountInput) (*entities.AccountOutput, error) {
						return nil, vos.ErrDocumentFormat
					},
				},
			},
			want:         fmt.Sprintf(`{"error": "%s"}`, vos.ErrDocumentFormat),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "invalid secret length should return error and status code 400",
			requestBody: bytes.NewReader([]byte(`{"name":"Elliot", "cpf":"44455566678", "secret":"123456"}`)),
			fields: fields{
				accUseCase: &mocks.AccountUseCaseMock{
					CreateAccountFunc: func(ctx context.Context, input *usecase.CreateAccountInput) (*entities.AccountOutput, error) {
						return nil, vos.ErrSmallSecret
					},
				},
			},
			want:         fmt.Sprintf(`{"error": "%s"}`, vos.ErrSmallSecret),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "empty required fields should return error and status code 400",
			requestBody: bytes.NewReader([]byte(`{"name":"", "cpf":"", "secret":""}`)),
			fields: fields{
				accUseCase: &mocks.AccountUseCaseMock{
					CreateAccountFunc: func(ctx context.Context, input *usecase.CreateAccountInput) (*entities.AccountOutput, error) {
						return nil, entities.ErrEmptyInput
					},
				},
			},
			want:         fmt.Sprintf(`{"error": "%s"}`, entities.ErrEmptyInput),
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		tt := tt
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
			ja := jsonassert.New(t)
			ja.Assertf(strings.TrimSpace(response.Body.String()), tt.want)
			assert.Equal(t, tt.expectedCode, response.Code)
		})
	}
}
