package repository

import "github.com/ABfry/album-battler/backend/internal/domain/entity"


type UserRepository interface {
	Create(user *entity.User) error
	FindByID(id string) (*entity.User, error)
	// todo
}