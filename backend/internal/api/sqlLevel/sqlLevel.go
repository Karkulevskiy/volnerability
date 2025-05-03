package sqllevel

import (
	"fmt"
	"log/slog"
	"net/http"
	"volnerability-game/internal/common"
	"volnerability-game/internal/db"
	"volnerability-game/internal/lib/api"
	"volnerability-game/internal/lib/logger/utils"
	sqllevels "volnerability-game/internal/sqlLevels"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

var (
	ErrInvalidLevelId = fmt.Errorf("invalid level id")
	ErrEmptyInput     = fmt.Errorf("empty input")
)

type Request struct {
	LevelId int    `json:"levelId"`
	Input   string `json:"input"`
	// Other fields
}

func New(l *slog.Logger, db *db.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "rest.sql.New"
		ctx := r.Context()
		reqId := ctx.Value(middleware.RequestIDKey).(string)
		l = l.With(
			slog.String("op", op),
			slog.String("request_id", reqId),
		)

		req := Request{}
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			l.Error("failed to parse request body", utils.Err(err))
			render.JSON(w, r, api.New(api.WithCode(http.StatusExpectationFailed)))
			return
		}

		l.Info("request body decoded", slog.Any("request", req))

		if err := validate(req); err != nil {
			l.Error("invalid request", utils.Err(err))
			render.JSON(w, r, api.New(api.WithCode(http.StatusExpectationFailed)))
			return
		}

		// Добавить sql для подсказок
		// Добавить ui на подсказки
		_, err := sqllevels.Run(ctx, db, req.LevelId, req.Input)
		if err != nil {
			l.Error("failed to run sql level", utils.Err(err))
			render.JSON(w, r, api.InternalError())
			return
		}
		render.JSON(w, r, api.OK())
	}
}

// TODO
// 1) Задание с sql инъекцией
// 2) Маппинг ошибок для пользователя
// 3) middleware auth
// 4) Получение информации о пользователе
// 5) Хранить в куки части фронта
// 6) Может добавить простые задания с выбором вариантов ответов?
// 7) Добавить разделение заданий по группам
// 8) Наверное стоит делать по одному запросу в бд по всех инфе для конкретной группы

func validate(req Request) error {
	if req.LevelId <= 0 || req.LevelId > common.MaxLevel {
		return ErrInvalidLevelId
	}
	if len(req.Input) == 0 {
		return ErrEmptyInput
	}
	return nil
}
