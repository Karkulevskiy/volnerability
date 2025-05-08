package hints

import (
	"context"
	"volnerability-game/internal/db"
	"volnerability-game/internal/domain"
)

func Run(ctx context.Context, storage *db.Storage, hintId int) (domain.Hint, error) {
	return storage.Hint(ctx, hintId)
}
