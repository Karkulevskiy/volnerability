package coderunner

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"
	"volnerability-game/internal/common"
	"volnerability-game/internal/db"
	"volnerability-game/internal/domain"
)

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

func (r *CodeRunner) NewTask(ctx context.Context, db *db.Storage, input, reqId string, levelId int) func(context.Context) (domain.Response, error) {
	return func(ctx context.Context) (domain.Response, error) {
		const op = "newTask.codeRunner.creatingTask"

		input, err := wrapCode(input, levelId)
		if err != nil {
			return domain.Response{}, fmt.Errorf("%s: %w", op, err)
		}

		task := domain.Task{
			Code:  input,
			ReqId: reqId,
			Resp:  make(chan domain.ExecuteResponse, 1),
		}
		defer close(task.Resp)

		ctx, cancel := context.WithTimeout(ctx, time.Second*10)
		defer cancel()

		r.queue <- task

		for {
			select {
			case resp := <-task.Resp:
				return handleResp(resp.Resp, levelId)
			case <-ctx.Done():
				r.l.Info(fmt.Sprintf("task runtime exceeded, reqId: %s", task.ReqId))
				return domain.Response{}, fmt.Errorf("task runtime exceeded, reqId: %s", task.ReqId)
			}
		}
	}
}

func wrapCode(input string, levelId int) (string, error) {

	switch levelId {
	case 2:
		return fmt.Sprintf(`
import resource

soft, hard = 100 * 1024, 100 * 1024
resource.setrlimit(resource.RLIMIT_AS, (soft, hard))

%s
		`, input), nil
	}

	return "", fmt.Errorf("invalid level id while wrapping code")
}

func handleResp(output string, levelId int) (domain.Response, error) {
	switch levelId {
	// TODO поставить нужный айди уровня
	case 1:
		if strings.Contains(output, "RecursionError: maximum recursion depth exceeded") {
			return domain.NewResponseOK(), nil
		}
	case 2:
		// 	import resource
		//
		// # Ограничим память до 100 MB
		// soft, hard = 100 * 1024 * 1024, 100 * 1024 * 1024
		// resource.setrlimit(resource.RLIMIT_AS, (soft, hard))
		//
		// # Пример: создаем большой список — MemoryError
		// a = []
		// while True:
		//     a.append("X" * 10**6)
		if strings.Contains(output, "MemoryError") {
			return domain.NewResponseOK(), nil
		}

	default:
		return domain.Response{}, fmt.Errorf("invalid level id")
	}
	return domain.NewResponseBadRequest("failed to overflow stack"), nil
}

func handleCmp(input, expectedInput string) bool {
	toRemove := []string{" ", ""}
	first := common.Remove(strings.Split(input, " "), toRemove...)
	second := common.Remove(strings.Split(expectedInput, " "), toRemove...)
	return slices.Equal(first, second)
}
