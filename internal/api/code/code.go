package code

import (
	"fmt"
	"log/slog"
	"net/http"
	"volnerability-game/internal/lib/logger/utils"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	Code string
	Lang string
}

type Runner interface {
	Run(code, lang string) (string, error) // TODO переделать респонс
}

func New(l *slog.Logger, runner Runner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l.With(
			slog.String("op", "rest.code.New"), // Задаем тип операции (чтобы это отображалось в логах)
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		req := Request{}
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			l.Error("failed parse request body", utils.Err(err))
			render.JSON(w, r, err)
			return
		}

		l.Info("request body decoded", slog.Any("request", req))

		if err := validate(req); err != nil {
			l.Error("invalid request", utils.Err(err))
			render.JSON(w, r, err)
			return
		}

		resp, err := runner.Run(req.Code, req.Lang)
		if err != nil {
			l.Error("failed run code", utils.Err(err))
			render.JSON(w, r, err)
			return
		}

		render.JSON(w, r, resp)
	}
}

func validate(req Request) error {
	if req.Lang != "c" && req.Lang != "py" {
		return fmt.Errorf("unsupported language format")
	}
	return nil
}
