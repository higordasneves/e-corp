package controller

import (
	"fmt"
	"net/http"

	"github.com/gofrs/uuid/v5"

	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/gateway/http/controller/reponses"
	"github.com/higordasneves/e-corp/pkg/gateway/http/controller/requests"
)

func (tController TransferController) Transfer(w http.ResponseWriter, r *http.Request) {
	var transferInput usecase.TransferInput
	if err := requests.ReadRequestBody(r, &transferInput); err != nil {
		reponses.HandleError(w, err, tController.log)
		return
	}

	accountOriginID := uuid.FromStringOrNil(fmt.Sprint(r.Context().Value("subject")))
	transferInput.AccountOriginID = accountOriginID

	transfer, err := tController.tUseCase.Transfer(r.Context(), &transferInput)
	if err != nil {
		reponses.HandleError(w, err, tController.log)
		return
	}
	reponses.SendResponse(w, http.StatusCreated, transfer, tController.log)
}
