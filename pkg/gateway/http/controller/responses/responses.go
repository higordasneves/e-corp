package responses

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

//JSON returns a json response to request
func JSON(w http.ResponseWriter, statusCode int, data interface{}, log *logrus.Logger) {
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.WithError(err).Print("An encoding error occurred")
	}

}

func Error(w http.ResponseWriter, statusCode int, err error, log *logrus.Logger) {
	JSON(w, statusCode, struct {
		Erro string `json:"error"`
	}{
		Erro: err.Error(),
	}, log)
}
