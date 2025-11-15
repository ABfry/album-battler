package repository

import (
	"context"

	"github.com/google/uuid"
)

type BattleUserRepository interface {
	SaveBatch(ctx context.Context, battleID uuid.UUID, userIDs []uuid.UUID) error
}
