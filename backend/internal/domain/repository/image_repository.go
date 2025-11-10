package repository

import (
	"context"

	"github.com/ABfry/album-battler/backend/internal/domain/entity"
	"github.com/google/uuid"
)

type ImageRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Image, error)
	FindByUserIDs(ctx context.Context, userID uuid.UUID) ([]*entity.Image, error)
	FindByBattleIDs(ctx context.Context, battleID uuid.UUID) ([]*entity.Image, error)
	Save(ctx context.Context, image *entity.Image) error
}
