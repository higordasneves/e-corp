package controller_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/higordasneves/e-corp/pkg/domain"
	"github.com/higordasneves/e-corp/pkg/gateway/controller"
	"github.com/higordasneves/e-corp/pkg/gateway/controller/mocks"
)

func TestAccountController_GetBalance(t *testing.T) {
	t.Parallel()

	type fields struct {
		accUseCase controller.AccountUseCase
	}

	tests := []struct {
		name         string
		fields       fields
		accID        uuid.UUID
		want         string
		expectedCode int
	}{
		{
			name: "with success, balance of 9700000 cents",
			fields: fields{
				accUseCase: &mocks.AccountUseCaseMock{
					GetBalanceFunc: func(ctx context.Context, id uuid.UUID) (int, error) {
						return 9700000, nil
					},
				},
			},
			accID: uuid.Must(uuid.NewV7()),
			want: `{"balance":9700000}

`,
			expectedCode: 200,
		},
		{
			name: "with success, balance of 5534513 cents",
			fields: fields{
				accUseCase: &mocks.AccountUseCaseMock{
					GetBalanceFunc: func(ctx context.Context, id uuid.UUID) (int, error) {
						return 5534513, nil
					},
				},
			},
			accID:        uuid.Must(uuid.NewV7()),
			want:         `{"balance":5534513}`,
			expectedCode: 200,
		},
		{
			name:  "account not found",
			accID: uuid.Must(uuid.NewV7()),
			fields: fields{
				accUseCase: &mocks.AccountUseCaseMock{
					GetBalanceFunc: func(ctx context.Context, id uuid.UUID) (int, error) {
						return 0, domain.ErrNotFound
					},
				},
			},
			want:         fmt.Sprintf(`{"error":"%s"}`, domain.ErrNotFound),
			expectedCode: 404,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			//setup
			accUseCase := tt.fields.accUseCase
			accCtrl := controller.NewAccountController(accUseCase)
			router := mux.NewRouter()
			router.HandleFunc("/accounts/{account_id}/balance", accCtrl.GetBalance).Methods(http.MethodGet)
			path := fmt.Sprintf("/accounts/%v/balance", tt.accID)
			req := httptest.NewRequest(http.MethodGet, path, nil)
			response := httptest.NewRecorder()

			// execute
			router.ServeHTTP(response, req)

			//assert
			assert.Equal(t, tt.expectedCode, response.Code)
			assert.Equal(t, strings.TrimSpace(tt.want), strings.TrimSpace(response.Body.String()))
		})
	}
}
