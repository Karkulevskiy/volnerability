package level

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"volnerability-game/internal/db"
	"volnerability-game/internal/levels"
	"volnerability-game/internal/lib/api"
	"volnerability-game/internal/lib/logger/utils"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

func New(l *slog.Logger, db *db.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "rest.level.New"
		ctx := r.Context()
		reqId := ctx.Value(middleware.RequestIDKey).(string)
		l = l.With(
			slog.String("op", op),
			slog.String("request_id", reqId),
		)

		levelIdStr := r.URL.Query().Get("id")
		levelId, err := strconv.Atoi(levelIdStr)
		if err != nil {
			l.Error("failed to parse query", utils.Err(err))
			render.JSON(w, r, api.New(api.WithCode(http.StatusExpectationFailed)))
			return
		}

		if err := validate(levelId); err != nil {
			l.Error("failed to validate request", utils.Err(err))
			render.JSON(w, r, api.InternalError())
			return
		}

		level, err := levels.Level(ctx, db, levelId)
		if err != nil {
			l.Error("failed to get level", utils.Err(err))
			render.JSON(w, r, api.InternalError())
			return
		}
		// shadow for user
		level.ExpectedInput = ""
		render.JSON(w, r, level)
	}
}

func validate(hintId int) error {
	if hintId <= 0 {
		return fmt.Errorf("invalid level id")
	}
	return nil
}
