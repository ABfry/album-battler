package repository

import (
	"context"

	"github.com/ABfry/album-battler/backend/internal/domain/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id string) (*entity.User, error)
	Save(ctx context.Context, user *entity.User) error
}
