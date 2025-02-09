package grpcapp

import (
	"fmt"
	"log/slog"
	"net"
	"strconv"
	authgrpc "volnerability-game/internal/api/auth"

	"google.golang.org/grpc"
)


type App struct{
	log *slog.Logger
	gRPCServer *grpc.Server
	port int
}

func (app *App) MustRun(){
	if err:= app.Run(); err!= nil{
		panic(err)
	}
}

func New(log *slog.Logger, port int) *App{
	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer) //Так она автоматически ребутается с вс кодом

	return &App{
		log,
		gRPCServer,
		port,
	}
}

func (app *App) Run() error{
	log := app.log.With(
		slog.String("op", "grpcapp.Run"),
	 	slog.Int("port", app.port),
	)

	log.Info("starting gRPC server")
	listener, err := net.Listen("tcp", strconv.Itoa(app.port))
	if err!= nil{
		return fmt.Errorf("grpcapp.Run: %w", err)
	}
	 
	if err:= app.gRPCServer.Serve(listener); err!= nil{
		return fmt.Errorf("grpcapp.Run: %w", err)
	}

	return nil
}

func (app *App) Stop(){
	app.log.With(slog.String("op","grpcapp.Stop")).
		Info("stopping gRPC server", slog.Int("port", app.port))

	app.gRPCServer.GracefulStop()
}