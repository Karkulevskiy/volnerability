package sqlrunner

import (
	"context"
	"fmt"
	"regexp"
	"volnerability-game/internal/db"
	"volnerability-game/internal/domain"
)

var (
	reFirstLevelInjection = regexp.MustCompile(`(?i)' *OR *'([^']+)' *= *'([^']+)'`)
	reFirstLevel          = regexp.MustCompile(`(?i)' *OR *'([^']+)' *= *'([^']+)' AND '([^']+)' *= *'([^']+)'`)

	reSecondLevelInjection = regexp.MustCompile(`^'\s*UNION\s+SELECT\s+username,\s+password\s+FROM\s+users\s*--.*$`)
	reSecondLevel          = regexp.MustCompile(`(?i)^SELECT\s+name\s*,\s*email\s+FROM\s+clients\s+WHERE\s+name\s+LIKE\s+'%'\s+UNION\s+SELECT\s+username\s*,\s*password\s+FROM\s+users--'?%'\s*$`)

	reThirdLevelInjection = regexp.MustCompile(`^';\s*DROP\s+TABLE\s+users;\s*--.*$`)
	reThirdLevel          = regexp.MustCompile(`^';\s*DROP\s+TABLE\s+users;\s*--.*$`)
)

var (
	levelsMap = map[int]*regexp.Regexp{
		7: reFirstLevelInjection,
		8: reSecondLevelInjection,
		9: reThirdLevelInjection,
	}
)

func runLevel(ctx context.Context, db *db.Storage, levelId int, input string) (domain.Response, error) {
	re, ok := levelsMap[levelId]
	if !ok {
		return domain.Response{}, fmt.Errorf("failed re for levelId: %d", levelId)
	}

	if similar := re.MatchString(input); !similar {
		return domain.NewResponseBadRequest("invalid sql query"), nil
	}
	// check prepare sql
	// db.IsQueryValid()
	return domain.NewResponseOK(), nil
}

func isFirstSqlInjection(input string) bool {
	matches := reFirstLevelInjection.FindAllStringSubmatch(input, -1)
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
	return reSecondLevelInjection.MatchString(input)
}

func isThirdSqlInjection(input string) bool {
	return reThirdLevelInjection.MatchString(input)
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

func NewTask(db *db.Storage, levelId int, input string) func(context.Context) (domain.Response, error) {
	return func(ctx context.Context) (domain.Response, error) {
		const op = "sqlRunner.NewTask.validation"
		if levelId >= len(levelsCheckers) {
			return domain.Response{}, fmt.Errorf("%s: level id is bigger then availabe sql levels. levelId: %d, total levels: %d", op, levelId, len(levelsCheckers))
		}
		return runLevel(ctx, db, levelId, input)
		// passed := levelsCheckers[levelId](input)
		// if passed {
		// 	return domain.NewResponseOK(), nil
		// }
		// TODO неправильный респонс, нужно подумать, что делать с sql
		// return domain.NewResponseBadRequest("invalid sql injection"), nil
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
	}
}
