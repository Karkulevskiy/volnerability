package coderunner

import (
	"context"
	"fmt"
	"log/slog"
	"time"
	"volnerability-game/internal/domain"
)

type Runner interface {
	Run(code string) (string, error)
}

type CodeRunner struct {
	queue chan domain.Task
	l     *slog.Logger
}

func New(l *slog.Logger, queue chan domain.Task) *CodeRunner {
	return &CodeRunner{
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	r.queue <- task

	for {
		select {
		case resp := <-task.Resp:
			return resp.Resp, nil
		case <-ctx.Done():
			r.l.Info(fmt.Sprintf("task runtime exceeded, reqId: %s", task.ReqId)) // TODO не хватает данных айди таски, чтобы потом по логам можно было нормально найти
			return "", nil                                                        // TODO создать ошибку с типом, что таска не успела выполниться, и прокидывать дальше
		}
	}
}
