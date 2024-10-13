package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"

	"github.com/higordasneves/e-corp/pkg/domain"
	"github.com/higordasneves/e-corp/pkg/gateway/controller/reponses"
)

// Authenticate validates the session token provided as input.
// Tokens are generated by the Login use case and contain an expiration time.
// Returns ErrTokenInvalid if the token does not match or if the token has expired.
func Authenticate(secretKey string, next http.HandlerFunc, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(header) != 2 {
			reponses.HandleError(w, domain.ErrUnauthorized, log)
			return
		}

		tokenString := header[1]
		token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(secretKey), nil
		})

		if err != nil {
			reponses.HandleError(w, domain.ErrUnauthorized, log)
			return
		}

		claims, ok := token.Claims.(*jwt.StandardClaims)
		if !(ok && token.Valid) {
			reponses.HandleError(w, domain.ErrUnauthorized, log)
			return
		}

		ctxWithValue := context.WithValue(r.Context(), "subject", claims.Subject)

		next(w, r.WithContext(ctxWithValue))
	}
}
