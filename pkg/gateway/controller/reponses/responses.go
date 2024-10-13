package reponses

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/higordasneves/e-corp/pkg/domain"
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

	switch {
	case errors.Is(err, domain.ErrUnauthorized):
		statusCode = http.StatusUnauthorized
	case errors.Is(err, domain.ErrInvalidParameter):
		statusCode = http.StatusBadRequest
	case errors.Is(err, domain.ErrNotFound):
		statusCode = http.StatusNotFound
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
