package repository

import (
	"context"

	"github.com/ABfry/album-battler/backend/internal/domain/entity"
)

type RoomRepository interface {
	Create(ctx context.Context, room *entity.Room) error
	FindByID(ctx context.Context, id string) (*entity.Room, error)
	FindByRoomID(ctx context.Context, roomID int) (*entity.Room, error)
	Save(ctx context.Context, room *entity.Room) error
}