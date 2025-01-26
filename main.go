package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"volnerability-game/application/logger"
	"volnerability-game/cfg"
	"volnerability-game/db"
	"volnerability-game/rest/auth"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	// open file for writting logs
	logFile, err := os.OpenFile("game.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	cfg := cfg.MustLoad() // parse cfg in cfg.json

	l := slog.New(slog.NewJSONHandler(logFile, nil)) // sutup logger

	db, err := db.New(cfg.StoragePath)
	if err != nil {
		l.Error("failed to init storage", err.Error())
		os.Exit(1)
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(logger.New(l))
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	// routing
	r.Route("", func(r chi.Router) {
		// TODO r.Use() add basic auth ??
		r.Post("/login", auth.New(l, db)) // TODO логин
		r.Post("/register", nil)          // TODO регистрация

		r.Route("/int", func(r chi.Router) {
			r.Post("/int", nil) // TODO сдесь будут имплементированы внутренние методы
		})
	})

	l.Info("starting server", slog.String("address", cfg.Address))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      r,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			l.Error("failed to start server")
		}
	}()

	l.Info("server started")

	<-done

	l.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		l.Error("failed to stop server: ", err) // TODO проверить, есть ли разница, если вместо err использовать err.Error()
	}

	//TODO close db

	l.Info("server stopped")
}
