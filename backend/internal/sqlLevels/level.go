package sqllevels

import (
	"context"
	"fmt"
	"volnerability-game/internal/db"
)

func Run(ctx context.Context, storage *db.Storage, levelId int, input string) (string, error) {
	level, err := storage.Level(ctx, levelId)
	if err != nil {
		return "", err
	}
	isValid, err := storage.IsQueryValid(input)
	if err != nil {
		if db.IsSyntaxError(err) {
			return "", fmt.Errorf("")
		}
		return "", fmt.Errorf("")
	}
	// TODO Тут прям напрашивается механгизм маппинга ошибок, чтобы по кайфу было
	if !isValid {
		return "", nil
	}
	if level.ExpectedInput != input {
		return "", nil
	}
	return "", nil
}
