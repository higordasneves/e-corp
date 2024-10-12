package controller

import (
	"fmt"
	"github.com/higordasneves/e-corp/pkg/gateway/controller/reponses"
	"net/http"
)

func (tController TransferController) FetchTransfers(w http.ResponseWriter, r *http.Request) {
	accountOriginID := fmt.Sprint(r.Context().Value("subject"))
	transferList, err := tController.tUseCase.FetchTransfers(r.Context(), accountOriginID)

	if err != nil {
		reponses.HandleError(w, err, tController.log)
		return
	}

	if len(transferList) > 0 {
		reponses.SendResponse(w, http.StatusOK, transferList, tController.log)

	} else {
		noTransfers := &map[string]string{"msg": "no transfers"}
		reponses.SendResponse(w, http.StatusOK, noTransfers, tController.log)
	}
}
