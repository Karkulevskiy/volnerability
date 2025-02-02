package auth

import (
	"log/slog"
	"net/http"
	"volnerability-game/internal/lib/logger/utils"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// TODO имплементировать логику логина, регистрации
type Auther interface {
	Login() error
	Register() error
}

type Request struct {
	// TODO Придумать поля для запрос на авторизацию / регистрацию
}

type Response struct {
}

func New(l *slog.Logger, auth Auther) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l.With(
			slog.String("op", "rest.auth.New"), // Задаем тип операции (чтобы это отображалось в логах)
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		req := Request{}
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			l.Error("failed parse request body", utils.Err(err))
			render.JSON(w, r, err) // TODO не стоит просто отправлять внутреннюю ошибку пользователю, нужно ее замаппить на кастомную
			return
		}

		l.Info("request body decoded", slog.Any("request", req))

		// TODO validate credentials

		render.JSON(w, r, http.StatusOK)
	}
}
