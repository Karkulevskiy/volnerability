package user

import (
	"log/slog"
	"net/http"
	"os"
	"volnerability-game/auth/lib/jwt"
	"volnerability-game/internal/cfg"
	"volnerability-game/internal/db"
	models "volnerability-game/internal/domain"
	"volnerability-game/internal/lib/api"
	"volnerability-game/internal/lib/logger/utils"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
)

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type ChangePasswordResponse struct {
	Token string `json:"token"`
}

func ChangePassword(l *slog.Logger, db *db.Storage, cfg *cfg.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const emailClaim = "email"
		const op = "rest.user.ChangePassword"
		ctx := r.Context()
		reqId := ctx.Value(middleware.RequestIDKey).(string)
		l = l.With(
			slog.String("op", op),
			slog.String("request_id", reqId),
		)

		email, ok := ctx.Value(emailClaim).(string)
		if !ok || email == "" {
			l.Error("failed to get email by JWT")
			render.JSON(w, r, api.BadRequest("wrong JWT token"))
			return
		}

		user, err := db.User(ctx, email)
		if err != nil {
			l.Error("failed to get user", utils.Err(err))
			render.JSON(w, r, api.BadRequest(err.Error()))
			return
		}

		req := ChangePasswordRequest{}
		if err = render.DecodeJSON(r.Body, &req); err != nil {
			l.Error("unable to decode json for changing password", utils.Err(err))
			render.JSON(w, r, api.BadRequest("wrong request"))
			return
		}

		if err = checkOldPassword(user, req.OldPassword); err != nil {
			l.Error("wrong old password", utils.Err(err))
			render.JSON(w, r, api.BadRequest("wrong old password"))
			return
		}

		passHash, err := cryptNewPassword(req.NewPassword)
		if err != nil {
			l.Error("unable to crypt new password", utils.Err(err))
			render.JSON(w, r, api.BadRequest(err.Error()))
			return
		}

		newUser := models.User{PassHash: passHash, Email: user.Email}
		if err = db.UpdateUser(ctx, newUser); err != nil {
			l.Error("unable to change password", utils.Err(err))
			render.JSON(w, r, api.BadRequest("unable to change password"))
			return
		}

		user.PassHash = passHash
		newToken, err := jwt.NewToken(user, cfg.TokenTTL, os.Getenv("JWT_SECRET"))
		if err != nil {
			l.Error("failed to generate token", utils.Err(err))
			render.JSON(w, r, api.BadRequest("failed to submit password changing"))
			return
		}

		resp := ChangePasswordResponse{Token: newToken}
		render.JSON(w, r, resp)

	}
}

func checkOldPassword(user models.User, password string) error {
	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		return err
	}

	return nil
}

func cryptNewPassword(password string) ([]byte, error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return passHash, nil
}
