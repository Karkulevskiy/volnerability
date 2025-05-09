package common

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken         = fmt.Errorf("invalid token or claims")
	ErrInvalidSigningMethod = fmt.Errorf("invalid signing method")
)

func ParseToken(tokenString string, appSecret []byte) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// Проверка метода подписи
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, ErrInvalidSigningMethod
		}
		return appSecret, nil
	})

	if err != nil {
		if errors.Is(err, ErrInvalidSigningMethod) {
			return nil, ErrInvalidSigningMethod
		}
		return nil, ErrInvalidToken
	}

	// Извлечение claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}
