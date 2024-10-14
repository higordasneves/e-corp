package controller_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"

	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/gateway/config"
	"github.com/higordasneves/e-corp/pkg/gateway/controller"
	"github.com/higordasneves/e-corp/pkg/gateway/controller/mocks"
	"github.com/higordasneves/e-corp/pkg/gateway/controller/server"
)

func TestAccountController_ListAccounts_Success(t *testing.T) {
	t.Parallel()

	uc := &mocks.AccountUseCaseMock{
		ListAccountsFunc: func(ctx context.Context, input usecase.ListAccountsInput) (usecase.ListAccountsOutput, error) {
			assert.Equal(t, []uuid.UUID{uuid.FromStringOrNil("019282db-ff95-76ce-8ddd-ec5abceffa25"), uuid.FromStringOrNil("019282db-ff95-76cd-8b7f-c3a07b52a57c")}, input.IDs)
			assert.Equal(t, 100, input.PageSize)

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
				},
				NextPage: &usecase.ListAccountsInput{
					LastFetchedID: uuid.FromStringOrNil("019282db-ff95-76d0-a96d-41f561a1af28"),
					PageSize:      100,
				},
			}, nil
		},
	}
	want := `{"accounts":[{"id":"019282db-ff95-76cd-8b7f-c3a07b52a57c","name":"Elliot","document":"55566677780","balance":9700000,"created_at":"2024-01-01T00:00:00Z"},{"id":"019282db-ff95-76ce-8ddd-ec5abceffa25","name":"Mr. Robot","document":"55566677781","balance":5596400,"created_at":"2024-01-01T00:00:00Z"}],"next_page":"eyJJRHMiOm51bGwsIkxhc3RGZXRjaGVkSUQiOiIwMTkyODJkYi1mZjk1LTc2ZDAtYTk2ZC00MWY1NjFhMWFmMjgiLCJQYWdlU2l6ZSI6MTAwfQ=="}`
	accCtrl := controller.NewAccountController(uc)
	api := controller.API{
		AccountController: accCtrl,
	}

	handler := server.HTTPHandler(zaptest.NewLogger(t), api, config.Config{})
	urlValues := url.Values{
		"ids":        []string{strings.Join([]string{"019282db-ff95-76ce-8ddd-ec5abceffa25", "019282db-ff95-76cd-8b7f-c3a07b52a57c"}, ",")},
		"page_size":  []string{"100"},
		"page_token": []string{""},
	}
	req := httptest.NewRequest(http.MethodGet, "/api/v1/accounts?"+urlValues.Encode(), nil)
	response := httptest.NewRecorder()

	// execute
	handler.ServeHTTP(response, req)

	// assert
	assert.Equal(t, strings.TrimSpace(want), strings.TrimSpace(response.Body.String()))
	assert.Equal(t, 200, response.Code)

	// request with page token
	uc = &mocks.AccountUseCaseMock{
		ListAccountsFunc: func(ctx context.Context, input usecase.ListAccountsInput) (usecase.ListAccountsOutput, error) {
			assert.Equal(t, usecase.ListAccountsInput{
				LastFetchedID: uuid.FromStringOrNil("019282db-ff95-76d0-a96d-41f561a1af28"),
				PageSize:      100,
			}, input)
			return usecase.ListAccountsOutput{}, nil
		},
	}
	accCtrl = controller.NewAccountController(uc)
	api = controller.API{
		AccountController: accCtrl,
	}

	handler = server.HTTPHandler(zaptest.NewLogger(t), api, config.Config{})
	req = httptest.NewRequest(http.MethodGet, "/api/v1/accounts?"+url.Values{
		"page_token": []string{"eyJJRHMiOm51bGwsIkxhc3RGZXRjaGVkSUQiOiIwMTkyODJkYi1mZjk1LTc2ZDAtYTk2ZC00MWY1NjFhMWFmMjgiLCJQYWdlU2l6ZSI6MTAwfQ=="},
	}.Encode(), nil)
	response = httptest.NewRecorder()

	// execute
	handler.ServeHTTP(response, req)
	assert.Equal(t, http.StatusOK, response.Code)
}

func TestAccountController_ListAccounts_Failure(t *testing.T) {
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
			want:         fmt.Sprintf(`{"error":"%s"}`, controller.ErrUnexpected),
			expectedCode: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// setup
			accUseCase := tt.fields.accUseCase
			accCtrl := controller.NewAccountController(accUseCase)
			api := controller.API{
				AccountController: accCtrl,
			}

			handler := server.HTTPHandler(zaptest.NewLogger(t), api, config.Config{})
			urlValues := url.Values{
				"ids":        []string{uuid.Must(uuid.NewV4()).String()},
				"page_size":  []string{"100"},
				"page_token": []string{""},
			}
			req := httptest.NewRequest(http.MethodGet, "/api/v1/accounts?"+urlValues.Encode(), nil)
			response := httptest.NewRecorder()

			// execute
			handler.ServeHTTP(response, req)

			// assert
			assert.Equal(t, strings.TrimSpace(tt.want), strings.TrimSpace(response.Body.String()))
			assert.Equal(t, tt.expectedCode, response.Code)
		})
	}
}
