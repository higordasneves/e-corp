package controller

import (
	"context"
	"errors"
	"github.com/gorilla/mux"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	ucmock "github.com/higordasneves/e-corp/pkg/domain/usecase/mock"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"github.com/higordasneves/e-corp/pkg/gateway/http/controller/interpreter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAccountController_FetchAccountsWithSuccess(t *testing.T) {
	accountsList := []entities.AccountOutput{
		{
			ID:        vos.NewUUID(),
			Name:      "Elliot",
			CPF:       "555.666.777-80",
			Balance:   9700000,
			CreatedAt: time.Now().Truncate(time.Second),
		},
		{
			ID:        vos.NewUUID(),
			Name:      "Mr. Robot",
			CPF:       "555.666.777-81",
			Balance:   5596400,
			CreatedAt: time.Now().Truncate(time.Second),
		},
		{
			ID:        vos.NewUUID(),
			Name:      "WhiteRose",
			CPF:       "555.666.777-82",
			Balance:   5534513,
			CreatedAt: time.Now().Truncate(time.Second),
		},
		{
			ID:        vos.NewUUID(),
			Name:      "Darlene",
			CPF:       "555.666.777-83",
			Balance:   12350,
			CreatedAt: time.Now().Truncate(time.Second),
		},
	}

	type fields struct {
		accUseCase usecase.AccountUseCase
	}

	tests := []struct {
		name         string
		fields       fields
		want         []entities.AccountOutput
		expectedCode int
	}{
		{
			name: "with success",
			fields: fields{
				accUseCase: ucmock.AccountUseCase{
					Fetch: func(ctx context.Context) ([]entities.AccountOutput, error) {
						return accountsList, nil
					},
				},
			},
			want:         accountsList,
			expectedCode: 200,
		},
		{
			name: "empty database",
			fields: fields{
				accUseCase: ucmock.AccountUseCase{
					Fetch: func(ctx context.Context) ([]entities.AccountOutput, error) {
						return []entities.AccountOutput{}, nil
					},
				},
			},
			want:         []entities.AccountOutput{},
			expectedCode: 200,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			t.Helper()

			accUseCase := test.fields.accUseCase
			accCtrl := NewAccountController(accUseCase, logTest)

			router := mux.NewRouter()
			router.HandleFunc("/accounts", accCtrl.FetchAccounts).Methods(http.MethodGet)

			req := httptest.NewRequest(http.MethodGet, "/accounts", nil)
			response := httptest.NewRecorder()

			router.ServeHTTP(response, req)

			var responseBody []entities.AccountOutput
			err := decodeResponse(response, &responseBody)
			if err != nil {
				require.NoError(t, err)
			}

			assert.Equal(t, test.want, responseBody)
			assert.Equal(t, test.expectedCode, response.Code)
		})
	}
}

func TestAccountController_FetchAccountsWithError(t *testing.T) {
	type fields struct {
		accUseCase usecase.AccountUseCase
	}

	tests := []struct {
		name         string
		fields       fields
		want         errJSON
		expectedCode int
	}{
		{
			name: "unexpected error",
			fields: fields{
				accUseCase: ucmock.AccountUseCase{
					Fetch: func(ctx context.Context) ([]entities.AccountOutput, error) {
						return nil, errors.New("unknown error")
					},
				},
			},
			want:         errorJSON(interpreter.ErrUnexpected),
			expectedCode: 500,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Helper()

			accUseCase := test.fields.accUseCase
			accCtrl := NewAccountController(accUseCase, logTest)

			router := mux.NewRouter()
			router.HandleFunc("/accounts", accCtrl.FetchAccounts).Methods(http.MethodGet)

			req := httptest.NewRequest(http.MethodGet, "/accounts", nil)
			response := httptest.NewRecorder()

			router.ServeHTTP(response, req)

			var responseBody errJSON
			err := decodeResponse(response, &responseBody)
			if err != nil {
				require.NoError(t, err)
			}

			assert.Equal(t, test.want, responseBody)
			assert.Equal(t, test.expectedCode, response.Code)
		})
	}
}
