package http_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/gofrs/uuid/v5"
	http2 "github.com/higordasneves/e-corp/pkg/gateway/http"
	"github.com/higordasneves/e-corp/pkg/gateway/http/mocks"
	"github.com/higordasneves/e-corp/pkg/gateway/http/reponses"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
)

func TestAuthController_Login(t *testing.T) {
	t.Parallel()

	type fields struct {
		authUC http2.AuthUseCase
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
			requestBody: bytes.NewReader([]byte(`{"cpf": "44455566678", "secret": "12345678"}`)),
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
			requestBody: bytes.NewReader([]byte(`{"cpf": "44455566690", "secret": "12345678"}`)),
			fields: fields{
				authUC: &mocks.AuthUseCaseMock{
					LoginFunc: func(ctx context.Context, input usecase.LoginInput) (usecase.LoginOutput, error) {
						return usecase.LoginOutput{}, entities.ErrAccNotFound
					},
				},
			},
			want:         fmt.Sprintf(`{"error":"%s"}`, entities.ErrAccNotFound),
			expectedCode: http.StatusNotFound,
		},
		{
			name:        "invalid password should return error and status code 400",
			requestBody: bytes.NewReader([]byte(`{"cpf": "44455566690", "secret": "123456"}`)),
			fields: fields{
				authUC: &mocks.AuthUseCaseMock{
					LoginFunc: func(ctx context.Context, input usecase.LoginInput) (usecase.LoginOutput, error) {
						return usecase.LoginOutput{}, vos.ErrInvalidPass
					},
				},
			},
			want:         fmt.Sprintf(`{"error":"%s"}`, vos.ErrInvalidPass),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "invalid cpf format should return error and status code 400",
			requestBody: bytes.NewReader([]byte(`{"cpf": "444.555.666-90", "secret": "12345678"}`)),
			fields: fields{
				authUC: &mocks.AuthUseCaseMock{
					LoginFunc: func(ctx context.Context, input usecase.LoginInput) (usecase.LoginOutput, error) {
						return usecase.LoginOutput{}, vos.ErrDocumentFormat
					},
				},
			},
			want:         fmt.Sprintf(`{"error":"%s"}`, vos.ErrDocumentFormat),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "unknown error should return status code 500",
			requestBody: bytes.NewReader([]byte(`{"cpf": "444.555.666-90", "secret": "12345678"}`)),
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
			authCtrl := http2.NewAuthController(authUC, "test_secret_key", logTest)

			router := mux.NewRouter()
			router.HandleFunc("/login", authCtrl.Login).Methods(http.MethodPost)
			req := httptest.NewRequest(http.MethodPost, "/login", tt.requestBody)
			response := httptest.NewRecorder()

			// execute
			router.ServeHTTP(response, req)

			// assert
			assert.Contains(t, strings.TrimSpace(response.Body.String()), tt.want)
			assert.Equal(t, tt.expectedCode, response.Code)
		})
	}
}
