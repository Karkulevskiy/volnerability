package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	authv1 "volnerability-game/auth/protos/gen/auth"

	"github.com/go-chi/chi/middleware"
)

type Credentials struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type RegisterResponse struct {
	UserID int64 `json:"userId"`
}

func NewLoginHandler(l *slog.Logger, grpcClient authv1.AuthClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		const op = "rest.auth.NewLoginHandler"
		reqId := ctx.Value(middleware.RequestIDKey).(string)
		l = l.With(
			slog.String("op", op),
			slog.String("request_id", reqId),
		)

		pass, email, err := parseCredentials(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid JSON: %v", err.Error()), http.StatusBadRequest)
			return
		}
		req := &authv1.LoginRequest{
			Email:    email,
			Password: pass,
		}

		grpcRes, err := grpcClient.Login(ctx, req)
		if err != nil {
			http.Error(w, fmt.Sprintf("gRPC login call failed: %v", err), http.StatusInternalServerError)
			return
		}
		res := &LoginResponse{
			Token: grpcRes.Token,
		}
		if err = json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, "Failed to encode login response", http.StatusInternalServerError)
			return
		}
	}
}

func NewRegisterHandler(l *slog.Logger, grpcClient authv1.AuthClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		const op = "rest.auth.NewRegisterHandler"
		reqId := ctx.Value(middleware.RequestIDKey).(string)
		l = l.With(
			slog.String("op", op),
			slog.String("request_id", reqId),
		)

		pass, email, err := parseCredentials(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid JSON: %v", err.Error()), http.StatusBadRequest)
			return
		}
		req := &authv1.RegisterRequest{
			Email:    email,
			Password: pass,
		}

		grpcRes, err := grpcClient.Register(ctx, req)
		if err != nil {
			http.Error(w, fmt.Sprintf("gRPC register call failed: %v", err), http.StatusInternalServerError)
			return
		}
		res := &RegisterResponse{
			UserID: grpcRes.UserId,
		}
		if err = json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, "Failed to encode register response", http.StatusInternalServerError)
			return
		}
	}
}

func parseCredentials(body io.ReadCloser) (string, string, error) {
	var creds Credentials
	err := json.NewDecoder(body).Decode(&creds)
	if err != nil {
		return "", "", err
	}
	pass := creds.Password
	email := creds.Email
	return pass, email, nil
}
