package jwt

import (
	"time"

	"volnerability-game/internal/domain/models"

	"github.com/golang-jwt/jwt/v5"
)

func NewToken(user models.User, duration time.Duration, Secret string) (string, error) {
	claims := jwt.MapClaims{
		"uid":   user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(Secret)) 
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
