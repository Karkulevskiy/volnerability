package code

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	Code string
}

type Runner interface {
	Run(code string) (string, error) // TODO переделать респонс
}

func New(l *slog.Logger, runner Runner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l.With(
			slog.String("op", "rest.code.New"), // Задаем тип операции (чтобы это отображалось в логах)
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		req := Request{}
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			l.Error("failed parse request body", err.Error())
			render.JSON(w, r, err) // TODO не стоит просто отправлять внутреннюю ошибку пользователю, нужно ее замаппить на кастомную
			return
		}

		l.Info("request body decoded", slog.Any("request", req))


		resp, err := runner.Run(req.Code)
		if err != nil {
			l.Error("failed run code", err.Error())
			render.JSON(w, r, err)
		}

		render.JSON(w, r, resp)
	}
}
