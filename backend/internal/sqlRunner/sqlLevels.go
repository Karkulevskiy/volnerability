package sqlrunner

import (
	"context"
	"fmt"
	"regexp"
	"volnerability-game/internal/db"
	"volnerability-game/internal/domain"
)

var (
	reFirstLevelInjection = regexp.MustCompile(`(?i)' *OR *'([^']+)' *= *'([^']+)`)
	reFirstLevel          = regexp.MustCompile(`(?i)' *OR *'([^']+)' *= *'([^']+)' AND '([^']+)' *= *'([^']+)'`)

	// reSecondLevelInjection = regexp.MustCompile(`^'\s*UNION\s+SELECT\s+username,\s+password\s+FROM\s+users\s*--.*$`) // Убрать username
	reSecondLevelInjection = regexp.MustCompile(`^'\s*UNION\s+SELECT\s+password\s+FROM\s+users\s*--.*$`) // Убрать username
	reSecondLevel          = regexp.MustCompile(`(?i)^SELECT\s+name\s*,\s*email\s+FROM\s+clients\s+WHERE\s+name\s+LIKE\s+'%'\s+UNION\s+SELECT\s+username\s*,\s*password\s+FROM\s+users--'?%'\s*$`)

	reThirdLevelInjection = regexp.MustCompile(`^';\s*DROP\s+TABLE\s+users;\s*--.*$`)
	reThirdLevel          = regexp.MustCompile(`^';\s*DROP\s+TABLE\s+users;\s*--.*$`)
)

var (
	levelsMap = map[int]func(string) bool{
		6: isFirstSqlInjection,
		7: isSecondSqlInjection,
		8: isThirdSqlInjection,
	}
)

func runLevel(levelId int, input string) (domain.Response, error) {
	re, ok := levelsMap[levelId]
	if !ok {
		return domain.Response{}, fmt.Errorf("failed re for levelId: %d", levelId)
	}

	if ok := re(input); !ok {
		return domain.NewResponseBadRequest("invalid sql query"), nil
	}

	return domain.NewResponseOK(), nil
}

func isFirstSqlInjection(input string) bool {
	matches := reFirstLevelInjection.FindAllStringSubmatch(input, -1)
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

func NewTask(db *db.Storage, levelId int, input string) func(context.Context) (domain.Response, error) {
	return func(ctx context.Context) (domain.Response, error) {
		const op = "sqlRunner.NewTask.validation"
		return runLevel(levelId, input)
	}
}
