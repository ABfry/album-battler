package repository

import (
	"context"

	"github.com/ABfry/album-battler/backend/internal/domain/entity"
	"github.com/google/uuid"
)

type BattleUserRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*entity.BattleUser, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) (*entity.BattleUser, error)
	FindByBattleID(ctx context.Context, battleID uuid.UUID) (*entity.BattleUser, error)
	Save(ctx context.Context, battleUser *entity.BattleUser) error
}
