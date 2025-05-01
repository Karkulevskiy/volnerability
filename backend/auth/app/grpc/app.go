package grpcapp

import (
	"fmt"
	"log/slog"
	"net"
	authgrpc "volnerability-game/auth/api"
	authservice "volnerability-game/auth/services"

	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	address    string
}

func (app *App) MustRun() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func New(log *slog.Logger, authService authservice.Auth, address string) *App {
	gRPCServer := grpc.NewServer()
	authgrpc.Register(gRPCServer, &authService)

	return &App{
		log,
		gRPCServer,
		address,
	}
}

func (app *App) Run() error {
	log := app.log.With(
		slog.String("op", "grpcapp.Run"),
		slog.String("address", app.address),
	)

	log.Info("starting gRPC server")
	listener, err := net.Listen("tcp", app.address)
	if err != nil {
		return fmt.Errorf("grpcapp.Run: %w", err)
	}

	if err := app.gRPCServer.Serve(listener); err != nil {
		return fmt.Errorf("grpcapp.Run: %w", err)
	}

	return nil
}

func (app *App) Stop() {
	app.log.With(slog.String("op", "grpcapp.Stop")).
		Info("stopping gRPC server", slog.String("port", app.address))

	app.gRPCServer.GracefulStop()
}
