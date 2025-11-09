package repository

import (
	"context"

	"github.com/ABfry/album-battler/backend/internal/domain/entity"
)

type BattleRepository interface {
	FindByID(ctx context.Context, id string) (*entity.Battle, error)
	FindByRoomID(ctx context.Context, roomID int) (*entity.Battle, error)
	Save(ctx context.Context, battle *entity.Battle) error
}
