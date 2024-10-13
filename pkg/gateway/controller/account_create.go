package controller

import (
	"net/http"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/gateway/controller/reponses"
	"github.com/higordasneves/e-corp/pkg/gateway/controller/requests"
)

type CreateAccountRequest struct {
	Name     string `json:"name"`
	Document string `json:"document"`
	Secret   string `json:"secret"`
}

// CreateAccountResponse represents information from a bank account that
// should be returned to the user after tha account creation.
type CreateAccountResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Document  string    `json:"document"`
	Balance   int       `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateAccount creates a banking account.
// Returns BadRequest error if:
// - the account name is not filled;
// - the number of characters of the document is not valid;
// - the format of the document is not valid;
// - the number of the characters of the secret is less than the minimum;
// - the account already exists.
func (accController AccountController) CreateAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateAccountRequest
	if err := requests.ReadRequestBody(r, &req); err != nil {
		reponses.HandleError(ctx, w, err)
		return
	}

	ucOutput, err := accController.accUseCase.CreateAccount(r.Context(), usecase.CreateAccountInput{
		Name:     req.Name,
		Document: req.Document,
		Secret:   req.Secret,
	})
	if err != nil {
		reponses.HandleError(ctx, w, err)
		return
	}

	response := CreateAccountResponse{
		ID:        ucOutput.Account.ID,
		Name:      ucOutput.Account.Name,
		Document:  ucOutput.Account.Document.String(),
		Balance:   ucOutput.Account.Balance,
		CreatedAt: ucOutput.Account.CreatedAt,
	}

	reponses.SendResponse(ctx, w, http.StatusCreated, response)
}
