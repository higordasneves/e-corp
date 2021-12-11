package middleware

import (
	"context"
	"fmt"
	"github.com/higordasneves/e-corp/pkg/domain/usecase"
	"github.com/higordasneves/e-corp/pkg/gateway/http/controller/io"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func Authenticate(authUC usecase.AuthUseCase, next http.HandlerFunc, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(header) != 2 {
			io.HandleError(w, io.ErrTokenFormat, log)
			return
		}

		tokenString := header[1]
		claims, err := authUC.ValidateToken(tokenString)
		if err != nil {
			io.HandleError(w, fmt.Errorf("%w: %s", io.ErrTokenInvalid, err), log)
			return
		}

		ctxWithValue := context.WithValue(r.Context(), "subject", claims.Subject)
		next(w, r.WithContext(ctxWithValue))
	}
}
