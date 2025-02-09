package auth

import (
	"log/slog"
	"net/http"
	"volnerability-game/internal/lib/logger/utils"

	authv1 "volnerability-game/protos/gen/auth"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TODO имплементировать логику логина, регистрации
type Auther interface {
	Login(ctx context.Context, email string, password string) (token string, err error)
	Register(ctx context.Context, email string, password string) (UserID int64, err error)
}

type serverApi struct {
	authv1.UnimplementedAuthServer
	auth Auther
}

func Register(gRPC *grpc.Server, auth Auther) {
	authv1.RegisterAuthServer(gRPC, &serverApi{})
}

func (s *serverApi) Login(ctx context.Context, req *authv1.LoginRequest) (res *authv1.LoginResponse, err error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	if token, err:= s.auth.Login(ctx, req.GetEmail(), req.GetPassword()); err!=nil{

	}

	return &authv1.LoginResponse{
		Token: "",
	}, nil
}

func (s *serverApi) Register(ctx context.Context, req *authv1.RegisterRequest) (res *authv1.RegisterResponse, err error) {
	panic("implement me")
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
