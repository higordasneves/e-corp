package controller

import (
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
)

func TestAccountController_GetBalance(t *testing.T) {
	t.Parallel()
	type fields struct {
		accUseCase usecase.AccountUseCase
	}

	tests := []struct {
		name         string
		fields       fields
		accID        vos.UUID
		want         string
		expectedCode int
	}{
		{
			name: "with success, balance of 9700000 cents",
			fields: fields{
				accUseCase: ucmock.AccountUseCase{
					GetAccBalance: func(ctx context.Context, id vos.UUID) (int, error) {
						return 9700000, nil
					},
				},
			},
			accID:        "uuid1",
			want:         `{"balance": 9700000}`,
			expectedCode: 200,
		},
		{
			name: "with success, balance of 5534513 cents",
			fields: fields{
				accUseCase: ucmock.AccountUseCase{
					GetAccBalance: func(ctx context.Context, id vos.UUID) (int, error) {
						return 5534513, nil
					},
				},
			},
			accID:        "uuid3",
			want:         `{"balance": 5534513}`,
			expectedCode: 200,
		},
		{
			name:  "account not found",
			accID: vos.NewUUID(),
			fields: fields{
				accUseCase: ucmock.AccountUseCase{
					GetAccBalance: func(ctx context.Context, id vos.UUID) (int, error) {
						return 0, entities.ErrAccNotFound
					},
				},
			},
			want:         fmt.Sprintf(`{"error": "%s"}`, entities.ErrAccNotFound),
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
			want:         fmt.Sprintf(`{"error": "%s"}`, vos.ErrInvalidID),
			expectedCode: 400,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			//setup
			accUseCase := tt.fields.accUseCase
			accCtrl := NewAccountController(accUseCase, logTest)
			router := mux.NewRouter()
			router.HandleFunc("/accounts/{account_id}/balance", accCtrl.GetBalance).Methods(http.MethodGet)
			path := fmt.Sprintf("/accounts/%v/balance", tt.accID)
			req := httptest.NewRequest(http.MethodGet, path, nil)
			response := httptest.NewRecorder()

			// execute
			router.ServeHTTP(response, req)

			//assert
			ja := jsonassert.New(t)
			assert.Equal(t, tt.expectedCode, response.Code)
			ja.Assertf(strings.TrimSpace(response.Body.String()), tt.want)
		})
	}
}
