package usecase

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

var (
	ErrTokenFormat  = errors.New("invalid token format")
	ErrTokenInvalid = errors.New("invalid token")
)

func (authUC authUseCase) ValidateToken(tokenString string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(authUC.secretKey), nil
	})

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !(ok && token.Valid) {
		return nil, err
	}

	return claims, nil
}
