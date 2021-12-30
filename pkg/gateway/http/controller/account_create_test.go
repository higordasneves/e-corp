package controller

import (
	"bytes"
	"context"
	"github.com/gorilla/mux"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	ucmock "github.com/higordasneves/e-corp/pkg/domain/usecase/mock"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const balanceInit = 1000000

func TestAccountController_CreateAccountWithSuccess(t *testing.T) {
	type fields struct {
		accUseCase usecase.AccountUseCase
	}

	tests := []struct {
		name         string
		requestBody  *bytes.Reader
		fields       fields
		want         entities.AccountOutput
		expectedCode int
	}{
		{
			name:        "with success",
			requestBody: bytes.NewReader([]byte(`{"name":"Elliot", "cpf":"44455566678", "secret":"12345678"}`)),
			fields: fields{
				accUseCase: ucmock.AccountUseCase{
					Create: func(ctx context.Context, input *usecase.AccountInput) (*entities.AccountOutput, error) {
						return &entities.AccountOutput{
							ID:        "uuid1",
							Name:      input.Name,
							CPF:       input.CPF.FormatOutput(),
							Balance:   balanceInit,
							CreatedAt: time.Now().Truncate(time.Minute),
						}, nil
					},
				},
			},
			want: entities.AccountOutput{
				ID:        "uuid1",
				Name:      "Elliot",
				CPF:       "444.555.666-78",
				Balance:   balanceInit,
				CreatedAt: time.Now().Truncate(time.Minute),
			},
			expectedCode: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()
			t.Parallel()

			accUseCase := tt.fields.accUseCase
			accController := NewAccountController(accUseCase, logTest)

			router := mux.NewRouter()
			router.HandleFunc("/accounts", accController.CreateAccount).Methods(http.MethodPost)
			req := httptest.NewRequest(http.MethodPost, "/accounts", tt.requestBody)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, req)

			var requestBody entities.AccountOutput
			err := decodeResponse(response, &requestBody)
			require.NoError(t, err)

			assert.Equal(t, tt.want, requestBody)
			assert.Equal(t, tt.expectedCode, response.Code)
		})
	}
}

func TestAccountController_CreateAccountWithError(t *testing.T) {
	type fields struct {
		accUseCase usecase.AccountUseCase
	}

	tests := []struct {
		name         string
		requestBody  *bytes.Reader
		fields       fields
		want         errJSON
		expectedCode int
	}{
		{
			name:        "when account already exists should return error and status code 400",
			requestBody: bytes.NewReader([]byte(`{"name":"Elliot", "cpf":"44455566678", "secret":"12345678"}`)),
			fields: fields{
				accUseCase: ucmock.AccountUseCase{
					Create: func(ctx context.Context, input *usecase.AccountInput) (*entities.AccountOutput, error) {
						return nil, entities.ErrAccAlreadyExists
					},
				},
			},
			want:         errorJSON(entities.ErrAccAlreadyExists),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "invalid cpf length should return error and status code 400",
			requestBody: bytes.NewReader([]byte(`{"name":"Elliot", "cpf":"111", "secret":"12345678"}`)),
			fields: fields{
				accUseCase: ucmock.AccountUseCase{
					Create: func(ctx context.Context, input *usecase.AccountInput) (*entities.AccountOutput, error) {
						return nil, vos.ErrCPFLen
					},
				},
			},
			want:         errorJSON(vos.ErrCPFLen),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "invalid cpf format should return error and status code 400",
			requestBody: bytes.NewReader([]byte(`{"name":"Elliot", "cpf":"111.233", "secret":"12345678"}`)),
			fields: fields{
				accUseCase: ucmock.AccountUseCase{
					Create: func(ctx context.Context, input *usecase.AccountInput) (*entities.AccountOutput, error) {
						return nil, vos.ErrCPFFormat
					},
				},
			},
			want:         errorJSON(vos.ErrCPFFormat),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "invalid secret length should return error and status code 400",
			requestBody: bytes.NewReader([]byte(`{"name":"Elliot", "cpf":"44455566678", "secret":"123456"}`)),
			fields: fields{
				accUseCase: ucmock.AccountUseCase{
					Create: func(ctx context.Context, input *usecase.AccountInput) (*entities.AccountOutput, error) {
						return nil, vos.ErrSmallSecret
					},
				},
			},
			want:         errorJSON(vos.ErrSmallSecret),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "empty required fields should return error and status code 400",
			requestBody: bytes.NewReader([]byte(`{"name":"", "cpf":"", "secret":""}`)),
			fields: fields{
				accUseCase: ucmock.AccountUseCase{
					Create: func(ctx context.Context, input *usecase.AccountInput) (*entities.AccountOutput, error) {
						return nil, entities.ErrEmptyInput
					},
				},
			},
			want:         errorJSON(entities.ErrEmptyInput),
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()
			t.Parallel()

			accUseCase := tt.fields.accUseCase
			accController := NewAccountController(accUseCase, logTest)

			router := mux.NewRouter()
			router.HandleFunc("/accounts", accController.CreateAccount).Methods(http.MethodPost)
			req := httptest.NewRequest(http.MethodPost, "/accounts", tt.requestBody)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, req)

			var requestBody errJSON
			err := decodeResponse(response, &requestBody)
			require.NoError(t, err)

			assert.Equal(t, tt.want, requestBody)
			assert.Equal(t, tt.expectedCode, response.Code)
		})
	}
}
