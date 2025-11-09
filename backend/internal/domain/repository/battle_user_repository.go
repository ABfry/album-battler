package repository

import (
	"context"

	"github.com/google/uuid"
)

type BattleUserRepository interface {
	Save(ctx context.Context, battleID, userID uuid.UUID) error
}
