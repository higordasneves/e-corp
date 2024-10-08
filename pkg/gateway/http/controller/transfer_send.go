package controller

import (
	"fmt"
	"github.com/gofrs/uuid/v5"
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/gateway/http/controller/interpreter"
	"net/http"
)

func (tController TransferController) Transfer(w http.ResponseWriter, r *http.Request) {
	var transferInput usecase.TransferInput
	if err := interpreter.ReadRequestBody(r, &transferInput); err != nil {
		interpreter.HandleError(w, err, tController.log)
		return
	}

	accountOriginID := uuid.FromStringOrNil(fmt.Sprint(r.Context().Value("subject")))
	transferInput.AccountOriginID = accountOriginID

	transfer, err := tController.tUseCase.Transfer(r.Context(), &transferInput)
	if err != nil {
		interpreter.HandleError(w, err, tController.log)
		return
	}
	interpreter.SendResponse(w, http.StatusCreated, transfer, tController.log)
}
