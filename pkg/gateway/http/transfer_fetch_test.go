package http_test

import (
	"context"
	"errors"
	"fmt"
	http2 "github.com/higordasneves/e-corp/pkg/gateway/http"
	"github.com/higordasneves/e-corp/pkg/gateway/http/mocks"
	"github.com/higordasneves/e-corp/pkg/gateway/http/reponses"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"

	"github.com/higordasneves/e-corp/pkg/domain/entities"
)

func TestTransferController_FetchTransfers(t *testing.T) {
	t.Parallel()

	type fields struct {
		tUseCase http2.TransferUseCase
	}

	type args struct {
		ctxWithValue context.Context
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		want         string
		expectedCode int
	}{
		{
			name: "with success",
			fields: fields{
				tUseCase: &mocks.TransferUseCaseMock{
					FetchTransfersFunc: func(ctx context.Context, id string) ([]entities.Transfer, error) {
						return []entities.Transfer{
							{
								ID:                   uuid.FromStringOrNil("8b07e65f-7fed-4387-ba84-d2213527c6f1"),
								AccountOriginID:      uuid.FromStringOrNil(id),
								AccountDestinationID: uuid.FromStringOrNil("9751fe39-976f-4b3d-9611-d6c8c6370b0f"),
								Amount:               2000,
								CreatedAt:            time.Now().Truncate(time.Minute),
							},
							{
								ID:                   uuid.FromStringOrNil("6ca1469e-1def-445c-b6ad-1028689d72f2"),
								AccountOriginID:      uuid.FromStringOrNil(id),
								AccountDestinationID: uuid.FromStringOrNil("9ee14852-1011-422e-b9f3-abd905d5103c"),
								Amount:               4598,
								CreatedAt:            time.Now().Truncate(time.Minute),
							},
						}, nil
					},
				},
			},
			args: args{ctxWithValue: context.WithValue(context.Background(), "subject", "0457c690-f884-4d57-810c-85cf09a50d8b")},
			want: `[{
						"id": "8b07e65f-7fed-4387-ba84-d2213527c6f1", 
						"account_origin_id": "0457c690-f884-4d57-810c-85cf09a50d8b",
						"account_destination_id": "9751fe39-976f-4b3d-9611-d6c8c6370b0f",
						"amount": 2000,
						"created_at": "<<PRESENCE>>"
					},
					{
						"id": "6ca1469e-1def-445c-b6ad-1028689d72f2",
						"account_origin_id": "0457c690-f884-4d57-810c-85cf09a50d8b",
						"account_destination_id": "9ee14852-1011-422e-b9f3-abd905d5103c",
						"amount": 4598,
						"created_at": "<<PRESENCE>>"
					}]`,
			expectedCode: http.StatusOK,
		},
		{
			name: "no transfers, should return empty list and status code 200",
			fields: fields{
				tUseCase: &mocks.TransferUseCaseMock{
					FetchTransfersFunc: func(ctx context.Context, id string) ([]entities.Transfer, error) {
						return []entities.Transfer{}, nil
					},
				},
			},
			args:         args{ctxWithValue: context.WithValue(context.Background(), "subject", "uuid_acc1")},
			want:         `{"msg": "no transfers"}`,
			expectedCode: http.StatusOK,
		},
		{
			name: "unknown error should return unexpected error and status code 500",
			fields: fields{
				tUseCase: &mocks.TransferUseCaseMock{
					FetchTransfersFunc: func(ctx context.Context, id string) ([]entities.Transfer, error) {
						return nil, errors.New("new error")
					},
				},
			},
			args:         args{ctxWithValue: context.WithValue(context.Background(), "subject", "uuid_acc1")},
			want:         fmt.Sprintf(`{"error": "%s"}`, reponses.ErrUnexpected),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// setup
			tUseCase := tt.fields.tUseCase
			tCtrl := http2.NewTransferController(tUseCase, logTest)

			router := mux.NewRouter()
			router.HandleFunc("/transfers", tCtrl.FetchTransfers).Methods(http.MethodGet)
			req := httptest.NewRequest(http.MethodGet, "/transfers", nil).WithContext(tt.args.ctxWithValue)
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
