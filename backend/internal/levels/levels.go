package levels

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	coderunner "volnerability-game/internal/codeRunner"
	"volnerability-game/internal/common"
	"volnerability-game/internal/curl"
	"volnerability-game/internal/db"
	"volnerability-game/internal/domain"
	sqlrunner "volnerability-game/internal/sqlRunner"
)

type Submit func(ctx context.Context) (domain.Response, error)

func New(ctx context.Context, r domain.Request, db *db.Storage, codeRunner *coderunner.CodeRunner) (Submit, error) {
	if err := validate(r); err != nil {
		return nil, err
	}
	if r.LevelId <= 5 { // [1-4] curl tasks
		return curl.NewTask(db, r.LevelId, r.Input), nil
	}
	if r.LevelId <= 8 { // [5-7] sql tasks
		return sqlrunner.NewTask(db, r.LevelId, r.Input), nil
	}
	if r.LevelId <= 10 { // [8-9]code tasks
		return codeRunner.NewTask(ctx, db, r.Input, r.Id, r.LevelId), nil
	}
	return nil, nil
}

func ProcessTask(ctx context.Context, db *db.Storage, req domain.Request, task Submit) (domain.Response, error) {
	const op = "internal.levels.ProcessTask"

	email, ok := ctx.Value("email").(string)
	if !ok {
		return domain.Response{}, fmt.Errorf("failed to get user email")
	}

	slog.Info(fmt.Sprintf("user email: %s", email)) // TODO просто для тестов

	user, err := db.User(ctx, email)
	if err != nil {
		return domain.Response{}, fmt.Errorf("%s: failed to get user by email: %s, due err: %w", op, email, err)
	}

	resp, err := task(ctx)
	if err != nil {
		return domain.Response{}, fmt.Errorf("%s: failed to do task, due err: %w", op, err)
	}

	if err := updateUserAttempt(ctx, db, user, req, resp); err != nil {
		return domain.Response{}, fmt.Errorf("%s: failed to update user attempt due err: %w", op, err)
	}

	return resp, nil
}

func updateUserAttempt(ctx context.Context, storage *db.Storage, user domain.User, req domain.Request, resp domain.Response) error {
	const op = "levels.updateUserAttempt"
	user.TotalAttempts++
	updatedAttempt := domain.UserLevel{
		UserId:          int(user.ID),
		LevelId:         req.LevelId,
		IsCompleted:     resp.IsCompleted,
		LastInput:       req.Input,
		AttemptResponse: resp.Response,
		Attempts:        user.TotalAttempts,
	}

	userCurrentLevel := user.PassLevels
	lastPassedLevelShouldBe := req.LevelId - 1

	if userCurrentLevel < lastPassedLevelShouldBe {
		return fmt.Errorf("%s: last user passed level: %d, request level: %d", op, userCurrentLevel, lastPassedLevelShouldBe)
	}

	userLevel, err := storage.UserLevel(ctx, user.Email, req.LevelId)
	if err != nil {
		// ранее уровень не был пройден, первая попытка
		if errors.Is(err, db.ErrUserLevelNotFound) {
			if resp.IsCompleted {
				user.PassLevels++
			}
			return storage.UpdateUserAndLevel(ctx, user, updatedAttempt)
		}
		return err
	}

	// Уровень был пройден ранее, тогда не будем обновлять инфу по уровню
	if userLevel.IsCompleted {
		return nil
	}

	// Ранее уровень не пройден и пользователь прошел уровень
	if !userLevel.IsCompleted && resp.IsCompleted {
		user.PassLevels++
	}

	return storage.UpdateUserAndLevel(ctx, user, updatedAttempt)
}

func Level(ctx context.Context, db *db.Storage, levelId int) (domain.Level, error) {
	const op = "level.getLevelById"
	level, err := db.Level(ctx, levelId)
	if err != nil {
		return domain.Level{}, fmt.Errorf("op: %s. failed to proceed get level by id: %d", op, levelId)
	}
	return level, nil
}

func validate(r domain.Request) error {
	if r.LevelId <= 0 || r.LevelId > common.MaxLevel {
		return common.ErrInvalidLevelId
	}
	if r.Input == "" {
		return common.ErrEmptyInput
	}
	return nil
}
