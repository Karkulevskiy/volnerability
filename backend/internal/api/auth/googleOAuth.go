package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"volnerability-game/internal/cfg"
	"volnerability-game/internal/db"

	"github.com/go-chi/chi/middleware"
	"golang.org/x/oauth2"
)

type GoogleUserInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type GoogleAuther struct {
	GoogleOauthCfg *oauth2.Config
}

func NewGoogleAuther(cfg *cfg.Config) *GoogleAuther {
	return &GoogleAuther{
		GoogleOauthCfg: cfg.GetGoogleOAuthConfig(),
	}
}

func (a *GoogleAuther) GoogleAuthHandler(w http.ResponseWriter, r *http.Request) {
	url := a.GoogleOauthCfg.AuthCodeURL("state")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (a *GoogleAuther) NewGoogleAuthCallbackHandler(l *slog.Logger, db *db.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		const op = "rest.auth.GoogleAuthCallbackHandler"
		reqId := ctx.Value(middleware.RequestIDKey).(string)
		l = l.With(
			slog.String("op", op),
			slog.String("request_id", reqId),
		)

		// Получаем код авторизации из query параметров
		code := r.URL.Query().Get("code")

		// Обмениваем код на access token
		token, err := a.GoogleOauthCfg.Exchange(ctx, code)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to exchange token: %v", err)
			l.Error(errMsg)
			http.Error(w, "Failed to authenticate", http.StatusInternalServerError)
			return
		}

		// Создаем HTTP клиент с полученным токеном
		client := a.GoogleOauthCfg.Client(ctx, token)

		// Запрашиваем данные пользователя
		resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
		if err != nil {
			errMsg := fmt.Sprintf("Failed to get user info: %v", err)
			l.Error(errMsg)
			http.Error(w, "Failed to get user info", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Читаем и парсим ответ
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to read response: %v", err)
			l.Error(errMsg)
			http.Error(w, "Failed to read response", http.StatusInternalServerError)
			return
		}

		var userInfo GoogleUserInfo
		if err := json.Unmarshal(data, &userInfo); err != nil {
			errMsg := fmt.Sprintf("Failed to parse user info: %v", err)
			l.Error(errMsg)
			http.Error(w, "Failed to parse user info", http.StatusInternalServerError)
			return
		}

		user, err := db.User(ctx, userInfo.Email)
		if err == sql.ErrNoRows {
			_, err := db.SaveUser(ctx, userInfo.Email, []byte{}, true)
			if err != nil {
				http.Error(w, "Failed to register user", http.StatusInternalServerError)
			}
		}

		if !user.IsOauth {

		}
		// Здесь вы должны:
		// 1. Создать/найти пользователя в вашей БД
		// 2. Сгенерировать JWT токен для вашего приложения
		// 3. Вернуть токен пользователю (например, через куки или redirect с токеном в URL)

		// Пример ответа (в реальном приложении используйте JWT)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Authenticated successfully",
			"user":    userInfo,
		})
	}
}
