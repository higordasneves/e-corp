package controller_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/gateway/controller"
	"github.com/higordasneves/e-corp/pkg/gateway/controller/mocks"
	"github.com/higordasneves/e-corp/pkg/gateway/controller/reponses"
)

func TestAccountController_ListAccounts(t *testing.T) {
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
					ListAccountsFunc: func(ctx context.Context, input usecase.ListAccountsInput) (usecase.ListAccountsOutput, error) {
						return usecase.ListAccountsOutput{
							Accounts: []entities.Account{
								{
									ID:        uuid.FromStringOrNil("019282db-ff95-76cd-8b7f-c3a07b52a57c"),
									Name:      "Elliot",
									Document:  "55566677780",
									Balance:   9700000,
									CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
								},
								{
									ID:        uuid.FromStringOrNil("019282db-ff95-76ce-8ddd-ec5abceffa25"),
									Name:      "Mr. Robot",
									Document:  "55566677781",
									Balance:   5596400,
									CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
								},
								{
									ID:        uuid.FromStringOrNil("019282db-ff95-76cf-a31f-101349333a13"),
									Name:      "WhiteRose",
									Document:  "55566677782",
									Balance:   5534513,
									CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
								},
								{
									ID:        uuid.FromStringOrNil("019282db-ff95-76d0-a96d-41f561a1af27"),
									Name:      "Darlene",
									Document:  "55566677783",
									Balance:   12350,
									CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
								},
							},
							NextPage: nil,
						}, nil
					},
				},
			},
			want:         `{"accounts":[{"id":"019282db-ff95-76cd-8b7f-c3a07b52a57c","name":"Elliot","document":"55566677780","balance":9700000,"created_at":"2024-01-01T00:00:00Z"},{"id":"019282db-ff95-76ce-8ddd-ec5abceffa25","name":"Mr. Robot","document":"55566677781","balance":5596400,"created_at":"2024-01-01T00:00:00Z"},{"id":"019282db-ff95-76cf-a31f-101349333a13","name":"WhiteRose","document":"55566677782","balance":5534513,"created_at":"2024-01-01T00:00:00Z"},{"id":"019282db-ff95-76d0-a96d-41f561a1af27","name":"Darlene","document":"55566677783","balance":12350,"created_at":"2024-01-01T00:00:00Z"}],"next_page":""}`,
			expectedCode: 200,
		},
		{
			name: "empty database",
			fields: fields{
				accUseCase: &mocks.AccountUseCaseMock{
					ListAccountsFunc: func(ctx context.Context, input usecase.ListAccountsInput) (usecase.ListAccountsOutput, error) {
						return usecase.ListAccountsOutput{}, nil
					},
				},
			},
			want:         `{"accounts":[],"next_page":""}`,
			expectedCode: 200,
		},
		{
			name: "unexpected error",
			fields: fields{
				accUseCase: &mocks.AccountUseCaseMock{
					ListAccountsFunc: func(ctx context.Context, input usecase.ListAccountsInput) (usecase.ListAccountsOutput, error) {
						return usecase.ListAccountsOutput{}, errors.New("unknown error")
					},
				},
			},
			want:         fmt.Sprintf(`{"error":"%s"}`, reponses.ErrUnexpected),
			expectedCode: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// setup
			accUseCase := tt.fields.accUseCase
			accCtrl := controller.NewAccountController(accUseCase, logTest)

			router := mux.NewRouter()
			router.HandleFunc("/accounts", accCtrl.ListAccounts).Methods(http.MethodGet)

			req := httptest.NewRequest(http.MethodGet, "/accounts", nil)
			response := httptest.NewRecorder()

			// execute
			router.ServeHTTP(response, req)

			// assert
			fmt.Println(response.Body.String())
			assert.Equal(t, strings.TrimSpace(tt.want), strings.TrimSpace(response.Body.String()))
			assert.Equal(t, tt.expectedCode, response.Code)
		})
	}
}
