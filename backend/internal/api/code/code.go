package code

import (
	"fmt"
	"log/slog"
	"net/http"
	"volnerability-game/internal/lib/logger/utils"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	Code  string `json:"code"`
	Lang  string `json:"lang"`
	Level int    `json:"level"`
	// Other fields
}

// Номер задания -> запрос в бд -> cmp req && resp

type Runner interface {
	Run(code, lang, reqId string) (string, error) // TODO переделать респонс
}

func New(l *slog.Logger, runner Runner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "rest.code.New"
		reqId := r.Context().Value(middleware.RequestIDKey).(string)
		l = l.With(
			slog.String("op", op),
			slog.String("request_id", reqId),
		)

		req := Request{}
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			l.Error("failed to parse request body", utils.Err(err))
			render.JSON(w, r, err)
			return
		}

		l.Info("request body decoded", slog.Any("request", req))

		// TODO подумать над валидацией. Вообще стоит обсудить это детальнее
		if err := validate(req); err != nil {
			l.Error("invalid request", utils.Err(err))
			render.JSON(w, r, err)
			return
		}

		resp, err := runner.Run(req.Code, req.Lang, reqId)
		if err != nil {
			l.Error("failed to run code", utils.Err(err))
			render.JSON(w, r, err)
			return
		}

		render.JSON(w, r, resp)
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
// 9) Запрос на получение инфы по текущему уровню

// TODO добавить ошибки 404, 500 и их маппинг
// при ошибке возвращается инфа, которая ничего не говорит...
func validate(req Request) error {
	if req.Lang != "c" && req.Lang != "py" {
		return fmt.Errorf("unsupported language format")
	}
	return nil
}
