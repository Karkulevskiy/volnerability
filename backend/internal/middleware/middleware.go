package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"volnerability-game/internal/lib/logger/utils"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

func New(log *slog.Logger, appSecret string) func(next http.Handler) http.Handler {
	const op = "middleware.New"
	const errKey = "parse token"
	const emailClaim = "email"
	log = log.With(slog.String("op", op))

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := extractBearerToken(r)
			if tokenStr == "" {
				next.ServeHTTP(w, r)
				return
			}

			// hmac alg. expects type of []byte, not string. NOT fix this!
			claims, err := utils.ParseToken(tokenStr, []byte(appSecret))
			if err != nil {
				log.Warn("failed to parse token", utils.Err(err))

				ctx := context.WithValue(r.Context(), errKey, ErrInvalidToken)
				next.ServeHTTP(w, r.WithContext(ctx))

				return
			}

			// прокинул дальше почту пользователя
			userEmail := claims[emailClaim]
			ctx := context.WithValue(r.Context(), emailClaim, userEmail)
			log.Info("user authorized", slog.Any("claims", claims))

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// В случае HTTP-запросов, JWT-токен обычно отправляют в заголовке вида:
// Authorization: "Bearer <jwt_token>"
func extractBearerToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 {
		return ""
	}
	return splitToken[1]
}
