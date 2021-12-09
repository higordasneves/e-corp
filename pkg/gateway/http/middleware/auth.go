package middleware

import (
	"context"
	"fmt"
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/gateway/http/controller/responses"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func Authenticate(authUC usecase.AuthUseCase, next http.HandlerFunc, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(header) != 2 {
			responses.SendError(w, http.StatusUnauthorized, usecase.ErrTokenFormat, log)
			return
		}

		tokenString := header[1]
		claims, err := authUC.ValidateToken(tokenString)
		if err != nil {
			responses.SendError(w, http.StatusUnauthorized, fmt.Errorf("%w: %s", usecase.ErrTokenInvalid, err), log)
			return
		}

		ctxWithValue := context.WithValue(r.Context(), "subject", claims.Subject)
		next(w, r.WithContext(ctxWithValue))
	}
}
