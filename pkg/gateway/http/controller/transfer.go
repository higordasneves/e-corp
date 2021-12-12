package controller

import (
	"fmt"
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/gateway/http/controller/io"
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
	var transferInput usecase.TransferInput
	if err := io.ReadRequestBody(r, &transferInput); err != nil {
		io.HandleError(w, err, tController.log)
		return
	}

	accountOriginID := fmt.Sprint(r.Context().Value("subject"))
	transferInput.AccountOriginID = accountOriginID

	transfer, err := tController.tUseCase.Transfer(r.Context(), &transferInput)
	if err != nil {
		io.HandleError(w, err, tController.log)
		return
	}
	io.SendResponse(w, http.StatusCreated, transfer, tController.log)
}

func (tController transferController) GetTransfers(w http.ResponseWriter, r *http.Request) {
	panic("implement me!")
}
