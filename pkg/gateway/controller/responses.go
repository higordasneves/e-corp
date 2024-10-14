package controller

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"

	"github.com/higordasneves/e-corp/pkg/domain"
	"github.com/higordasneves/e-corp/utils/logger"
)

type ErrorResponse struct {
	Err string `json:"error"`
}

var ErrUnexpected = errors.New("internal server error")

// SendResponse sends formatted json response to request
func SendResponse(ctx context.Context, w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Error(ctx, "an unexpected encoding error occurred", zap.Error(err))
	}
}

func HandleError(ctx context.Context, w http.ResponseWriter, err error) {
	var statusCode int

	switch {
	case errors.Is(err, domain.ErrUnauthorized):
		statusCode = http.StatusUnauthorized
	case errors.Is(err, domain.ErrInvalidParameter):
		statusCode = http.StatusBadRequest
	case errors.Is(err, domain.ErrNotFound):
		statusCode = http.StatusNotFound
	default:
		statusCode = http.StatusInternalServerError
		logger.Error(ctx, "an unexpected error occurred", zap.Error(err))
		err = ErrUnexpected
	}

	SendError(ctx, w, statusCode, err)
}

func SendError(ctx context.Context, w http.ResponseWriter, statusCode int, err error) {
	jsonError := errorJSON(err)
	SendResponse(ctx, w, statusCode, jsonError)
}

func errorJSON(err error) ErrorResponse {
	return ErrorResponse{
		Err: err.Error(),
	}
}
