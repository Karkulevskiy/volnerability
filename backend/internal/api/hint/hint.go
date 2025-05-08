package hint

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"volnerability-game/internal/db"
	"volnerability-game/internal/hints"
	"volnerability-game/internal/lib/api"
	"volnerability-game/internal/lib/logger/utils"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	HintId int `json:"hintId"`
}

func New(l *slog.Logger, db *db.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "rest.hint.New"
		ctx := r.Context()
		reqId := ctx.Value(middleware.RequestIDKey).(string)
		l = l.With(
			slog.String("op", op),
			slog.String("request_id", reqId),
		)

		hintIdStr := r.URL.Query().Get("hintId")
		hintId, err := strconv.Atoi(hintIdStr)
		if err != nil {
			l.Error("failed to parse query", utils.Err(err))
			render.JSON(w, r, api.New(api.WithCode(http.StatusExpectationFailed)))
			return
		}

		if err := validate(hintId); err != nil {
			l.Error("failed to validate request", utils.Err(err))
			render.JSON(w, r, api.InternalError())
			return
		}

		hint, err := hints.Run(ctx, db, hintId)
		if err != nil {
			l.Error("failed to get hint", utils.Err(err))
			render.JSON(w, r, api.InternalError())
			return
		}
		render.JSON(w, r, hint)
	}
}

func validate(hintId int) error {
	if hintId <= 0 {
		return fmt.Errorf("invalid hint id")
	}
	return nil
}
