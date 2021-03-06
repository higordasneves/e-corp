package controller

import (
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
	"time"
)

func TestTransferController_FetchTransfers(t *testing.T) {
	t.Parallel()

	type fields struct {
		tUseCase usecase.TransferUseCase
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
				tUseCase: &ucmock.TransferUseCase{
					FetchTransfersFunc: func(ctx context.Context, id string) ([]entities.Transfer, error) {
						return []entities.Transfer{
							{
								ID:                   "transfer_id1",
								AccountOriginID:      vos.UUID(id),
								AccountDestinationID: vos.UUID("uuid_destination_acc"),
								Amount:               2000,
								CreatedAt:            time.Now().Truncate(time.Minute),
							},
							{
								ID:                   "transfer_id2",
								AccountOriginID:      vos.UUID(id),
								AccountDestinationID: vos.UUID("uuid_destination_acc2"),
								Amount:               4598,
								CreatedAt:            time.Now().Truncate(time.Minute),
							},
						}, nil
					},
				},
			},
			args: args{ctxWithValue: context.WithValue(context.Background(), "subject", "uuid_acc1")},
			want: `[{
						"id": "transfer_id1", "account_origin_id": "uuid_acc1",
						"account_destination_id": "uuid_destination_acc",
						"amount": 2000,
						"created_at": "<<PRESENCE>>"
					},
					{
						"id": "transfer_id2",
						"account_origin_id": "uuid_acc1",
						"account_destination_id": "uuid_destination_acc2",
						"amount": 4598,
						"created_at": "<<PRESENCE>>"
					}]`,
			expectedCode: http.StatusOK,
		},
		{
			name: "no transfers, should return empty list and status code 200",
			fields: fields{
				tUseCase: &ucmock.TransferUseCase{
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
				tUseCase: &ucmock.TransferUseCase{
					FetchTransfersFunc: func(ctx context.Context, id string) ([]entities.Transfer, error) {
						return nil, errors.New("new error")
					},
				},
			},
			args:         args{ctxWithValue: context.WithValue(context.Background(), "subject", "uuid_acc1")},
			want:         fmt.Sprintf(`{"error": "%s"}`, interpreter.ErrUnexpected),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// setup
			tUseCase := tt.fields.tUseCase
			tCtrl := NewTransferController(tUseCase, logTest)

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
