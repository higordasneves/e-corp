package controller

import (
	"context"
	"fmt"
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/higordasneves/e-corp/pkg/domain/entities"
	ucmock "github.com/higordasneves/e-corp/pkg/domain/usecase/mock"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
)

func TestAccountController_GetBalanceWithSuccess(t *testing.T) {
	type fields struct {
		accUseCase usecase.AccountUseCase
	}

	accountsList := []entities.AccountOutput{
		{
			ID:        "uuid1",
			Name:      "Elliot",
			CPF:       "555.666.777-80",
			Balance:   9700000,
			CreatedAt: time.Now().Truncate(time.Second),
		},
		{
			ID:        "uuid2",
			Name:      "Mr. Robot",
			CPF:       "555.666.777-81",
			Balance:   5596400,
			CreatedAt: time.Now().Truncate(time.Second),
		},
		{
			ID:        "uuid3",
			Name:      "WhiteRose",
			CPF:       "555.666.777-82",
			Balance:   5534513,
			CreatedAt: time.Now().Truncate(time.Second),
		},
		{
			ID:        "uuid4",
			Name:      "Darlene",
			CPF:       "555.666.777-83",
			Balance:   12350,
			CreatedAt: time.Now().Truncate(time.Second),
		},
	}

	tests := []struct {
		name         string
		fields       fields
		accID        vos.UUID
		want         interface{}
		expectedCode int
	}{
		{
			name: "with success, balance of 9700000 cents",
			fields: fields{
				accUseCase: ucmock.AccountUseCase{
					GetAccBalance: func(ctx context.Context, id vos.UUID) (int, error) {
						for _, acc := range accountsList {
							if acc.ID == id {
								return acc.Balance, nil
							}
						}
						return 0, entities.ErrAccNotFound
					},
				},
			},
			accID:        "uuid1",
			want:         map[string]int{"balance": 9700000},
			expectedCode: 200,
		},
		{
			name: "with success, balance of 5534513 cents",
			fields: fields{
				accUseCase: ucmock.AccountUseCase{
					GetAccBalance: func(ctx context.Context, id vos.UUID) (int, error) {
						for _, acc := range accountsList {
							if acc.ID == id {
								return acc.Balance, nil
							}
						}
						return 0, entities.ErrAccNotFound
					},
				},
			},
			accID:        "uuid3",
			want:         map[string]int{"balance": 5534513},
			expectedCode: 200,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Helper()

			accUseCase := test.fields.accUseCase
			accCtrl := NewAccountController(accUseCase, logTest)

			router := mux.NewRouter()
			router.HandleFunc("/accounts/{account_id}/balance", accCtrl.GetBalance).Methods(http.MethodGet)

			path := fmt.Sprintf("/accounts/%v/balance", test.accID)
			req := httptest.NewRequest(http.MethodGet, path, nil)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, req)

			var responseBody map[string]int

			err := decodeResponse(response, &responseBody)
			if err != nil {
				require.NoError(t, err)
			}

			assert.Equal(t, test.want, responseBody)
			assert.Equal(t, test.expectedCode, response.Code)
		})
	}
}

func TestAccountController_GetBalanceWithError(t *testing.T) {
	type fields struct {
		accUseCase usecase.AccountUseCase
	}

	tests := []struct {
		name         string
		accID        string
		fields       fields
		want         interface{}
		expectedCode int
		domainErr    error
	}{
		{
			name:  "account not found",
			accID: vos.NewUUID().String(),
			fields: fields{
				accUseCase: ucmock.AccountUseCase{
					GetAccBalance: func(ctx context.Context, id vos.UUID) (int, error) {
						return 0, entities.ErrAccNotFound
					},
				},
			},
			want:         errorJSON(entities.ErrAccNotFound),
			expectedCode: 404,
		},
		{
			name:  "invalid id",
			accID: "invalid",
			fields: fields{
				accUseCase: ucmock.AccountUseCase{
					GetAccBalance: func(ctx context.Context, id vos.UUID) (int, error) {
						return 0, vos.ErrInvalidID
					},
				},
			},
			want:         errorJSON(vos.ErrInvalidID),
			expectedCode: 400,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Helper()

			accUseCase := test.fields.accUseCase
			accCtrl := NewAccountController(accUseCase, logTest)

			router := mux.NewRouter()
			router.HandleFunc("/accounts/{account_id}/balance", accCtrl.GetBalance).Methods(http.MethodGet)

			path := fmt.Sprintf("/accounts/%v/balance", test.accID)
			req := httptest.NewRequest(http.MethodGet, path, nil)
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
