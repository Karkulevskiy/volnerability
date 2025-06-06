package sqlrunner

import (
	"context"
	"regexp"
	"volnerability-game/internal/common"
	"volnerability-game/internal/db"
)

var (
	reFirstLevel = regexp.MustCompile(`(?i)' *OR *'([^']+)' *= *'([^']+)'`)
)

func isFirstSqlInjection(input string) bool {
	matches := reFirstLevel.FindAllStringSubmatch(input, -1)
	if len(matches) != 2 {
		return false
	}
	for _, match := range matches {
		if len(match) != 3 || match[1] != match[2] {
			return false
		}
	}
	return true
}

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
