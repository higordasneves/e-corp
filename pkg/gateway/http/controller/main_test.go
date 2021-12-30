package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http/httptest"
	"os"
	"testing"
)

var logTest *logrus.Logger

type errJSON struct {
	Err string `json:"error"`
}

func TestMain(m *testing.M) {
	logTest = logrus.New()

	code := m.Run()
	os.Exit(code)
}

func decodeResponse(w *httptest.ResponseRecorder, obj interface{}) error {
	if err := json.NewDecoder(w.Body).Decode(&obj); err != nil {
		return fmt.Errorf("%w: %s", errors.New("error decoding response"), err)
	}
	return nil
}

func errorJSON(err error) errJSON {
	return errJSON{
		Err: err.Error(),
	}
}
