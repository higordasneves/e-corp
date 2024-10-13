package requests

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/higordasneves/e-corp/pkg/domain"
)

func ReadRequestBody(r *http.Request, obj interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(&obj); err != nil {
		return fmt.Errorf("%w: invalid request body: %s", domain.ErrInvalidParameter, err)
	}

	return nil
}
