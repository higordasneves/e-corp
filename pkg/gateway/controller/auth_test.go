package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
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

	"github.com/higordasneves/e-corp/pkg/domain"
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/gateway/controller"
	"github.com/higordasneves/e-corp/pkg/gateway/controller/mocks"
	"github.com/higordasneves/e-corp/pkg/gateway/controller/reponses"
)

func TestAuthController_Login(t *testing.T) {
	t.Parallel()

	type fields struct {
		authUC controller.AuthUseCase
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
			requestBody: bytes.NewReader([]byte(`{"document": "44455566678", "secret": "12345678"}`)),
			fields: fields{
				authUC: &mocks.AuthUseCaseMock{
					LoginFunc: func(ctx context.Context, input usecase.LoginInput) (usecase.LoginOutput, error) {
						return usecase.LoginOutput{
							AccountID: uuid.Must(uuid.NewV7()),
							IssuedAt:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local),
							ExpiresAt: time.Date(2024, 1, 1, 0, 1, 0, 0, time.Local),
						}, nil
					},
				},
			},
			want:         "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			expectedCode: http.StatusOK,
		},
		{
			name:        "when account not found should return error and status code 400",
			requestBody: bytes.NewReader([]byte(`{"document": "44455566690", "secret": "12345678"}`)),
			fields: fields{
				authUC: &mocks.AuthUseCaseMock{
					LoginFunc: func(ctx context.Context, input usecase.LoginInput) (usecase.LoginOutput, error) {
						return usecase.LoginOutput{}, domain.ErrNotFound
					},
				},
			},
			want:         fmt.Sprintf(`{"error":"%s"}`, domain.ErrNotFound),
			expectedCode: http.StatusNotFound,
		},
		{
			name:        "invalid password should return error and status code 400",
			requestBody: bytes.NewReader([]byte(`{"document": "44455566690", "secret": "123456"}`)),
			fields: fields{
				authUC: &mocks.AuthUseCaseMock{
					LoginFunc: func(ctx context.Context, input usecase.LoginInput) (usecase.LoginOutput, error) {
						return usecase.LoginOutput{}, domain.ErrInvalidParameter
					},
				},
			},
			want:         fmt.Sprintf(`{"error":"%s"}`, domain.ErrInvalidParameter),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "invalid document format should return error and status code 400",
			requestBody: bytes.NewReader([]byte(`{"document": "444.555.666-90", "secret": "12345678"}`)),
			fields: fields{
				authUC: &mocks.AuthUseCaseMock{
					LoginFunc: func(ctx context.Context, input usecase.LoginInput) (usecase.LoginOutput, error) {
						return usecase.LoginOutput{}, domain.ErrInvalidParameter
					},
				},
			},
			want:         fmt.Sprintf(`{"error":"%s"}`, domain.ErrInvalidParameter),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "unknown error should return status code 500",
			requestBody: bytes.NewReader([]byte(`{"document": "444.555.666-90", "secret": "12345678"}`)),
			fields: fields{
				authUC: &mocks.AuthUseCaseMock{
					LoginFunc: func(ctx context.Context, input usecase.LoginInput) (usecase.LoginOutput, error) {
						return usecase.LoginOutput{}, errors.New("something")
					},
				},
			},

			want:         fmt.Sprintf(`{"error":"%s"}`, reponses.ErrUnexpected),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// setup
			authUC := tt.fields.authUC
			authCtrl := controller.NewAuthController(authUC, "test_secret_key", logTest)

			router := mux.NewRouter()
			router.HandleFunc("/login", authCtrl.Login).Methods(http.MethodPost)
			req := httptest.NewRequest(http.MethodPost, "/login", tt.requestBody)
			response := httptest.NewRecorder()

			// execute
			router.ServeHTTP(response, req)

			// assert
			assert.Contains(t, strings.TrimSpace(response.Body.String()), tt.want)
			assert.Equal(t, tt.expectedCode, response.Code)
			if response.Code == http.StatusOK {
				var got controller.LoginResponse
				err := json.NewDecoder(response.Body).Decode(&got)
				assert.NoError(t, err)
				assert.Contains(t, got.Token, tt.want)
			}
		})
	}
}
