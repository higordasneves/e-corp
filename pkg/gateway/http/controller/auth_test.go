package controller

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	ucmock "github.com/higordasneves/e-corp/pkg/domain/usecase/mock"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"github.com/higordasneves/e-corp/pkg/gateway/http/controller/interpreter"
	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAuthController_Login(t *testing.T) {
	t.Parallel()
	type fields struct {
		authUC usecase.AuthUseCase
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
				authUC: ucmock.AuthUseCase{
					AuthLogin: func(ctx context.Context, input *usecase.LoginInput) (*usecase.Token, error) {
						var token usecase.Token = "fake_token"
						return &token, nil
					},
				},
			},
			want:         `"fake_token"`,
			expectedCode: http.StatusOK,
		},
		{
			name:        "when account not found should return error and status code 400",
			requestBody: bytes.NewReader([]byte(`{"cpf": "44455566690", "secret": "12345678"}`)),
			fields: fields{
				authUC: ucmock.AuthUseCase{
					AuthLogin: func(ctx context.Context, input *usecase.LoginInput) (*usecase.Token, error) {
						return nil, entities.ErrAccNotFound
					},
				},
			},
			want:         fmt.Sprintf(`{"error": "%s"}`, entities.ErrAccNotFound),
			expectedCode: http.StatusNotFound,
		},
		{
			name:        "invalid password should return error and status code 400",
			requestBody: bytes.NewReader([]byte(`{"cpf": "44455566690", "secret": "123456"}`)),
			fields: fields{
				authUC: ucmock.AuthUseCase{
					AuthLogin: func(ctx context.Context, input *usecase.LoginInput) (*usecase.Token, error) {
						return nil, vos.ErrInvalidPass
					},
				},
			},
			want:         fmt.Sprintf(`{"error": "%s"}`, vos.ErrInvalidPass),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "invalid cpf format should return error and status code 400",
			requestBody: bytes.NewReader([]byte(`{"cpf": "444.555.666-90", "secret": "12345678"}`)),
			fields: fields{
				authUC: ucmock.AuthUseCase{
					AuthLogin: func(ctx context.Context, input *usecase.LoginInput) (*usecase.Token, error) {
						return nil, vos.ErrCPFFormat
					},
				},
			},
			want:         fmt.Sprintf(`{"error": "%s"}`, vos.ErrCPFFormat),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:        "unknown error should return status code 500",
			requestBody: bytes.NewReader([]byte(`{"cpf": "444.555.666-90", "secret": "12345678"}`)),
			fields: fields{
				authUC: ucmock.AuthUseCase{
					AuthLogin: func(ctx context.Context, input *usecase.LoginInput) (*usecase.Token, error) {
						return nil, errors.New("something")
					},
				},
			},

			want:         fmt.Sprintf(`{"error": "%s"}`, interpreter.ErrUnexpected),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// setup
			authUC := tt.fields.authUC
			authCtrl := NewAuthController(authUC, logTest)

			router := mux.NewRouter()
			router.HandleFunc("/login", authCtrl.Login).Methods(http.MethodPost)
			req := httptest.NewRequest(http.MethodPost, "/login", tt.requestBody)
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
