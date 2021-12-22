package controller

import (
	"fmt"
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/gateway/http/controller/interpreter"
	"github.com/sirupsen/logrus"
	"net/http"
)

type TransferController interface {
	Transfer(w http.ResponseWriter, r *http.Request)
	FetchTransfers(w http.ResponseWriter, r *http.Request)
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
	if err := interpreter.ReadRequestBody(r, &transferInput); err != nil {
		interpreter.HandleError(w, err, tController.log)
		return
	}

	accountOriginID := fmt.Sprint(r.Context().Value("subject"))
	transferInput.AccountOriginID = accountOriginID

	transfer, err := tController.tUseCase.Transfer(r.Context(), &transferInput)
	if err != nil {
		interpreter.HandleError(w, err, tController.log)
		return
	}
	interpreter.SendResponse(w, http.StatusCreated, transfer, tController.log)
}

func (tController transferController) FetchTransfers(w http.ResponseWriter, r *http.Request) {
	accountOriginID := fmt.Sprint(r.Context().Value("subject"))
	transferList, err := tController.tUseCase.FetchTransfers(r.Context(), accountOriginID)

	if err != nil {
		interpreter.HandleError(w, err, tController.log)
		return
	}

	if len(transferList) > 0 {
		interpreter.SendResponse(w, http.StatusOK, transferList, tController.log)

	} else {
		noTransfers := &map[string]string{"msg": "no transfers"}
		interpreter.SendResponse(w, http.StatusOK, noTransfers, tController.log)
	}
}
