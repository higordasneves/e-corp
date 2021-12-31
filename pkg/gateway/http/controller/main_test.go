package controller

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
