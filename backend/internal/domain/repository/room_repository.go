package repository

import "github.com/ABfry/album-battler/backend/internal/domain/entity"

type RoomRepository interface {
	Create(room *entity.Room) error
	FindByID(id string) (*entity.Room, error)
	FindByRoomID(roomID int) (*entity.Room, error)
	Save(room *entity.Room) error
}