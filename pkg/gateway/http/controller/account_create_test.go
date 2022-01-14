package controller

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	ucmock "github.com/higordasneves/e-corp/pkg/domain/usecase/mock"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const balanceInit = 1000000

func TestAccountController_CreateAccount(t *testing.T) {
	t.Parallel()
	type fields struct {
		accUseCase usecase.AccountUseCase
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
			requestBody: bytes.NewReader([]byte(`{"name":"Elliot", "cpf":"44455566678", "secret":"12345678"}`)),
			fields: fields{
				accUseCase: &ucmock.AccountUseCase{
					CreateAccountFunc: func(ctx context.Context, input *usecase.AccountInput) (*entities.AccountOutput, error) {
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
			want:         fmt.Sprintf(`{"id": "uuid1", "name": "Elliot", "cpf": "444.555.666-78", "balance": %v, "created_at": "<<PRESENCE>>"}`, balanceInit),
			expectedCode: http.StatusCreated,
		},
		{
			name:        "when account already exists should return error and status code 400",
			requestBody: bytes.NewReader([]byte(`{"name":"Elliot", "cpf":"44455566678", "secret":"12345678"}`)),
			fields: fields{
				accUseCase: &ucmock.AccountUseCase{
					CreateAccountFunc: func(ctx context.Context, input *usecase.AccountInput) (*entities.AccountOutput, error) {
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
				accUseCase: &ucmock.AccountUseCase{
					CreateAccountFunc: func(ctx context.Context, input *usecase.AccountInput) (*entities.AccountOutput, error) {
						return nil, vos.ErrCPFLen
					},
				},
			},
			want:         fmt.Sprintf(`{"error": "%s"}`, vos.ErrCPFLen),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "invalid cpf format should return error and status code 400",
			requestBody: bytes.NewReader([]byte(`{"name":"Elliot", "cpf":"111.233", "secret":"12345678"}`)),
			fields: fields{
				accUseCase: &ucmock.AccountUseCase{
					CreateAccountFunc: func(ctx context.Context, input *usecase.AccountInput) (*entities.AccountOutput, error) {
						return nil, vos.ErrCPFFormat
					},
				},
			},
			want:         fmt.Sprintf(`{"error": "%s"}`, vos.ErrCPFFormat),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "invalid secret length should return error and status code 400",
			requestBody: bytes.NewReader([]byte(`{"name":"Elliot", "cpf":"44455566678", "secret":"123456"}`)),
			fields: fields{
				accUseCase: &ucmock.AccountUseCase{
					CreateAccountFunc: func(ctx context.Context, input *usecase.AccountInput) (*entities.AccountOutput, error) {
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
				accUseCase: &ucmock.AccountUseCase{
					CreateAccountFunc: func(ctx context.Context, input *usecase.AccountInput) (*entities.AccountOutput, error) {
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
			accController := NewAccountController(accUseCase, logTest)

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
