package repository

import (
	"context"

	"github.com/ABfry/album-battler/backend/internal/domain/entity"
	"github.com/google/uuid"
)

type RoomRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Room, error)
	FindByRoomNumber(ctx context.Context, roomNumber int) (*entity.Room, error)
	// FindAll は全ルームを返す。利用側でソート/フィルタが必要な場合は別途処理する。
	FindAll(ctx context.Context) ([]*entity.Room, error)
	Save(ctx context.Context, room *entity.Room) error
}
