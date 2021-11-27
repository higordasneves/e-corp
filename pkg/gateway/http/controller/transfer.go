package controller

import (
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/sirupsen/logrus"
	"net/http"
)

type TransferController interface {
	Transfer(w http.ResponseWriter, r *http.Request)
	GetTransfers(w http.ResponseWriter, r *http.Request)
}

type transferController struct {
	tUseCase usecase.TransferUseCase
	log      *logrus.Logger
}

func NewTransferController(tUseCase usecase.TransferUseCase, log *logrus.Logger) TransferController {
	return &transferController{tUseCase: tUseCase, log: log}
}
func (tController transferController) Transfer(w http.ResponseWriter, r *http.Request) {

}

func (tController transferController) GetTransfers(w http.ResponseWriter, r *http.Request) {

}
