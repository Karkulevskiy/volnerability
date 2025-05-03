package sqllevels

import (
	"context"
	"volnerability-game/internal/db"
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
