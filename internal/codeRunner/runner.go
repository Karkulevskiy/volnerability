package coderunner

import (
	"context"
	"fmt"
	"log/slog"
	"volnerability-game/internal/domain"
)

const executionTimeout = 5

type Runner interface {
	Run(code string) (string, error)
}

type CodeRunner struct {
	dir   string
	queue chan domain.Task
	l     *slog.Logger
}

func New(l *slog.Logger, dir string, queue chan domain.Task) *CodeRunner {
	return &CodeRunner{
		dir:   dir,
		queue: queue,
		l:     l,
	}
}

func (r *CodeRunner) Run(code, lang, reqId string) (string, error) {
	// TODO нужно быстро уметь валидировать, что код вообще билдиться
	// Механизм кеширования, используя очередь
	task := domain.Task{
		Code:  code,
		Lang:  lang,
		ReqId: reqId,
		Resp:  make(chan domain.ExecuteResponse, 1),
	}
	defer close(task.Resp)

	ctx, cancel := context.WithTimeout(context.Background(), executionTimeout)
	defer cancel()
	// TODO или тут прокидывать контекст с дедлайном, прям в таску?
	// Тут утечка горутин будет 100 %, нужен дедлайн
	for {
		select {
		case r.queue <- task:
		case resp := <-task.Resp:
			return resp.Resp, nil
		case <-ctx.Done():
			r.l.Info(fmt.Sprintf("task runtime exceeded, reqId: %s", task.ReqId)) // TODO не хватает данных айди таски, чтобы потом по логам можно было нормально найти
			return "", fmt.Errorf("task runtime exceeded")
		}
	}
}
