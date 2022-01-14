package controller

import (
	"bytes"
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
	"time"
)

func TestTransferController_Transfer(t *testing.T) {
	t.Parallel()

	type fields struct {
		tUseCase usecase.TransferUseCase
	}

	type args struct {
		ctxWithValue context.Context
		requestBody  *bytes.Reader
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
					TransferFunc: func(ctx context.Context, transferInput *usecase.TransferInput) (*entities.Transfer, error) {
						return &entities.Transfer{
							ID:                   "transfer_id",
							AccountOriginID:      vos.UUID(transferInput.AccountOriginID),
							AccountDestinationID: vos.UUID(transferInput.AccountDestinationID),
							Amount:               transferInput.Amount,
							CreatedAt:            time.Now().Truncate(time.Minute),
						}, nil
					},
				},
			},
			args: args{
				ctxWithValue: context.WithValue(context.Background(), "subject", "uuid_acc1"),
				requestBody:  bytes.NewReader([]byte(`{"destinationID": "uuid_acc2", "amount": 10827}`)),
			},
			want: `{"id": "transfer_id",
				"account_origin_id": "uuid_acc1",
				"account_destination_id": "uuid_acc2",
				"amount": 10827,
				"created_at": "<<PRESENCE>>"}`,
			expectedCode: http.StatusCreated,
		},
		{
			name: "same account id in origin and destination should return an error and status code 400",
			fields: fields{
				tUseCase: &ucmock.TransferUseCase{
					TransferFunc: func(ctx context.Context, transferInput *usecase.TransferInput) (*entities.Transfer, error) {
						return nil, entities.ErrSelfTransfer
					},
				},
			},
			args: args{
				ctxWithValue: context.WithValue(context.Background(), "subject", "uuid_acc1"),
				requestBody:  bytes.NewReader([]byte(`{"destinationID": "uuid_acc1", "amount": 10}`)),
			},
			want:         fmt.Sprintf(`{"error": "%s"}`, entities.ErrSelfTransfer),
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid destination id should return an error and status code 400",
			fields: fields{
				tUseCase: &ucmock.TransferUseCase{
					TransferFunc: func(ctx context.Context, transferInput *usecase.TransferInput) (*entities.Transfer, error) {
						return nil, entities.ErrDestAccID
					},
				},
			},
			args: args{
				ctxWithValue: context.WithValue(context.Background(), "subject", "uuid_acc1"),
				requestBody:  bytes.NewReader([]byte(`{"destinationID": "invalid", "amount": 10}`)),
			},
			want:         fmt.Sprintf(`{"error": "%s"}`, entities.ErrDestAccID),
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid origin id should return an error and status code 400",
			fields: fields{
				tUseCase: &ucmock.TransferUseCase{
					TransferFunc: func(ctx context.Context, transferInput *usecase.TransferInput) (*entities.Transfer, error) {
						return nil, entities.ErrOriginAccID
					},
				},
			},
			args: args{
				ctxWithValue: context.WithValue(context.Background(), "subject", "invalid"),
				requestBody:  bytes.NewReader([]byte(`{"destinationID": "uuid_acc1", "amount": 10}`)),
			},
			want:         fmt.Sprintf(`{"error": "%s"}`, entities.ErrOriginAccID),
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "when transfer amount < 0 should return an error and status code 400",
			fields: fields{
				tUseCase: &ucmock.TransferUseCase{
					TransferFunc: func(ctx context.Context, transferInput *usecase.TransferInput) (*entities.Transfer, error) {
						return nil, entities.ErrTransferAmount
					},
				},
			},
			args: args{
				ctxWithValue: context.WithValue(context.Background(), "subject", "uuid_acc1"),
				requestBody:  bytes.NewReader([]byte(`{"destinationID": "uuid_acc2", "amount": -10}`)),
			},
			want:         fmt.Sprintf(`{"error": "%s"}`, entities.ErrTransferAmount),
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "when origin account balance < transfer amount should return an error and status code 400",
			fields: fields{
				tUseCase: &ucmock.TransferUseCase{
					TransferFunc: func(ctx context.Context, transferInput *usecase.TransferInput) (*entities.Transfer, error) {
						return nil, entities.ErrTransferInsufficientFunds
					},
				},
			},
			args: args{
				ctxWithValue: context.WithValue(context.Background(), "subject", "uuid_acc1"),
				requestBody:  bytes.NewReader([]byte(`{"destinationID": "uuid_acc2", "amount": 1000000000000000000}`)),
			},
			want:         fmt.Sprintf(`{"error": "%s"}`, entities.ErrTransferInsufficientFunds),
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "when destination account doesn't exists should return an error and status code 400",
			fields: fields{
				tUseCase: &ucmock.TransferUseCase{
					TransferFunc: func(ctx context.Context, transferInput *usecase.TransferInput) (*entities.Transfer, error) {
						return nil, entities.ErrAccNotFound
					},
				},
			},
			args: args{
				ctxWithValue: context.WithValue(context.Background(), "subject", "uuid_acc1"),
				requestBody:  bytes.NewReader([]byte(`{"destinationID": "uuid_acc2", "amount": 1000000000000000000}`)),
			},
			want:         fmt.Sprintf(`{"error": "%s"}`, entities.ErrAccNotFound),
			expectedCode: http.StatusNotFound,
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
			router.HandleFunc("/transfers", tCtrl.Transfer).Methods(http.MethodPost)
			req := httptest.NewRequest(http.MethodPost, "/transfers", tt.args.requestBody).WithContext(tt.args.ctxWithValue)
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
