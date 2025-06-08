package curlServer

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"
	"volnerability-game/internal/domain"
	"volnerability-game/internal/lib/logger/utils"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

func AbountHandler(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, domain.NewResponseOK())
}

func GetDbScheme(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("user")
	if user != "admin" {
		render.JSON(w, r, domain.NewResponseBadRequest("unknown user role"))
		return
	}
	render.JSON(w, r, domain.NewResponseOK())
}

func Login(w http.ResponseWriter, r *http.Request) {
	user, password := r.URL.Query().Get("user"), r.URL.Query().Get("password")
	if user == "admin" && password != "" {
		render.JSON(w, r, domain.NewResponseOK())
		return
	}
	render.JSON(w, r, domain.NewResponseBadRequest("invalid credentials"))
}

func FilesCmd(w http.ResponseWriter, r *http.Request) {
	req := domain.Request{}
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		slog.Error("failed to decode body: ")
		render.JSON(w, r, domain.NewResponseBadRequest("failed to decode json body"))
		return
	}

	input := strings.Trim(req.Input, " ")
	if input == "ls" || input == "cat" {
		render.JSON(w, r, domain.NewResponseOK())
		return
	}
	render.JSON(w, r, domain.NewResponseBadRequest("invalid cmd"))
}

// def send_email(name, email, message):
//    email_body = f"From: {email}\nTo: support@example.com\nSubject: Feedback\n\n{message}"
//    smtp.send(email_body)

// curl -X POST "http://example.com/contact" \
//     -d "email=attacker@example.com%0ABcc:admin@example.com" \
//     -d "message=Hello!"

// Тут подразумевается, что пользователь введет такой курл, и отправит сообщение свои коллегам: Привет, я создал новую БД вот скрипт - посмотрите пожалуйста
func Email(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		render.JSON(w, r, domain.NewResponseBadRequest("invalid form"))
		return
	}

	email := r.FormValue("email")
	if email == "" {
		render.JSON(w, r, domain.NewResponseBadRequest("invalid email"))
		return
	}

	message := r.FormValue("message")
	if message == "" {
		render.JSON(w, r, domain.NewResponseBadRequest("invalid message"))
		return
	}

	render.JSON(w, r, domain.NewResponseOK())
}

func MustRun() {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowCredentials: true,
	}))
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	r.Get("/about", AbountHandler)
	r.Post("/db/about", GetDbScheme)
	r.Post("/login", Login)
	r.Post("/files", FilesCmd)

	srv := &http.Server{
		Addr:         "127.0.0.1:9086",
		Handler:      r,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
		IdleTimeout:  time.Second * 5,
	}

	run(srv) // TODO gracefully pohui
}

func run(srv *http.Server) {
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				slog.Info("server stopped")
				return
			}
			slog.Error("failed to start server", utils.Err(err))
		}
	}()
	slog.Info("server started")
}
