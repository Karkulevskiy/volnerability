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

func handleResp(output string, levelId int) (domain.Response, error) {
	switch levelId {
	case 1:
		if strings.Contains(output, "RecursionError: maximum recursion depth exceeded") {
			return domain.NewResponseOK(), nil
		}
		return domain.NewResponseBadRequest("failed to overflow stack"), nil
	default:
		return domain.Response{}, fmt.Errorf("invalid level id")
	}
}

func handleCmp(input, expectedInput string) bool {
	toRemove := []string{" ", ""}
	first := common.Remove(strings.Split(input, " "), toRemove...)
	second := common.Remove(strings.Split(expectedInput, " "), toRemove...)
	return slices.Equal(first, second)
}
