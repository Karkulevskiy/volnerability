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
	"volnerability-game/internal/cfg"
	//coderunner "volnerability-game/internal/codeRunner"
	containermgr "volnerability-game/internal/containerMgr"
	grpcmgr "volnerability-game/auth/app/grpc"
	"volnerability-game/internal/db"
	"volnerability-game/internal/lib/logger"
	"volnerability-game/internal/lib/logger/utils"
	authservice "volnerability-game/auth/services"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/mattn/go-sqlite3" // SQLite драйвер
)

func main() {
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

	if err := db.Init(cfg.StoragePath); err != nil {
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

	//codeRunner := coderunner.New(orchestrator.Dir)

	//запуск grpcApp
	authSerivce := authservice.New(l, db, db, cfg.TokenTTL, "12345") // TODO поменять JWTSecret
	grpcSrv := grpcmgr.New(l, *authSerivce, cfg.GRPCCfg.Port)
	go grpcSrv.MustRun()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(logger.New(l))
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	//codeRunner = coderunner.New(l, orchestrator.Queue)

	// r.Post("/code", code.New(l, codeRunner))

	l.Info("starting server", slog.String("address", cfg.Address))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      r,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
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

	//TODO close db
}
