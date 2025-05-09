package submit

import (
	"errors"
	"log/slog"
	"net/http"
	coderunner "volnerability-game/internal/codeRunner"
	"volnerability-game/internal/common"
	"volnerability-game/internal/db"
	"volnerability-game/internal/levels"
	"volnerability-game/internal/lib/api"
	"volnerability-game/internal/lib/logger/utils"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

func New(l *slog.Logger, db *db.Storage, codeRunner *coderunner.CodeRunner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "rest.submit.New"
		ctx := r.Context()
		reqId := ctx.Value(middleware.RequestIDKey).(string)
		l = l.With(
			slog.String("op", op),
			slog.String("request_id", reqId),
		)

		req := levels.Request{Id: reqId}
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			l.Error("failed to parse request body", utils.Err(err))
			render.JSON(w, r, api.BadRequest(err.Error()))
			return
		}

		l.Info("request body decoded", slog.Any("request", req))

		submit, err := levels.New(req, db, codeRunner)
		if err != nil {
			l.Info("failed to create task", utils.Err(err))
			if common.IsValidateErr(err) {
				render.JSON(w, r, api.BadRequest(err.Error()))
				return
			}
			render.JSON(w, r, api.InternalError(api.WithMsg("oops, smth went wrong :(")))
			return
		}

		resp, err := submit(ctx)
		if err != nil {
			l.Info("error while running submit", utils.Err(err))
			if errors.Is(err, common.ErrBadSubmit) {
				render.JSON(w, r, api.BadSubmit(err.Error()))
				return
			}
			render.JSON(w, r, api.InternalError(api.WithMsg("oops, smth went wrong :(")))
			return
		}

		render.JSON(w, r, resp)
	}
}
