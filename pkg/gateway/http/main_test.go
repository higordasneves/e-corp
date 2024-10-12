package http_test

import (
	"github.com/sirupsen/logrus"
	"os"
	"testing"
)

var logTest *logrus.Logger

func TestMain(m *testing.M) {
	logTest = logrus.New()

	code := m.Run()
	os.Exit(code)
}
