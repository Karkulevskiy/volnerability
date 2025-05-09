package levels

import (
	"context"
	"volnerability-game/internal/db"
	"volnerability-game/internal/domain"
)

func Run(ctx context.Context, storage *db.Storage, levelId int, input string) (string, error) {
	level, err := storage.Level(ctx, levelId)
	if err != nil {
		return "", err
	}
	if err := storage.IsQueryValid(input); err != nil {
		if db.IsSyntaxError(err) {
			return "invalid SQL syntax", nil
		}
		return "internal err", err
	}
	if level.ExpectedInput != input {
		return "not expected sql injection", nil
	}
	return "accepted", nil
}

func Submit(ctx context.Context, userLevel domain.UserLevel) error {

	return nil
}
