package responses

import (
	"encoding/json"
	"errors"
	"github.com/higordasneves/e-corp/pkg/domain/entities"
	"github.com/sirupsen/logrus"
	"net/http"
)

type errJSON struct {
	Err string `json:"error"`
}

var (
	ErrUnexpected = errors.New("an unexpected error has occurred")
)

//SendResponse sends formatted json response to request
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
	case errors.Is(err, entities.ErrAccNotFound):
		statusCode = http.StatusNotFound
	case errors.Is(err, entities.ErrBadAccRequest):
		statusCode = http.StatusBadRequest
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
