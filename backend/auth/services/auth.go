package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
	"volnerability-game/auth/lib/jwt"
	"volnerability-game/internal/db"
	models "volnerability-game/internal/domain"
	"volnerability-game/internal/lib/logger/utils"

	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	usrSaver    UserSaver
	usrProvider UserProvider
	log         *slog.Logger
	tokenTTL    time.Duration
	jwtSecret   string
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
}

var (
	ErrInvalidCredentials = fmt.Errorf("invalid credentials")
	ErrUserExists         = fmt.Errorf("user already exists")
)

func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	tokenTTL time.Duration,
	jwtSecret string,
) *Auth {
	return &Auth{
		usrSaver:    userSaver,
		usrProvider: userProvider,
		log:         log,
		tokenTTL:    tokenTTL,
		jwtSecret:   jwtSecret,
	}
}

func (a *Auth) Login(
	ctx context.Context,
	email string,
	password string,
) (string, error) {
	const op = "Auth.LoginUser"
	log := a.log.With(slog.String("op", op))

	log.Info("attempting to login user")

	user, err := a.usrProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, db.ErrUserNotFound) {
			return "", fmt.Errorf("user not found. %s: %w", op, ErrInvalidCredentials)
		}
		return "", fmt.Errorf("failed to get user. %s: %w", op, err)
	}
	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("failed to validate credentials", utils.Err(err))
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	token, err := jwt.NewToken(user, a.tokenTTL, a.jwtSecret)
	if err != nil {
		a.log.Warn("failed to generate token", utils.Err(err))
		return "", fmt.Errorf("failed to generate token. %s: %w", op, err)
	}

	return token, nil
}

func (a *Auth) Register(ctx context.Context, email string, password string) (int64, error) {
	const op = "auth.RegisterNewUser"
	log := a.log.With(slog.String("op", op))

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", utils.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.usrSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, db.ErrUserExists) {
			log.Warn("user already exists", utils.Err(err))

			return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		}
		log.Error("failed to save user", utils.Err(err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
