package repository

import (
	"context"

	"github.com/ABfry/album-battler/backend/internal/domain/entity"
	"github.com/google/uuid"
)

type UserRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	FindByIDs(ctx context.Context, ids []uuid.UUID) ([]*entity.User, error)
	Save(ctx context.Context, user *entity.User) error
}
