package controller

import (
	"net/http"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/gateway/controller/requests"
)

type CreateAccountRequest struct {
	// Name represents the name of the customer.
	Name string `json:"name"`
	// Document is the document number of the customer.
	Document string `json:"document"`
	// Secret is the password. Must have at least 8 digits.
	Secret string `json:"secret"`
}

// CreateAccountResponse represents information from a bank account that
// should be returned to the user after tha account creation.
type CreateAccountResponse struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Document string    `json:"document"`
	// Balance represents the balance of the account.
	Balance   int       `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateAccount creates a banking account.
// @Summary Create Account
// @Description Creates a banking account.
// @Description Returns bad request error if:
// @Description - the account name is not filled;
// @Description - the number of characters of the document is not valid;
// @Description - the format of the document is not valid;
// @Description - the number of the characters of the secret is less than the minimum;
// @Description - the account already exists.
// @Tags Accounts
// @Param Body body CreateAccountRequest true "Request body"
// @Accept json
// @Produce json
// @Success 200 {object} CreateAccountResponse "Account created"
// @Failure 400 {object} ErrorResponse "Invalid parameter"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/accounts [POST]
func (accController AccountController) CreateAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateAccountRequest
	if err := requests.ReadRequestBody(r, &req); err != nil {
		HandleError(ctx, w, err)
		return
	}

	ucOutput, err := accController.accUseCase.CreateAccount(r.Context(), usecase.CreateAccountInput{
		Name:     req.Name,
		Document: req.Document,
		Secret:   req.Secret,
	})
	if err != nil {
		HandleError(ctx, w, err)
		return
	}

	response := CreateAccountResponse{
		ID:        ucOutput.Account.ID,
		Name:      ucOutput.Account.Name,
		Document:  ucOutput.Account.Document.String(),
		Balance:   ucOutput.Account.Balance,
		CreatedAt: ucOutput.Account.CreatedAt,
	}

	SendResponse(ctx, w, http.StatusCreated, response)
}
