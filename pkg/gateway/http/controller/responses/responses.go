package responses

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

type errJSON struct {
	Err string `json:"error"`
}

//SendResponse sends formatted json response to request
func SendResponse(w http.ResponseWriter, statusCode int, data interface{}, log *logrus.Logger) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.WithError(err).Print("An encoding error occurred")
	}
}

func ErrorJSON(err error) *errJSON {
	return &errJSON{
		Err: err.Error(),
	}
}
