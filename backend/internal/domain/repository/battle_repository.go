package repository

import (
	"context"

	"github.com/ABfry/album-battler/backend/internal/domain/entity"
	"github.com/google/uuid"
)

type BattleRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Battle, error)
	FindByRoomID(ctx context.Context, roomID uuid.UUID) (*entity.Battle, error)
	Save(ctx context.Context, battle *entity.Battle) error
}
