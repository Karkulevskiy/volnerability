package sqlrunner

import (
	"context"
	"fmt"
	"regexp"
	"volnerability-game/internal/common"
	"volnerability-game/internal/db"
)

var (
	reFirstLevel  = regexp.MustCompile(`(?i)' *OR *'([^']+)' *= *'([^']+)'`)
	reSecondLevel = regexp.MustCompile(`^'\s*UNION\s+SELECT\s+username,\s+password\s+FROM\s+users\s*--.*$`)
	reThirdLevel  = regexp.MustCompile(`^';\s*DROP\s+TABLE\s+users;\s*--.*$`)
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

func isSecondSqlInjection(input string) bool {
	return reSecondLevel.MatchString(input)
}

func isThirdSqlInjection(input string) bool {
	return reThirdLevel.MatchString(input)
}

var (
	levelsCheckers = []func(input string) bool{
		func(input string) bool {
			return isFirstSqlInjection(input)
		},
		func(input string) bool {
			return isSecondSqlInjection(input)
		},
		func(input string) bool {
			return isThirdSqlInjection(input)
		},
	}
)

func NewTask(storage *db.Storage, levelId int, input string) (func(context.Context) (string, bool, error), error) {
	return func(ctx context.Context) (string, bool, error) {
		const op = "sqlRunner.NewTask.validation"
		if levelId >= len(levelsCheckers) {
			return "", false, fmt.Errorf("%s: level id is bigger then availabe sql levels. levelId: %d, total levels: %d", op, levelId, len(levelsCheckers))
		}
		passed := levelsCheckers[levelId](input)
		if passed {
			return common.LevelPassed, true, nil
		}
		return common.LevelNotPassed, false, nil
		// TODO решил оставить на всякий
		// level, err := storage.Level(ctx, levelId)
		// if err != nil {
		// 	return "", false, err
		// }
		// if err := storage.IsQueryValid(input); err != nil {
		// 	if db.IsSyntaxError(err) {
		// 		return "", false, common.NewBadSubmitErr("invalid SQL syntax")
		// 	}
		// 	return "", false, err
		// }
		// if level.ExpectedInput != input {
		// 	return "", false, common.NewBadSubmitErr("not expected sql injection")
		// }
		// return "", true, nil
	}, nil
}
