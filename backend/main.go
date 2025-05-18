package main

import (
	"context"
	"errors"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	grpcmgr "volnerability-game/auth/app/grpc"
	authservice "volnerability-game/auth/services"
	"volnerability-game/internal/api/auth"
	"volnerability-game/internal/api/hint"
	"volnerability-game/internal/api/submit"
	"volnerability-game/internal/api/user"
	"volnerability-game/internal/cfg"
	coderunner "volnerability-game/internal/codeRunner"
	containermgr "volnerability-game/internal/containerMgr"
	"volnerability-game/internal/db"
	"volnerability-game/internal/lib/logger"
	"volnerability-game/internal/lib/logger/utils"
	cstmMiddleware "volnerability-game/internal/middleware"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Fatal(err)
		}
		log.Println(".env not exists, auth not working")
	}

	logFile, err := os.OpenFile("game.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	cfg := cfg.MustLoad()

	l := slog.New(slog.NewTextHandler(
		io.MultiWriter(logFile, os.Stdout),
		&slog.HandlerOptions{AddSource: true}))

	db, err := db.New(cfg.StoragePath)
	if err != nil {
		panic(err)
	}

	orchestrator, err := containermgr.New(l, cfg.OrchestratorConfig)
	if err != nil {
		panic(err)
	}

	defer func() {
		l.Info("stopping containers")
		if err := orchestrator.Stop(); err != nil {
			l.Error("failed stop containers", utils.Err(err))
		}
	}()

	l.Info("start containers")
	if err := orchestrator.RunContainers(); err != nil {
		panic(err)
	}
	l.Info("containers started")

	appSecret := os.Getenv("JWT_SECRET")

	authSerivce := authservice.New(l, db, db, time.Duration(cfg.TokenTTL), appSecret)
	grpcSrv, err := grpcmgr.New(l, *authSerivce, cfg.GRPCConfig.Address)
	if err != nil {
		panic(err)
	}
	go grpcSrv.MustRun()
	l.Info("auth server started", slog.String("address", cfg.GRPCConfig.Address))
	grpcClnt := grpcSrv.GetGRPCClient()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(logger.New(l))
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://*"},
	}))
	r.Use(cstmMiddleware.New(l, appSecret))

	codeRunner := coderunner.New(l, orchestrator.Queue)

	r.Post("/submit", submit.New(l, db, codeRunner))
	r.Post("/login", auth.NewLoginHandler(l, grpcClnt))
	r.Post("/register", auth.NewRegisterHandler(l, grpcClnt))
	r.Get("/hint", hint.New(l, db)) // чисто для подсказок
	r.Get("/user", user.New(l, db))

	l.Info("starting server", slog.String("address", cfg.HttpServer.Address))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.HttpServer.Address,
		Handler:      r,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				l.Info("server stopped")
				return
			}
			l.Error("failed to start server", utils.Err(err))
		}
	}()
	l.Info("server started")

	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		l.Error("failed to stop server: ", utils.Err(err))
		return
	}

	if err := db.Close(); err != nil {
		l.Error("failed to close db: ", utils.Err(err))
		return
	}
}
