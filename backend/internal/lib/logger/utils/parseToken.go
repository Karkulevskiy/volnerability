package utils

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func ParseToken(tokenString string, appSecret []byte) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// Проверка метода подписи
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return appSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// Извлечение claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token or claims")
}
