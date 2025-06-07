package levels

import (
	"context"
	"fmt"
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
	Lang    string `json:"lang"`
}

type Submit func(ctx context.Context) (any, error)

func New(r Request, db *db.Storage, codeRunner *coderunner.CodeRunner) (Submit, error) {
	if err := validate(r); err != nil {
		return nil, err
	}
	// TODO need to create level groups
	if r.LevelId < 3 {
		return codeRunner.NewTask(r.Input, r.Lang, r.Id)
	}
	if r.LevelId < 8 {
		return sqlrunner.NewTask(db, r.LevelId, r.Input)
	}
	// TODO other tasks
	return nil, nil
}

func Level(ctx context.Context, db *db.Storage, levelId int) (domain.Level, error) {
	const op = "level.getLevelById"
	level, err := db.Level(ctx, levelId)
	if err != nil {
		return domain.Level{}, fmt.Errorf("op: %s. failed to proceed get level by id: %d", op, levelId)
	}
	return level, nil
	// TODO новое поле в бд
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
