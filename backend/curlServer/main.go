package curlServer

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"
	"volnerability-game/internal/domain"
	"volnerability-game/internal/lib/logger/utils"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

func newResponseOK(levelId int) domain.Response {
	return domain.NewResponseOK(domain.WithCurlLevelId(levelId))
}

func AbountHandler(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, newResponseOK(1))
}

func GetDbScheme(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("user")
	if user != "admin" {
		render.JSON(w, r, domain.NewResponseBadRequest("unknown user role"))
		return
	}
	render.JSON(w, r, newResponseOK(2))
}

func Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		render.JSON(w, r, domain.NewResponseBadRequest(fmt.Sprintf("failed to parse form: %s", err.Error())))
	}
	user := r.Form.Get("user")
	password := r.Form.Get("password")
	if user == "admin" && password != "" {
		render.JSON(w, r, newResponseOK(3))
		return
	}
	render.JSON(w, r, domain.NewResponseBadRequest("invalid credentials"))
}

func FilesCmd(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		render.JSON(w, r, domain.NewResponseBadRequest(fmt.Sprintf("failed to parse form: %s", err.Error())))
	}
	cmd := r.Form.Get("cmd")
	if cmd == "ls" {
		render.JSON(w, r, newResponseOK(4))
		return
	}
	fmt.Printf("INPUT: %s\n", cmd)
	if cmd == "cat db.sql" {
		render.JSON(w, r, newResponseOK(5))
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
	r.Get("/db/about", GetDbScheme)
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
