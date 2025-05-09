package sqlrunner

import (
	"context"
	"volnerability-game/internal/common"
	"volnerability-game/internal/db"
)

func NewTask(storage *db.Storage, levelId int, input string) (func(context.Context) (any, error), error) {
	return func(ctx context.Context) (any, error) {
		level, err := storage.Level(ctx, levelId)
		if err != nil {
			return "", err
		}
		if err := storage.IsQueryValid(input); err != nil {
			if db.IsSyntaxError(err) {
				return "", common.NewBadSubmitErr("invalid SQL syntax")
			}
			return "", err
		}
		if level.ExpectedInput != input {
			return "", common.NewBadSubmitErr("not expected sql injection")
		}
		return "", nil
	}, nil
}
