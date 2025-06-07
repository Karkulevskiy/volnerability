package levels

import (
	"context"
	"fmt"
	"log/slog"
	coderunner "volnerability-game/internal/codeRunner"
	"volnerability-game/internal/common"
	"volnerability-game/internal/db"
	"volnerability-game/internal/domain"
	sqlrunner "volnerability-game/internal/sqlRunner"
)

type Request struct {
	Id      string
	LevelId int    `json:"levelId"`
	Input   string `json:"input"`
}

type Response struct {
	Status   string `json:"status"`
	Response string `json:"response"`
}

type Submit func(ctx context.Context) (string, bool, error)

func New(ctx context.Context, r Request, db *db.Storage, codeRunner *coderunner.CodeRunner) (Submit, error) {
	if err := validate(r); err != nil {
		return nil, err
	}
	// TODO need to create level groups
	if r.LevelId < 8 {
		return sqlrunner.NewTask(db, r.LevelId, r.Input)
	}
	if r.LevelId < 3 {
		return codeRunner.NewTask(ctx, db, r.Input, r.Id, r.LevelId)
	}
	// TODO other tasks
	return nil, nil
}

func ProcessTask(ctx context.Context, db *db.Storage, r Request, task Submit) (Response, error) {
	const op = "internal.levels.ProcessTask"

	email, ok := ctx.Value("email").(string)
	if !ok {
		return Response{}, fmt.Errorf("failed to get user email")
	}

	slog.Info(fmt.Sprintf("user email: %s", email)) // TODO просто для тестов

	user, err := db.User(ctx, email)
	if err != nil {
		return Response{}, fmt.Errorf("%s: failed to get user by email: %s, due err: %w", op, email, err)
	}

	output, isCompleted, err := task(ctx)
	if err != nil {
		return Response{}, fmt.Errorf("%s: failed to do task, due err: %w", op, err)
	}

	if err := updateUserAttempt(ctx, db, user, r, output, isCompleted); err != nil {
		return Response{}, fmt.Errorf("%s: failed to update user attempt due err: %w", op, err)
	}

	return toResponse(output, isCompleted), nil
}

func toResponse(output string, isCompleted bool) Response {
	if isCompleted {
		return Response{Status: common.LevelPassed}
	}
	return Response{Status: common.LevelNotPassed, Response: output}
}

func updateUserAttempt(ctx context.Context, db *db.Storage, user domain.User, r Request, resp string, isCompleted bool) error {
	updatedAttempt := domain.UserLevel{
		UserId:          int(user.ID),
		LevelId:         r.LevelId,
		IsCompleted:     isCompleted,
		LastInput:       r.Input,
		AttemptResponse: resp,
		Attempts:        user.TotalAttempts,
	}
	user.TotalAttempts++
	if isCompleted {
		user.PassLevels++
	}
	return db.UpdateUserAndLevel(ctx, user, updatedAttempt)
}

func Level(ctx context.Context, db *db.Storage, levelId int) (domain.Level, error) {
	const op = "level.getLevelById"
	level, err := db.Level(ctx, levelId)
	if err != nil {
		return domain.Level{}, fmt.Errorf("op: %s. failed to proceed get level by id: %d", op, levelId)
	}
	return level, nil
}

func validate(r Request) error {
	if r.LevelId <= 0 || r.LevelId > common.MaxLevel {
		return common.ErrInvalidLevelId
	}
	if r.Input == "" {
		return common.ErrEmptyInput
	}
	return nil
}
