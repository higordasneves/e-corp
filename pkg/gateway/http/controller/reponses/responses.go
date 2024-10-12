package reponses

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/higordasneves/e-corp/pkg/domain"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/higordasneves/e-corp/pkg/domain/vos"
	"github.com/higordasneves/e-corp/pkg/gateway/http/controller/requests"
)

type errJSON struct {
	Err string `json:"error"`
}

var (
	ErrUnexpected = errors.New("an unexpected error occurred")
)

// SendResponse sends formatted json response to request
func SendResponse(w http.ResponseWriter, statusCode int, data interface{}, log *logrus.Logger) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.WithError(err).Print("an encoding error occurred")
	}
}

func HandleError(w http.ResponseWriter, err error, log *logrus.Logger) {
	var statusCode int
	var dbError *domain.DBError

	switch {
	case errors.Is(err, requests.ErrReadRequest):
		statusCode = http.StatusBadRequest
	case errors.Is(err, requests.ErrTokenFormat):
		statusCode = http.StatusUnauthorized
	case errors.Is(err, requests.ErrReadRequest):
		statusCode = http.StatusBadRequest
	case errors.Is(err, entities.ErrAccNotFound) || errors.Is(err, domain.ErrNotFound):
		statusCode = http.StatusNotFound
	case errors.Is(err, domain.ErrInvalidParameter):
		statusCode = http.StatusBadRequest
	case errors.Is(err, entities.ErrEmptyInput):
		statusCode = http.StatusBadRequest
	case errors.Is(err, entities.ErrAccAlreadyExists):
		statusCode = http.StatusBadRequest
		statusCode = http.StatusBadRequest
	case errors.Is(err, vos.ErrDocumentFormat):
		statusCode = http.StatusBadRequest
	case errors.Is(err, vos.ErrDocumentLen):
		statusCode = http.StatusBadRequest
	case errors.Is(err, vos.ErrInvalidPass):
		statusCode = http.StatusBadRequest
	case errors.Is(err, vos.ErrSmallSecret):
		statusCode = http.StatusBadRequest
	case errors.Is(err, entities.ErrZeroRowsAffectedUpdateBalance):
		statusCode = http.StatusInternalServerError
		err = ErrUnexpected
	case errors.As(err, &dbError):
		statusCode = http.StatusInternalServerError
		log.WithField("query", dbError.Query).WithError(dbError.DBErr).Error("unexpected sql error has occurred")
		err = dbError.GenericErr
	default:
		statusCode = http.StatusInternalServerError
		log.WithError(err).Println(ErrUnexpected)
		err = ErrUnexpected
	}

	SendError(w, statusCode, err, log)
}

func SendError(w http.ResponseWriter, statusCode int, err error, log *logrus.Logger) {
	jsonError := errorJSON(err)
	SendResponse(w, statusCode, jsonError, log)
}

func errorJSON(err error) errJSON {
	return errJSON{
		Err: err.Error(),
	}
}
