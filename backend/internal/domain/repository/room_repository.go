package repository

import (
	"context"

	"github.com/ABfry/album-battler/backend/internal/domain/entity"
	"github.com/google/uuid"
)

type RoomRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Room, error)
	FindByRoomNumber(ctx context.Context, roomNumber int) (*entity.Room, error)
	FindByHostUserID(ctx context.Context, hostUserID uuid.UUID) (*entity.Room, error)
	FindAll(ctx context.Context) ([]*entity.Room, error)
	Save(ctx context.Context, room *entity.Room) error
}
