package user

import (
	"fmt"
	"log/slog"
	"net/http"
	"volnerability-game/internal/common"
	"volnerability-game/internal/db"
	"volnerability-game/internal/domain"
	"volnerability-game/internal/lib/api"
	"volnerability-game/internal/lib/logger/utils"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	Email         string `json:"email"`
	PassLevels    int    `json:"passLevels"`
	TotalAttempts int    `json:"totalAttempts"`
}

func New(l *slog.Logger, db *db.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "rest.user.New"
		ctx := r.Context()
		reqId := ctx.Value(middleware.RequestIDKey).(string)
		l = l.With(
			slog.String("op", op),
			slog.String("request_id", reqId),
		)

		userEmail := r.URL.Query().Get("email")
		if valid := common.IsEmailValid(userEmail); !valid {
			l.Error("invalid request", utils.Err(fmt.Errorf("invalid user email: %s", userEmail)))
			render.JSON(w, r, api.New(api.WithCode(http.StatusBadRequest)))
			return
		}

		user, err := db.User(ctx, userEmail)
		if err != nil {
			l.Error(fmt.Sprintf("failed to user by email: %s", userEmail), utils.Err(err))
			render.JSON(w, r, api.InternalError())
			return
		}
		render.JSON(w, r, toResponse(user))
	}
}

func toResponse(u domain.User) Response {
	// skip internal credentials
	return Response{
		Email:         u.Email,
		TotalAttempts: u.TotalAttempts,
		PassLevels:    u.PassLevels,
	}
}
