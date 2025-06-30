package token

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mrtuuro/driver-location-api/internal/apperror"
	"github.com/mrtuuro/driver-location-api/internal/code"
)

type TokenManager struct {
	Secret string
}

func NewTokenManager(secret string) *TokenManager {
	return &TokenManager{Secret: secret}
}

type TokenClaims struct {
	Authenticated bool
	jwt.RegisteredClaims
}

func (tm *TokenManager) ValidateJWT(tokenStr string) (*TokenClaims, error) {

	token, err := jwt.ParseWithClaims(
		tokenStr,
		&TokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(tm.Secret), nil
		})
	if err != nil || !token.Valid {
		fmt.Println("here 2")
		return nil, apperror.NewAppError(
			code.ErrAuthInvalidToken,
			errors.New(code.ErrAuthInvalidToken),
			code.GetErrorMessage(code.ErrAuthInvalidToken),
			)
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, apperror.NewAppError(
			"TOKEN_INVALID_CLAIMS",
			errors.New("Invalid claims"),
			"Claims are not valid.",
			)
	}

	return claims, nil
}
