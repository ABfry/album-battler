package repository

import (
	"github.com/ABfry/album-battler/backend/internal/domain/entity"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user *entity.User) error
	FindByID(id uuid.UUID) (*entity.User, error)
	// todo
}