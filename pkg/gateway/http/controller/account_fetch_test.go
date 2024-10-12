package controller_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"

	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/gateway/http/controller"
	"github.com/higordasneves/e-corp/pkg/gateway/http/controller/mocks"
	"github.com/higordasneves/e-corp/pkg/gateway/http/controller/reponses"
)

func TestAccountController_FetchAccounts(t *testing.T) {
	t.Parallel()
	type fields struct {
		accUseCase controller.AccountUseCase
	}

	tests := []struct {
		name         string
		fields       fields
		want         string
		expectedCode int
	}{
		{
			name: "with success",
			fields: fields{
				accUseCase: &mocks.AccountUseCaseMock{
					FetchAccountsFunc: func(ctx context.Context) ([]entities.AccountOutput, error) {
						return []entities.AccountOutput{
							{
								ID:        uuid.Must(uuid.NewV7()),
								Name:      "Elliot",
								CPF:       "555.666.777-80",
								Balance:   9700000,
								CreatedAt: time.Now().Truncate(time.Second),
							},
							{
								ID:        uuid.Must(uuid.NewV7()),
								Name:      "Mr. Robot",
								CPF:       "555.666.777-81",
								Balance:   5596400,
								CreatedAt: time.Now().Truncate(time.Second),
							},
							{
								ID:        uuid.Must(uuid.NewV7()),
								Name:      "WhiteRose",
								CPF:       "555.666.777-82",
								Balance:   5534513,
								CreatedAt: time.Now().Truncate(time.Second),
							},
							{
								ID:        uuid.Must(uuid.NewV7()),
								Name:      "Darlene",
								CPF:       "555.666.777-83",
								Balance:   12350,
								CreatedAt: time.Now().Truncate(time.Second),
							},
						}, nil
					},
				},
			},
			want: `[{"id":"<<PRESENCE>>","name":"Elliot","cpf":"555.666.777-80","balance":9700000,"created_at":"<<PRESENCE>>"},
					{"id":"<<PRESENCE>>","name":"Mr. Robot","cpf":"555.666.777-81","balance":5596400,"created_at":"<<PRESENCE>>"},
					{"id":"<<PRESENCE>>","name":"WhiteRose","cpf":"555.666.777-82","balance":5534513,"created_at":"<<PRESENCE>>"},
					{"id":"<<PRESENCE>>","name":"Darlene","cpf":"555.666.777-83","balance":12350,"created_at":"<<PRESENCE>>"}]`,
			expectedCode: 200,
		},
		{
			name: "empty database",
			fields: fields{
				accUseCase: &mocks.AccountUseCaseMock{
					FetchAccountsFunc: func(ctx context.Context) ([]entities.AccountOutput, error) {
						return []entities.AccountOutput{}, nil
					},
				},
			},
			want:         `[]`,
			expectedCode: 200,
		},
		{
			name: "unexpected error",
			fields: fields{
				accUseCase: &mocks.AccountUseCaseMock{
					FetchAccountsFunc: func(ctx context.Context) ([]entities.AccountOutput, error) {
						return nil, errors.New("unknown error")
					},
				},
			},
			want:         fmt.Sprintf(`{"error": "%s"}`, reponses.ErrUnexpected),
			expectedCode: 500,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// setup
			accUseCase := tt.fields.accUseCase
			accCtrl := controller.NewAccountController(accUseCase, logTest)

			router := mux.NewRouter()
			router.HandleFunc("/accounts", accCtrl.FetchAccounts).Methods(http.MethodGet)

			req := httptest.NewRequest(http.MethodGet, "/accounts", nil)
			response := httptest.NewRecorder()

			// execute
			router.ServeHTTP(response, req)

			// assert
			ja := jsonassert.New(t)
			ja.Assertf(strings.TrimSpace(response.Body.String()), tt.want)
			assert.Equal(t, tt.expectedCode, response.Code)
		})
	}
}
